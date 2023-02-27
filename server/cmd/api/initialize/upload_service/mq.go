package upload_service

import (
	"fmt"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// InitMq to init rabbitMQ
func initMq() {
	info := config.GlobalServerConfig.UploadServiceInfo.RabbitMqInfo
	amqpConn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqURI, info.User, info.Password, info.Host, info.Port))
	if err != nil {
		klog.Fatal("cannot dial amqp", err)
	}
	if publisher, err = uploadService.NewPublisher(amqpConn, info.Exchange); err != nil {
		klog.Fatal("create publisher err", err)
	}
	if subscriber, err = uploadService.NewSubscriber(amqpConn, info.Exchange); err != nil {
		klog.Fatal("create subscriber err", err)
	}
}
