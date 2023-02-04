package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	AmqpConn     *amqp.Connection

	VideoClient videoservice.Client
)
