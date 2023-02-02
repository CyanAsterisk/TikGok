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
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
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
	GetVideos(ctx context.Context, list []int64, viewerId int64) ([]*base.Video, error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(_ context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	resp = new(interaction.DouyinFavoriteActionResponse)
	faInfo, err := dao.GetFavoriteInfo(req.UserId, req.VideoId)
	if err == nil && faInfo == nil {
		err = dao.CreateFavorite(&model.Favorite{
			UserId:     req.UserId,
			VideoId:    req.VideoId,
			ActionType: consts.IsLike,
		})
		if err != nil {
			klog.Error("favorite error", err)
			resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
			return resp, nil
		}
		resp.BaseResp = sTools.BuildBaseResp(nil)
		return resp, nil
	}
	if err != nil {
		klog.Error("favorite error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
		return resp, nil
	}
	err = dao.UpdateFavorite(req.UserId, req.VideoId, req.ActionType)
	if err != nil {
		klog.Error("favorite error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
		return resp, nil
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FavoriteList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteList(ctx context.Context, req *interaction.DouyinFavoriteListRequest) (resp *interaction.DouyinFavoriteListResponse, err error) {
	resp = new(interaction.DouyinFavoriteListResponse)
	list, err := dao.GetFavoriteVideoIdListByUserId(req.OwnerId)
	if err != nil {
		klog.Error("get user favorite video list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	videos, err := s.VideoManager.GetVideos(ctx, list, req.ViewerId)
	if err != nil {
		klog.Error("get videos by video manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCVideoErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	resp.VideoList = videos
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FavoriteCount implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteCount(_ context.Context, req *interaction.DouyinFavoriteCountRequest) (resp *interaction.DouyinFavoriteCountResponse, err error) {
	resp = new(interaction.DouyinFavoriteCountResponse)
	count, err := dao.FavoriteCountByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get favorite count error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get favorite count error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// CheckFavorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CheckFavorite(ctx context.Context, req *interaction.DouyinCheckFavoriteRequest) (resp *interaction.DouyinCheckFavoriteResponse, err error) {
	resp = new(interaction.DouyinCheckFavoriteResponse)
	info, err := dao.GetFavoriteInfo(req.UserId, req.VideoId)
	if err != nil {
		klog.Error("check favorite error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("check favorite error"))
		return resp, nil
	}
	if info.ActionType == consts.IsLike {
		resp.Check = true
	} else {
		resp.Check = false
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// Comment implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Comment(_ context.Context, req *interaction.DouyinCommentActionRequest) (resp *interaction.DouyinCommentActionResponse, err error) {
	resp = new(interaction.DouyinCommentActionResponse)
	cmt, err := s.GetResp(req)
	if err != nil {
		klog.Error("comment uses get response error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("comment error"))
		return resp, nil
	}
	resp.Comment = tools.Comment(cmt)
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// CommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentList(_ context.Context, req *interaction.DouyinCommentListRequest) (resp *interaction.DouyinCommentListResponse, err error) {
	resp = new(interaction.DouyinCommentListResponse)
	list, err := dao.GetCommentListByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment list by video id error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
		return resp, nil
	}
	resp.CommentList = tools.Comments(list)
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// CommentCount implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentCount(_ context.Context, req *interaction.DouyinCommentCountRequest) (resp *interaction.DouyinCommentCountResponse, err error) {
	resp = new(interaction.DouyinCommentCountResponse)
	count, err := dao.CommentCountByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment count error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment count error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}
