package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct {
	CommentManager
	VideoManager

	CommentPublisher
	CommentSubscriber
	FavoritePublisher
	FavoriteSubscriber
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

// CommentPublisher defines the comment action publisher interface.
type CommentPublisher interface {
	Publish(context.Context, *interaction.DouyinCommentActionRequest) error
}

// CommentSubscriber defines a comment action subscriber.
type CommentSubscriber interface {
	Subscribe(context.Context) (ch chan *interaction.DouyinCommentActionRequest, cleanUp func(), err error)
}

// FavoritePublisher defines the favorite action publisher interface.
type FavoritePublisher interface {
	Publish(context.Context, *interaction.DouyinFavoriteActionRequest) error
}

// FavoriteSubscriber  defines a favorite action subscriber interface.
type FavoriteSubscriber interface {
	Subscribe(context.Context) (ch chan *interaction.DouyinFavoriteActionRequest, cleanUp func(), err error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(ctx context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	resp = new(interaction.DouyinFavoriteActionResponse)
	err = s.FavoritePublisher.Publish(ctx, req)
	if err != nil {
		klog.Error("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite action error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil

	//faInfo, err := dao.GetFavoriteInfo(req.UserId, req.VideoId)
	//if err == nil && faInfo == nil {
	//	err = dao.CreateFavorite(&model.Favorite{
	//		UserId:     req.UserId,
	//		VideoId:    req.VideoId,
	//		ActionType: consts.IsLike,
	//	})
	//	if err != nil {
	//		klog.Error("favorite error", err)
	//		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
	//		return resp, nil
	//	}
	//	resp.BaseResp = tools.BuildBaseResp(nil)
	//	return resp, nil
	//}
	//if err != nil {
	//	klog.Error("favorite error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
	//	return resp, nil
	//}
	//err = dao.UpdateFavorite(req.UserId, req.VideoId, req.ActionType)
	//if err != nil {
	//	klog.Error("favorite error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite error"))
	//	return resp, nil
	//}
	//resp.BaseResp = tools.BuildBaseResp(nil)
	//return resp, nil
}

// FavoriteList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteList(ctx context.Context, req *interaction.DouyinFavoriteListRequest) (resp *interaction.DouyinFavoriteListResponse, err error) {
	resp = new(interaction.DouyinFavoriteListResponse)
	list, err := dao.GetFavoriteVideoIdListByUserId(req.OwnerId)
	if err != nil {
		klog.Error("get user favorite video list error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	videos, err := s.VideoManager.GetVideos(ctx, list, req.ViewerId)
	if err != nil {
		klog.Error("get videos by video manager error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCVideoErr.WithMessage("get user favorite video list error"))
		return resp, nil
	}
	resp.VideoList = videos
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// FavoriteCount implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteCount(_ context.Context, req *interaction.DouyinFavoriteCountRequest) (resp *interaction.DouyinFavoriteCountResponse, err error) {
	resp = new(interaction.DouyinFavoriteCountResponse)
	count, err := dao.FavoriteCountByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get favorite count error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get favorite count error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// CheckFavorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CheckFavorite(_ context.Context, req *interaction.DouyinCheckFavoriteRequest) (resp *interaction.DouyinCheckFavoriteResponse, err error) {
	resp = new(interaction.DouyinCheckFavoriteResponse)
	info, err := dao.GetFavoriteInfo(req.UserId, req.VideoId)
	if err != nil {
		klog.Error("check favorite error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("check favorite error"))
		return resp, nil
	}
	if info == nil {
		resp.Check = false
	} else {
		if info.ActionType == consts.IsLike {
			resp.Check = true
		} else {
			resp.Check = false
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// Comment implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Comment(ctx context.Context, req *interaction.DouyinCommentActionRequest) (resp *interaction.DouyinCommentActionResponse, err error) {
	resp = new(interaction.DouyinCommentActionResponse)
	err = s.CommentPublisher.Publish(ctx, req)
	if err != nil {
		klog.Error("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("comment action error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
	//cmt, err := s.GetResp(req)
	//if err != nil {
	//	klog.Error("comment uses get response error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("comment error"))
	//	return resp, nil
	//}
	//resp.Comment = pkg.Comment(cmt)
	//resp.BaseResp = tools.BuildBaseResp(nil)
	//return resp, nil
}

// CommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentList(_ context.Context, req *interaction.DouyinCommentListRequest) (resp *interaction.DouyinCommentListResponse, err error) {
	resp = new(interaction.DouyinCommentListResponse)
	list, err := dao.GetCommentListByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment list by video id error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
		return resp, nil
	}
	resp.CommentList = pkg.Comments(list)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// CommentCount implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentCount(_ context.Context, req *interaction.DouyinCommentCountRequest) (resp *interaction.DouyinCommentCountResponse, err error) {
	resp = new(interaction.DouyinCommentCountResponse)
	count, err := dao.CommentCountByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment count error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment count error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
