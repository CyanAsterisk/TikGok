package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	AmqpConn     *amqp.Connection
	RedisClient  *redis.Client

	InteractClient interactionserver.Client
	UserClient     userservice.Client
)
