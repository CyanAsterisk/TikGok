package initialize

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisUserClientDB,
	})
}
