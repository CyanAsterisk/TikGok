package global

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/config"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

var (
	DB                  *gorm.DB
	ServerConfig        config.ServerConfig
	NacosConfig         config.NacosConfig
	AmqpConn            *amqp.Connection
	RedisCommentClient  *redis.Client
	RedisFavoriteClient *redis.Client
	VideoClient         videoservice.Client
)
