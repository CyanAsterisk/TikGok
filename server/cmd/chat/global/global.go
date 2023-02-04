package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/config"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB                 *gorm.DB
	ServerConfig       config.ServerConfig
	NacosConfig        config.NacosConfig
	AmqpConn           *amqp.Connection
	RedisSentClient    *redis.Client
	RedisReceiveClient *redis.Client
)
