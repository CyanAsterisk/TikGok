package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	AmqpConn     *amqp.Connection

	UserClient userservice.Client
	ChatClient chatservice.Client
)
