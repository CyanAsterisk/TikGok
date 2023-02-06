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
	"time"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct {
	CommentPublisher
	CommentSubscriber
	FavoritePublisher
	FavoriteSubscriber

	CommentRedisManager
	FavoriteRedisManager
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

// CommentRedisManager defines the comment redis interface.
type CommentRedisManager interface {
	CommentCountByVideoId(videoId int64) (int64, error)
	CreateComment(comment *model.Comment) (*model.Comment, error)
	DeleteComment(id int64) error
	GetCommentListByVideoId(videoId int64) ([]*model.Comment, error)
}

// FavoriteRedisManager defines the favorite redis interface.
type FavoriteRedisManager interface {
	FavoriteCountByVideoId(videoId int64) (int64, error)
	CreateFavorite(fav *model.Favorite) error
	UpdateFavorite(userId, videoId int64, actionType int8) error
	GetFavoriteInfo(userId, videoId int64) (*model.Favorite, error)
	GetFavoriteVideoIdListByUserId(userId int64) ([]int64, error)
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
	faInfo, err := s.FavoriteRedisManager.GetFavoriteInfo(req.UserId, req.VideoId)
	if err == nil && faInfo == nil {
		if err = s.FavoriteRedisManager.CreateFavorite(&model.Favorite{
			UserId:     req.UserId,
			VideoId:    req.VideoId,
			ActionType: req.ActionType,
		}); err != nil {
			klog.Error("create favorite by redis err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("create favorite by redis err"))
			return resp, nil
		}
	}
	if err != nil {
		klog.Error("favorite by redis err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite by redis error"))
		return resp, nil
	}
	err = s.FavoriteRedisManager.UpdateFavorite(req.UserId, req.VideoId, req.ActionType)
	if err != nil {
		klog.Error("update favorite by redis error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("update favorite by redis error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil

	//
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

// GetFavoriteVideoIdList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetFavoriteVideoIdList(ctx context.Context, req *interaction.DouyinGetFavoriteVideoIdListRequest) (resp *interaction.DouyinGetFavoriteVideoIdListResponse, err error) {
	resp = new(interaction.DouyinGetFavoriteVideoIdListResponse)
	resp.VideoIdList, err = s.FavoriteRedisManager.GetFavoriteVideoIdListByUserId(req.UserId)
	if err != nil {
		klog.Error("get videoIdList by redis err", err)
		resp.VideoIdList, err = dao.GetFavoriteVideoIdListByUserId(req.UserId)
		if err != nil {
			klog.Error("get user favorite video list error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video id list error"))
			return resp, nil
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil

	//resp.VideoIdList, err = dao.GetFavoriteVideoIdListByUserId(req.UserId)
	//if err != nil {
	//	klog.Error("get user favorite video list error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video id list error"))
	//	return resp, nil
	//}
	//resp.BaseResp = tools.BuildBaseResp(nil)
	//return resp, nil
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
	var comment *model.Comment
	if req.ActionType == consts.ValidComment {
		comment, err = s.CommentRedisManager.CreateComment(&model.Comment{
			UserId:      req.UserId,
			VideoId:     req.VideoId,
			ActionType:  consts.ValidComment,
			CommentText: req.CommentText,
			CreateDate:  time.Now(),
		})
		if err != nil {
			klog.Errorf("create comment by redis error", err)
		}
	} else if req.ActionType == consts.InvalidComment {
		err = s.CommentRedisManager.DeleteComment(req.CommentId)
		if err != nil {
			klog.Errorf("delete comment from redis error")
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("delete comment from redis error"))
			return resp, nil
		}
	} else {
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("invalid ation type"))
		return resp, nil
	}

	resp.Comment = pkg.Comment(comment)
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

// GetCommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetCommentList(_ context.Context, req *interaction.DouyinGetCommentListRequest) (resp *interaction.DouyinGetCommentListResponse, err error) {
	resp = new(interaction.DouyinGetCommentListResponse)
	list, err := s.CommentRedisManager.GetCommentListByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment list by redis err", err)
		list, err = dao.GetCommentListByVideoId(req.VideoId)
		if err != nil {
			klog.Error("get comment list by video id error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
			return resp, nil
		}
	}
	resp.CommentList = pkg.Comments(list)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil

	//list, err := dao.GetCommentListByVideoId(req.VideoId)
	//if err != nil {
	//	klog.Error("get comment list by video id error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
	//	return resp, nil
	//}
	//resp.CommentList = pkg.Comments(list)
	//resp.BaseResp = tools.BuildBaseResp(nil)
	//return resp, nil
}

// GetInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetInteractInfo(_ context.Context, req *interaction.DouyinGetInteractInfoRequest) (resp *interaction.DouyinGetInteractInfoResponse, err error) {
	resp = new(interaction.DouyinGetInteractInfoResponse)
	if resp.InteractInfo, err = s.getInteractInfo(req.VideoId, req.ViewerId); err != nil {
		klog.Error("get interact info err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) BatchGetInteractInfo(ctx context.Context, req *interaction.DouyinBatchGetInteractInfoRequest) (resp *interaction.DouyinBatchGetInteractInfoResponse, err error) {
	resp = new(interaction.DouyinBatchGetInteractInfoResponse)
	for _, vid := range req.VideoIdList {
		info, err := s.getInteractInfo(vid, req.ViewerId)
		if err != nil {
			klog.Error("get interact info err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
			return resp, nil
		}
		resp.InteractInfoList = append(resp.InteractInfoList, info)
	}
	return resp, nil
}

func (s *InteractionServerImpl) getInteractInfo(videoId int64, viewerId int64) (info *base.InteractInfo, err error) {
	info = new(base.InteractInfo)
	if info.CommentCount, err = s.CommentRedisManager.CommentCountByVideoId(videoId); err != nil {
		klog.Error("get comment count by redis err", err)
		if info.CommentCount, err = dao.CommentCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	if info.FavoriteCount, err = s.FavoriteRedisManager.FavoriteCountByVideoId(videoId); err != nil {
		klog.Error("get favorite count by redis err", err)
		if info.FavoriteCount, err = dao.FavoriteCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	fav, err := s.FavoriteRedisManager.GetFavoriteInfo(viewerId, videoId)
	if err != nil {
		klog.Error("get favorite info by redis err", err)
		if fav, err = dao.GetFavoriteInfo(viewerId, videoId); err != nil {
			return nil, err
		}
	}
	if fav != nil && fav.ActionType == consts.IsLike {
		info.IsFavorite = true
	} else {
		info.IsFavorite = false
	}
	return info, nil

	//if info.CommentCount, err = dao.CommentCountByVideoId(videoId); err != nil {
	//	return nil, err
	//}
	//if info.FavoriteCount, err = dao.FavoriteCountByVideoId(videoId); err != nil {
	//	return nil, err
	//}
	//fav, err := dao.GetFavoriteInfo(viewerId, videoId)
	//if err != nil {
	//	return nil, err
	//}
	//if fav != nil && fav.ActionType == consts.IsLike {
	//	info.IsFavorite = true
	//} else {
	//	info.IsFavorite = false
	//}
	//return info, nil
}
