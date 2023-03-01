package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNoSuchRecord       = errors.New("no such record")
	ErrRecordAlreadyExist = errors.New("record already exist")
)

type CommentRedisManager struct {
	RedisClient *redis.Client
}

func NewCommentRedisManager(client *redis.Client) *CommentRedisManager {
	return &CommentRedisManager{
		RedisClient: client,
	}
}

func (r *CommentRedisManager) CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	count, err := r.RedisClient.ZCard(ctx, videoIdStr).Result()
	if err != nil {
		return -1, err
	}
	return count, err
}

func (r *CommentRedisManager) CreateComment(ctx context.Context, comment *model.Comment) error {
	commentIdStr := fmt.Sprintf("%d", comment.ID)
	err := r.RedisClient.Get(ctx, commentIdStr).Err()
	if err == nil {
		return ErrRecordAlreadyExist
	}
	if err != redis.Nil {
		return err
	}
	videoIdStr := fmt.Sprintf("%d", comment.VideoId)
	commentJson, err := sonic.Marshal(comment)
	if err != nil {
		return err
	}
	pl := r.RedisClient.TxPipeline()
	if err := pl.ZAdd(ctx, videoIdStr, &redis.Z{
		Score:  float64(comment.CreateDate),
		Member: commentIdStr,
	}).Err(); err != nil {
		return err
	}
	if err := pl.Set(ctx, commentIdStr, commentJson, 0).Err(); err != nil {
		return err
	}
	_, err = pl.Exec(ctx)
	return err
}

func (r *CommentRedisManager) DeleteComment(ctx context.Context, commentId int64) error {
	commentIdStr := fmt.Sprintf("%d", commentId)
	commentJson, err := r.RedisClient.Get(ctx, commentIdStr).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrNoSuchRecord
		}
		return err
	}
	var comment model.Comment
	if err = sonic.Unmarshal([]byte(commentJson), &comment); err != nil {
		return err
	}
	videoIdStr := fmt.Sprintf("%d", comment.VideoId)
	pl := r.RedisClient.TxPipeline()
	if err = pl.ZRem(ctx, videoIdStr, commentIdStr).Err(); err != nil {
		return err
	}
	if err = pl.Del(ctx, commentIdStr).Err(); err != nil {
		return err
	}
	_, err = pl.Exec(ctx)
	return err
}

func (r *CommentRedisManager) GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	commentIdList, err := r.RedisClient.ZRevRange(ctx, videoIdStr, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	length := len(commentIdList)
	cl := make([]*model.Comment, length)
	var wg sync.WaitGroup
	wg.Add(length)
	for i := 0; i < length; i++ {
		go func(idx int) {
			defer wg.Done()
			var c model.Comment
			commentJson, _ := r.RedisClient.Get(ctx, commentIdList[idx]).Result()
			err = sonic.Unmarshal([]byte(commentJson), &c)
			cl[idx] = &c
		}(i)
	}
	wg.Wait()
	return cl, err
}
