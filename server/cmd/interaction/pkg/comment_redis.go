package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

type CommentRedisManager struct {
	RedisClient *redis.Client
}

const (
	videoIdFiled     = "videoId"
	commentJsonFiled = "commentJson"
)

func (r *CommentRedisManager) CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	count, err := r.RedisClient.ZCard(ctx, videoIdStr).Result()
	if err != nil {
		return -1, err
	}
	return count, err
}
func (r *CommentRedisManager) CreateComment(ctx context.Context, comment *model.Comment) error {
	pl := r.RedisClient.TxPipeline()
	videoIdStr := fmt.Sprintf("%d", comment.VideoId)
	commentIdStr := fmt.Sprintf("%d", comment.ID)
	commentJson, err := sonic.Marshal(comment)
	if err != nil {
		return err
	}
	batchData := make(map[string]string)
	batchData[videoIdFiled] = videoIdStr
	batchData[commentJsonFiled] = string(commentJson)
	if err = pl.ZAdd(ctx, videoIdStr, &redis.Z{
		Score:  float64(comment.CreateDate.UnixNano()),
		Member: commentJson,
	}).Err(); err != nil {
		return err
	}
	if err = pl.HMSet(ctx, commentIdStr, batchData).Err(); err != nil {
		return err
	}
	if _, err = pl.Exec(ctx); err != nil {
		return err
	}
	return nil
}
func (r *CommentRedisManager) DeleteComment(ctx context.Context, commentId int64) error {
	commentIdStr := fmt.Sprintf("%d", commentId)
	values, err := r.RedisClient.HMGet(ctx, commentIdStr, videoIdFiled, commentJsonFiled).Result()
	if err != nil {
		return err
	}
	videoIdStr := values[0].(string)
	commentJson := values[1].(string)
	if err = r.RedisClient.ZRem(ctx, videoIdStr, commentJson).Err(); err != nil {
		return err
	}
	return nil
}
func (r *CommentRedisManager) GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	values, err := r.RedisClient.ZRange(ctx, videoIdStr, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var cl []*model.Comment
	for _, val := range values {
		var c model.Comment
		err = sonic.Unmarshal([]byte(val), &c)
		if err != nil {
			return nil, err
		}
		cl = append(cl, &c)
	}
	return cl, nil
}
