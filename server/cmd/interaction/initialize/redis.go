package initialize

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/global"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	global.RedisCommentClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisCommentClientDB,
	})
	global.RedisFavoriteClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password,
		DB:       consts.RedisFavoriteClientDB,
	})
}
