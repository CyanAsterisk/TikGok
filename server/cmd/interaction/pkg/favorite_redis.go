package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/go-redis/redis/v8"
)

type FavoriteRedisManager struct {
	client *redis.Client
}

func NewFavoriteRedisManager(client *redis.Client) *FavoriteRedisManager {
	return &FavoriteRedisManager{
		client: client,
	}
}

func (r *FavoriteRedisManager) GetFavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	count, err := r.client.ZCard(ctx, videoIdStr).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (r *FavoriteRedisManager) Like(ctx context.Context, userId, videoId, time int64) error {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	pl := r.client.TxPipeline()
	if err := pl.ZAdd(ctx, userIdStr, &redis.Z{
		Score:  float64(time),
		Member: videoId,
	}).Err(); err != nil {
		return err
	}
	if err := pl.ZAdd(ctx, videoIdStr, &redis.Z{
		Score:  float64(time),
		Member: userId,
	}).Err(); err != nil {
		return err
	}
	if _, err := pl.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRedisManager) Unlike(ctx context.Context, userId, videoId int64) error {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	pl := r.client.TxPipeline()
	if err := pl.ZRem(ctx, userIdStr, videoId).Err(); err != nil {
		return err
	}
	if err := pl.ZRem(ctx, videoIdStr, userId).Err(); err != nil {
		return err
	}
	if _, err := pl.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *FavoriteRedisManager) Check(ctx context.Context, userId, videoId int64) (bool, error) {
	userIdStr := fmt.Sprintf("%d", userId)
	videoIdStr := fmt.Sprintf("%d", videoId)
	_, err1 := r.client.ZScore(ctx, userIdStr, videoIdStr).Result()
	if err1 != nil && err1 != redis.Nil {
		return false, err1
	}
	_, err2 := r.client.ZScore(ctx, videoIdStr, userIdStr).Result()
	if err2 != nil && err2 != redis.Nil {
		return false, err2
	}
	if err1 != err2 {
		return false, errno.InteractionServerErr.WithMessage("dirty data in redis")
	}
	return err1 == nil, nil
}

func (r *FavoriteRedisManager) GetFavoriteVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error) {
	vidList := make([]int64, 0)
	userIdStr := fmt.Sprintf("%d", userId)
	err := r.client.ZRevRange(ctx, userIdStr, 0, -1).ScanSlice(&vidList)
	if err != nil {
		return nil, err
	}
	return vidList, err
}

func (r *FavoriteRedisManager) GetFavoriteVideoCountByUserId(ctx context.Context, userId int64) (int64, error) {
	UserIdStr := fmt.Sprintf("%d", userId)
	count, err := r.client.ZCard(ctx, UserIdStr).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}
