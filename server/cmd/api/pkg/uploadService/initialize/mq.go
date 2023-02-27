package initialize

import (
	"fmt"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/pkg/uploadService/config"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// InitMq to init rabbitMQ
func initMq() *amqp.Connection {
	info := config.GlobalServiceConfig.RabbitMqInfo
	amqpConn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqURI, info.User, info.Password, info.Host, info.Port))
	if err != nil {
		klog.Fatal("cannot dial amqp", err)
	}
	return amqpConn
}
