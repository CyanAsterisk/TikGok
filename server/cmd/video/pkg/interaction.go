package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
)

type InteractionManager struct {
	InteractionService interactionserver.Client
}

// GetVideoInteractInfo get video interactInfo.
func (i *InteractionManager) GetVideoInteractInfo(ctx context.Context, videoId, viewerId int64) (*base.VideoInteractInfo, error) {
	resp, err := i.InteractionService.GetVideoInteractInfo(ctx, &interaction.DouyinGetVideoInteractInfoRequest{
		VideoId:  videoId,
		ViewerId: viewerId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.InteractInfo, nil
}

// GetFavoriteVideoIdList gets the favorite video id list.
func (i *InteractionManager) GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error) {
	resp, err := i.InteractionService.GetFavoriteVideoIdList(ctx, &interaction.DouyinGetFavoriteVideoIdListRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.VideoIdList, nil
}

// BatchGetVideoInteractInfo batch get video interactInfo.
func (i *InteractionManager) BatchGetVideoInteractInfo(ctx context.Context, videoIdList []int64, viewerId int64) ([]*base.VideoInteractInfo, error) {
	resp, err := i.InteractionService.BatchGetInteractInfo(ctx, &interaction.DouyinBatchGetInteractInfoRequest{
		VideoIdList: videoIdList,
		ViewerId:    viewerId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.InteractInfoList, nil
}
