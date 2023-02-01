package tools

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
)

type VideoManager struct {
	VideoService videoservice.Client
}

// GetVideos gets video list.
func (m *VideoManager) GetVideos(ctx context.Context, list []int64, viewerId int64) ([]*base.Video, error) {
	var videos []*base.Video
	for _, vid := range list {
		v, err := m.VideoService.GetVideo(ctx, &video.DouyinGetVideoRequest{
			VideoId:  vid,
			ViewerId: viewerId,
		})
		if err != nil {
			return nil, err
		}
		videos = append(videos, v.Video)
	}
	return videos, nil
}
