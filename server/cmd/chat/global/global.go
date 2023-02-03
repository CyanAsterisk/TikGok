package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/config"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	AmqpConn     *amqp.Connection
)
