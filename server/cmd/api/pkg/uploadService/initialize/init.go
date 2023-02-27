package initialize

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/pkg"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Init() *uploadService.Service {
	initLogger()
	initConfig()
	minioClient := initMinio()
	amqpC := initMq()

	publisher, err := pkg.NewPublisher(amqpC, config.GlobalServiceConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create publisher")
	}
	subscriber, err := pkg.NewSubscriber(amqpC, config.GlobalServiceConfig.RabbitMqInfo.Exchange)
	if err != nil {
		klog.Fatal("cannot create subscriber")
	}

	return uploadService.NewUploadService(minioClient, subscriber, publisher, config.GlobalServiceConfig)
}
