package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisClient *redis.Client
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{RedisClient: client}
}

func (r *RedisManager) CreateVideo(ctx context.Context, video *model.Video) error {
	authorIdStr := fmt.Sprintf("%d", video.AuthorId)
	videoIdStr := fmt.Sprintf("%d", video.ID)
	_, err := r.RedisClient.Get(ctx, videoIdStr).Result()
	if err == nil {
		return errno.VideoServerErr.WithMessage("video record already exists")
	} else if err != redis.Nil {
		return err
	}
	pl := r.RedisClient.TxPipeline()

	z := &redis.Z{
		Score:  float64(video.CreateTime),
		Member: video.ID,
	}
	if err = pl.ZAdd(ctx, authorIdStr, z).Err(); err != nil {
		return err
	}
	if err = pl.ZAdd(ctx, consts.AllVideoSortSetKey, z).Err(); err != nil {
		return err
	}
	videoRecord, err := sonic.Marshal(video)
	if err != nil {
		return errno.VideoServerErr.WithMessage("marshal video error")
	}
	if err = pl.Set(ctx, videoIdStr, videoRecord, 0).Err(); err != nil {
		return err
	}
	_, err = pl.Exec(ctx)
	return err
}

func (r *RedisManager) DeleteVideoById(ctx context.Context, videoId int64) error {
	v, err := r.GetVideoByVideoId(ctx, videoId)
	if err != nil {
		return err
	}
	pl := r.RedisClient.TxPipeline()
	authorIdStr := fmt.Sprintf("%d", v.AuthorId)
	videoIdStr := fmt.Sprintf("%d", v.ID)

	if err = pl.ZRem(ctx, authorIdStr, v.ID).Err(); err != nil {
		return err
	}
	if err = pl.ZRem(ctx, consts.AllVideoSortSetKey, v.ID).Err(); err != nil {
		return err
	}
	if err = pl.Del(ctx, videoIdStr).Err(); err != nil {
		return err
	}
	_, err = pl.Exec(ctx)
	return err
}

func (r *RedisManager) GetVideoListByLatestTime(ctx context.Context, latestTime int64) ([]*model.Video, error) {
	op := &redis.ZRangeBy{
		Min:    "",
		Max:    fmt.Sprintf("%d", latestTime),
		Offset: 0,
		Count:  consts.VideosLimit,
	}
	vidList := make([]int64, 0)
	err := r.RedisClient.ZRevRangeByScore(ctx, consts.AllVideoSortSetKey, op).ScanSlice(&vidList)
	if err != nil {
		return nil, err
	}
	videoList, err := r.BatchGetVideoByVideoId(ctx, vidList)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (r *RedisManager) GetVideoListByAuthorId(ctx context.Context, authorId int64) ([]*model.Video, error) {
	vidList := make([]int64, 0)
	err := r.RedisClient.ZRevRange(ctx, fmt.Sprintf("%d", authorId), 0, -1).ScanSlice(&vidList)
	if err != nil {
		return nil, err
	}
	videoList, err := r.BatchGetVideoByVideoId(ctx, vidList)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (r *RedisManager) GetVideoIdListByAuthorId(ctx context.Context, authorId int64) ([]int64, error) {
	vidList := make([]int64, 0)
	err := r.RedisClient.ZRevRange(ctx, fmt.Sprintf("%d", authorId), 0, -1).ScanSlice(&vidList)
	if err != nil {
		return nil, err
	}
	return vidList, nil
}

func (r *RedisManager) GetVideoByVideoId(ctx context.Context, videoId int64) (*model.Video, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	videoJSONStr, err := r.RedisClient.Get(ctx, videoIdStr).Result()
	if err != nil {
		return nil, err
	}
	var video model.Video
	err = sonic.Unmarshal([]byte(videoJSONStr), &video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *RedisManager) BatchGetVideoByVideoId(ctx context.Context, videoIdList []int64) ([]*model.Video, error) {
	if videoIdList == nil {
		return nil, nil
	}
	vl := make([]*model.Video, 0)
	for _, vid := range videoIdList {
		video, err := r.GetVideoByVideoId(ctx, vid)
		if err != nil {
			return nil, err
		}
		vl = append(vl, video)
	}
	return vl, nil
}
