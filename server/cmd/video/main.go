package main

import (
	"context"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/initialize"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	video "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
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
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.ServerConfig.Name),
		provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	initialize.InitRedis()
	initialize.InitMq()
	initialize.InitInteraction()
	initialize.InitUser()

	publisher, err := pkg.NewPublisher(global.AmqpConn, global.ServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher", err)
	}
	subscriber, err := pkg.NewSubscriber(global.AmqpConn, global.ServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber", err)
	}
	videoDao := dao.NewVideo(global.DB)
	go func() {
		err = pkg.SubscribeRoutine(subscriber, videoDao)
		if err != nil {
			klog.Fatal("subscribe err", err)
		}
	}()

	impl := &VideoServiceImpl{
		Publisher:          publisher,
		UserManager:        &pkg.UserManager{UserService: global.UserClient},
		InteractionManager: &pkg.InteractionManager{InteractionService: global.InteractClient},
		RedisManager:       pkg.NewRedisManager(global.RedisClient),
		Dao:                videoDao,
	}
	// Create new server.
	srv := video.NewServer(impl,
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
