package main

import (
	"context"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/dao"
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
	db := initialize.InitDB()
	rC := initialize.InitRedis()
	amqpC := initialize.InitMq()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	userC := initialize.InitUser()
	chatC := initialize.InitChat()

	publisher, err := pkg.NewPublisher(amqpC, config.GlobalServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher", err)
	}
	subscriber, err := pkg.NewSubscriber(amqpC, config.GlobalServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber", err.Error())
	}

	followDao := dao.NewFollow(db)
	go pkg.SubscribeRoutine(subscriber, followDao)

	impl := &SocialityServiceImpl{
		UserManager: &pkg.UserManager{
			UserService: userC,
			ChatService: chatC,
		},
		Publisher:    publisher,
		RedisManager: pkg.NewRedisManager(rC),
		Dao:          followDao,
	}
	// Create new server.
	srv := sociality.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)

	err = srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
