package initialize

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	global.RedisFollowerClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RabbitMqInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisFollowerClientDB,
	})
	global.RedisFollowingClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RabbitMqInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisFollowingClientDB,
	})
}
