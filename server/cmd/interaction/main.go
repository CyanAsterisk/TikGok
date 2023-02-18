package main

import (
	"context"
	"net"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/global"
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
	initialize.InitDB()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(global.ServerConfig.Name),
		provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	initialize.InitRedis()
	initialize.InitMq()
	initialize.InitVideo()

	mqInfo := global.ServerConfig.RabbitMqInfo
	commentPublisher, err := pkg.NewCommentPublisher(global.AmqpConn, mqInfo.CommentExchange)
	if err != nil {
		klog.Fatal("cannot create comment publisher")
	}
	favoritePublisher, err := pkg.NewFavoritePublisher(global.AmqpConn, mqInfo.FavoriteExchange)
	if err != nil {
		klog.Fatal("cannot create favorite publisher")
	}
	commentSubscriber, err := pkg.NewCommentSubscriber(global.AmqpConn, mqInfo.CommentExchange)
	if err != nil {
		klog.Fatal("cannot create comment subscriber")
	}

	commentDao := dao.NewComment(global.DB)
	favoriteDao := dao.NewFavorite(global.DB)

	go func() {
		if err = pkg.CommentSubscribeRoutine(commentSubscriber, commentDao); err != nil {
			klog.Fatal("comment subscribe routine", err)
		}
	}()

	favoriteSubscriber, err := pkg.NewFavoriteSubscriber(global.AmqpConn, mqInfo.FavoriteExchange)
	if err != nil {
		klog.Fatal("cannot create favorite subscriber")
	}
	go func() {
		if err = pkg.FavoriteSubscribeRoutine(favoriteSubscriber, favoriteDao); err != nil {
			klog.Fatal("favorite subscribe routine", err)
		}
	}()

	impl := &InteractionServerImpl{
		VideoManager:         pkg.NewVideoManager(global.VideoClient),
		CommentPublisher:     commentPublisher,
		FavoritePublisher:    favoritePublisher,
		CommentRedisManager:  pkg.NewCommentRedisManager(global.RedisCommentClient),
		FavoriteRedisManager: pkg.NewFavoriteRedisManager(global.RedisFavoriteClient),
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
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.ServerConfig.Name}),
	)

	err = srv.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
