package initialize

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/config"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, *redis.Client) {
	commentClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisCommentClientDB,
	})
	favoriteClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalServerConfig.RedisInfo.Host, config.GlobalServerConfig.RedisInfo.Port),
		Password: config.GlobalServerConfig.RedisInfo.Password,
		DB:       consts.RedisFavoriteClientDB,
	})
	return commentClient, favoriteClient
}
