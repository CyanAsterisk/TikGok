package main

import (
	"context"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/initialize"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	sociality "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {
	// initialization
	initialize.InitLogger()
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitNacos(Port)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitMq()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.ServerConfig.Name),
		provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	initialize.InitUser()
	initialize.InitChat()

	Publisher, err := pkg.NewPublisher(global.AmqpConn, global.ServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher", err)
	}
	Subscriber, err := pkg.NewSubscriber(global.AmqpConn, global.ServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber", err.Error())
	}
	go pkg.SubscribeRoutine(Subscriber)

	impl := &SocialityServiceImpl{
		UserManager: &pkg.UserManager{
			UserService: global.UserClient,
			ChatService: global.ChatClient,
		},
		Publisher:  Publisher,
		Subscriber: Subscriber,
		RedisManager: &pkg.RedisManager{
			RedisFollowingClient: global.RedisFollowingClient,
			RedisFollowerClient:  global.RedisFollowerClient,
		},
	}
	// Create new server.
	srv := sociality.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.ServerConfig.Name}),
	)

	err = srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
