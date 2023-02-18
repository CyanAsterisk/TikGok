package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video/videoservice"
)

type VideoManager struct {
	client videoservice.Client
}

func NewVideoManager(client videoservice.Client) *VideoManager {
	return &VideoManager{client: client}
}

func (m *VideoManager) GetPublishedVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	resp, err := m.client.GetPublishedVideoIdList(ctx, &video.DouyinGetPublishedVideoIdListRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.VideoServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.VideoIdList, nil
}
