package tools

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
	"github.com/cloudwego/kitex/pkg/klog"
)

type InteractionManager struct {
	InteractionService interactionserver.Client
}

// GetCommentCount  get one video comment.
func (i *InteractionManager) GetCommentCount(ctx context.Context, videoId int64) (int64, error) {
	res, err := i.InteractionService.CommentCount(ctx, &interaction.DouyinCommentCountRequest{VideoId: videoId})
	if err != nil {
		klog.Errorf("get comment count err", err.Error())
		return 0, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		klog.Errorf("get comment count err", res.BaseResp.StatusMsg)
		return 0, errno.InteractionServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Count, nil
}

// CheckFavorite check one favorite the video or not.
func (i *InteractionManager) CheckFavorite(ctx context.Context, userId int64, videoId int64) (bool, error) {
	res, err := i.InteractionService.CheckFavorite(ctx, &interaction.DouyinCheckFavoriteRequest{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		klog.Errorf("check favorite err", err.Error())
		return false, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		klog.Errorf("check favorite err", res.BaseResp.StatusMsg)
		return false, errno.InteractionServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Check, nil
}

// GetFavoriteCount get the favorite num of the video.
func (i *InteractionManager) GetFavoriteCount(ctx context.Context, videoId int64) (int64, error) {
	res, err := i.InteractionService.FavoriteCount(ctx, &interaction.DouyinFavoriteCountRequest{VideoId: videoId})
	if err != nil {
		klog.Errorf("get favorite count err", err.Error())
		return 0, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		klog.Errorf("get favorite count err", res.BaseResp.StatusMsg)
		return 0, errno.InteractionServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Count, nil
}
