package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/config"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
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
	db := initialize.InitDB()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	rC := initialize.InitRedis()
	amqpC := initialize.InitMq()
	interactionC := initialize.InitInteraction()
	userC := initialize.InitUser()

	publisher, err := pkg.NewPublisher(amqpC, config.GlobalServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher", err)
	}
	subscriber, err := pkg.NewSubscriber(amqpC, config.GlobalServerConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber", err)
	}
	videoDao := dao.NewVideo(db)
	go func() {
		err = pkg.SubscribeRoutine(subscriber, videoDao)
		if err != nil {
			klog.Fatal("subscribe err", err)
		}
	}()

	impl := &VideoServiceImpl{
		Publisher:          publisher,
		UserManager:        &pkg.UserManager{UserService: userC},
		InteractionManager: &pkg.InteractionManager{InteractionService: interactionC},
		RedisManager:       pkg.NewRedisManager(rC),
		Dao:                videoDao,
	}
	// Create new server.
	srv := video.NewServer(impl,
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
