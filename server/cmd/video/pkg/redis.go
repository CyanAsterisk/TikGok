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
	pl := r.RedisClient.TxPipeline()
	authorIdStr := fmt.Sprintf("%d", video.AuthorId)
	videoIdStr := fmt.Sprintf("%d", video.ID)
	videoRecord, err := sonic.Marshal(video)
	if err != nil {
		return errno.VideoServerErr.WithMessage("marshal video error")
	}
	if err = pl.LPush(ctx, authorIdStr, videoIdStr).Err(); err != nil {
		return err
	}
	if err = pl.ZAdd(ctx, consts.AllVideoSortSetKey, &redis.Z{
		Score:  float64(video.CreateTime),
		Member: videoRecord,
	}).Err(); err != nil {
		return err
	}
	if err = pl.Set(ctx, videoIdStr, videoRecord, 0).Err(); err != nil {
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
	videoJSONList, err := r.RedisClient.ZRangeByScore(ctx, consts.AllVideoSortSetKey, op).Result()
	if err != nil {
		return nil, err
	}
	videoList, err := videoListJsonToStruct(videoJSONList)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (r *RedisManager) GetVideoListByAuthorId(ctx context.Context, authorId int64) ([]*model.Video, error) {
	videoJSONList, err := r.RedisClient.LRange(ctx, fmt.Sprintf("%d", authorId), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	videoList, err := videoListJsonToStruct(videoJSONList)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func (r *RedisManager) GetVideoByVideoId(ctx context.Context, videoId int64) (*model.Video, error) {
	videoIdStr := fmt.Sprintf("%d", videoId)
	videoJSONStr, err := r.RedisClient.Get(ctx, videoIdStr).Result()
	if err != nil {
		return nil, err
	}
	video := &model.Video{}
	err = sonic.Unmarshal([]byte(videoJSONStr), video)
	if err != nil {
		return nil, err
	}
	return video, nil
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

func videoListJsonToStruct(videoJsonList []string) ([]*model.Video, error) {
	if videoJsonList == nil {
		return nil, nil
	}
	videoList := make([]*model.Video, 0)
	for _, videoJson := range videoJsonList {
		video := &model.Video{}
		if err := sonic.Unmarshal([]byte(videoJson), video); err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
