package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
)

type InteractionManager struct {
	InteractionService interactionserver.Client
}

// GetCommentCount  get one video comment.
func (i *InteractionManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	resp, err := i.InteractionService.GetCommentCount(ctx, &interaction.DouyinGetCommentCountRequest{VideoId: videoId})
	if err != nil {
		return 0, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Count, nil
}

// CheckFavorite check one favorite the video or not.
func (i *InteractionManager) CheckFavorite(ctx context.Context, userId int64, videoId int64) (bool, error) {
	resp, err := i.InteractionService.CheckFavorite(ctx, &interaction.DouyinCheckFavoriteRequest{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		return false, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return false, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Check, nil
}

// GetFavoriteCount get the favorite num of the video.
func (i *InteractionManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	resp, err := i.InteractionService.GetFavoriteCount(ctx, &interaction.DouyinGetFavoriteCountRequest{VideoId: videoId})
	if err != nil {
		return 0, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Count, nil
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
	return resp.VideoList, nil
}

func (i *InteractionManager) BatchGetCommentCount(ctx context.Context, videoIdList []int64) ([]int64, error) {
	resp, err := i.InteractionService.BatchGetCommentCount(ctx, &interaction.DouyinBatchGetCommentCountRequest{VideoIdList: videoIdList})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.CountList, nil
}
func (i *InteractionManager) BatchCheckFavorite(ctx context.Context, userId int64, videoIdList []int64) ([]bool, error) {
	resp, err := i.InteractionService.BatchCheckFavorite(ctx, &interaction.DouyinBatchCheckFavoriteRequest{
		UserId:      userId,
		VideoIdList: videoIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.CheckList, nil
}
func (i *InteractionManager) BatchGetFavoriteCount(ctx context.Context, videoId []int64) ([]int64, error) {
	resp, err := i.InteractionService.BatchGetFavoriteCount(ctx, &interaction.DouyinBatchGetFavoriteCountRequest{VideoIdList: videoId})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.CountList, nil
}
