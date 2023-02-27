package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/config"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/initialize"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	interaction "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
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
	cC, fC := initialize.InitRedis()
	amqpC := initialize.InitMq()
	videoC := initialize.InitVideo()

	mqInfo := config.GlobalServerConfig.RabbitMqInfo
	commentPublisher, err := pkg.NewCommentPublisher(amqpC, mqInfo.CommentExchange)
	if err != nil {
		klog.Fatal("cannot create comment publisher")
	}
	favoritePublisher, err := pkg.NewFavoritePublisher(amqpC, mqInfo.FavoriteExchange)
	if err != nil {
		klog.Fatal("cannot create favorite publisher")
	}
	commentSubscriber, err := pkg.NewCommentSubscriber(amqpC, mqInfo.CommentExchange)
	if err != nil {
		klog.Fatal("cannot create comment subscriber")
	}

	commentDao := dao.NewComment(db)
	favoriteDao := dao.NewFavorite(db)

	go func() {
		if err = pkg.CommentSubscribeRoutine(commentSubscriber, commentDao); err != nil {
			klog.Fatal("comment subscribe routine", err)
		}
	}()

	favoriteSubscriber, err := pkg.NewFavoriteSubscriber(amqpC, mqInfo.FavoriteExchange)
	if err != nil {
		klog.Fatal("cannot create favorite subscriber")
	}
	go func() {
		if err = pkg.FavoriteSubscribeRoutine(favoriteSubscriber, favoriteDao); err != nil {
			klog.Fatal("favorite subscribe routine", err)
		}
	}()

	impl := &InteractionServerImpl{
		VideoManager:         pkg.NewVideoManager(videoC),
		CommentPublisher:     commentPublisher,
		FavoritePublisher:    favoritePublisher,
		CommentRedisManager:  pkg.NewCommentRedisManager(cC),
		FavoriteRedisManager: pkg.NewFavoriteRedisManager(fC),
		CommentDao:           commentDao,
		FavoriteDao:          favoriteDao,
	}
	// Create new server.
	srv := interaction.NewServer(impl,
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
