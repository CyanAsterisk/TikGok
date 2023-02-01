package tools

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
)

type InteractionManager struct {
	InteractionService interactionserver.Client
}

// GetCommentCount  get one video comment.
func (i *InteractionManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	resp, err := i.InteractionService.CommentList(ctx, &interaction.DouyinCommentListRequest{VideoId: videoId})
	if err != nil {
		return 0, err
	}
	if int64(resp.BaseResp.StatusCode) != errno.Success.ErrCode {
		return 0, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return int64(len(resp.CommentList)), nil
}

// CheckFavorite check one favorite the video or not.
func (i *InteractionManager) CheckFavorite(ctx context.Context, userId int64, videoId int64) (bool, error) {
	resp, err := i.InteractionService.FavoriteList(ctx, &interaction.DouyinFavoriteListRequest{OwnerId: userId})
	if err != nil {
		return false, err
	}
	if int64(resp.BaseResp.StatusCode) != errno.Success.ErrCode {
		return false, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	for _, v := range resp.VideoList {
		if v.Id == videoId {
			return true, nil
		}
	}
	return false, nil
}

// GetFavoriteCount get the favorite num of the video.
func (i *InteractionManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	return 0, nil
}
