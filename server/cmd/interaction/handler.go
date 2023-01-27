package main

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/pack"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct {
	CommentManager
	VideoManager
}

// CommentManager manage comment status.
type CommentManager interface {
	GetResp(req *interaction.DouyinCommentActionRequest) (comment *model.Comment, err error)
}

// VideoManager defines the Anti Corruption Layer
// for get video logic.
type VideoManager interface {
	GetVideos([]int64) ([]*base.Video, error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(_ context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	resp = new(interaction.DouyinFavoriteActionResponse)
	_, err = dao.GetFavoriteInfo(req.UserId, req.VideoId)
	if err == gorm.ErrRecordNotFound {
		err = dao.CreateFavorite(&model.Favorite{
			UserId:     req.UserId,
			VideoId:    req.VideoId,
			ActionType: consts.IsLike,
		})
		if err != nil {
			klog.Error("favorite error", err)
			resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
			return resp, nil
		}
		resp.BaseResp = pack.BuildBaseResp(nil)
		return resp, nil
	}
	if err != nil {
		klog.Error("favorite error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
		return resp, nil
	}
	err = dao.UpdateFavorite(req.UserId, req.VideoId, req.ActionType)
	if err != nil {
		klog.Error("favorite error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}

// FavoriteList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteList(_ context.Context, req *interaction.DouyinFavoriteListRequest) (resp *interaction.DouyinFavoriteListResponse, err error) {
	resp = new(interaction.DouyinFavoriteListResponse)
	list, err := dao.GetFavoriteVideoIdListByUserId(req.UserId)
	if err != nil {
		klog.Error("get user favorite video list error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	videos, err := s.VideoManager.GetVideos(list)
	if err != nil {
		klog.Error("get videos by video manager error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	resp.VideoList = videos
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}

// Comment implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Comment(_ context.Context, req *interaction.DouyinCommentActionRequest) (resp *interaction.DouyinCommentActionResponse, err error) {
	resp = new(interaction.DouyinCommentActionResponse)
	cmt, err := s.GetResp(req)
	if err != nil {
		klog.Error("comment uses get response error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("comment error"))
		return resp, nil
	}
	resp.Comment = tools.Comment(cmt)
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}

// CommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentList(_ context.Context, req *interaction.DouyinCommentListRequest) (resp *interaction.DouyinCommentListResponse, err error) {
	resp = new(interaction.DouyinCommentListResponse)
	list, err := dao.GetCommentListByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment list by video id error", err)
		resp.BaseResp = pack.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
		return resp, nil
	}
	resp.CommentList = tools.Comments(list)
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}
