package initialize

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	global.RedisSentClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RabbitMqInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisSentClientDB,
	})
	global.RedisReceiveClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RabbitMqInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisReceiveClientDB,
	})
}
