package pkg

import (
	"context"
	"fmt"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/go-redis/redis/v8"
)

type FavoriteRedisManager struct {
	RedisClient *redis.Client
}

func (r *FavoriteRedisManager) FavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	count, err := r.RedisClient.SCard(ctx, videoIdStr).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (r *FavoriteRedisManager) Like(ctx context.Context, userId int64, videoId int64) error {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	pl := r.RedisClient.TxPipeline()
	if err := pl.SAdd(ctx, userIdStr, videoIdStr).Err(); err != nil {
		return err
	}
	if err := pl.SAdd(ctx, videoIdStr, userId).Err(); err != nil {
		return err
	}
	if _, err := pl.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRedisManager) Unlike(ctx context.Context, userId int64, videoId int64) error {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	pl := r.RedisClient.TxPipeline()
	if err := pl.ZRem(ctx, userIdStr, videoId).Err(); err != nil {
		return err
	}
	if err := pl.SRem(ctx, videoIdStr, userId).Err(); err != nil {
		return err
	}
	if _, err := pl.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRedisManager) Check(ctx context.Context, userId int64, videoId int64) (bool, error) {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	flag1, err := r.RedisClient.SIsMember(ctx, userIdStr, videoId).Result()
	if err != nil {
		return false, err
	}
	flag2, err := r.RedisClient.SIsMember(ctx, videoIdStr, userId).Result()
	if err != nil {
		return false, err
	}
	if flag1 != flag2 {
		return false, errno.InteractionServerErr.WithMessage("dirty data in redis")
	}
	return flag1, nil
}

func (r *FavoriteRedisManager) GetFavoriteVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	userIdStr := fmt.Sprintf("%d", userId)
	args := r.RedisClient.SMembers(ctx, userIdStr).Args()
	list, err := redis.NewIntSliceCmd(ctx, args).Result()
	if err != nil {
		return nil, err
	}
	return list, err
}
