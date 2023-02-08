package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct {
	CommentPublisher
	FavoritePublisher

	CommentRedisManager
	FavoriteRedisManager
}

// CommentPublisher defines the comment action publisher interface.
type CommentPublisher interface {
	Publish(context.Context, *model.Comment) error
}

// FavoritePublisher defines the favorite action publisher interface.
type FavoritePublisher interface {
	Publish(context.Context, *model.Favorite) error
}

// CommentRedisManager defines the comment redis interface.
type CommentRedisManager interface {
	CommentCountByVideoId(ctx context.Context, videoId int64) (int64, error)
	CreateComment(ctx context.Context, comment *model.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	GetCommentListByVideoId(ctx context.Context, videoId int64) ([]*model.Comment, error)
}

// FavoriteRedisManager defines the favorite redis interface.
type FavoriteRedisManager interface {
	FavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error)
	Like(ctx context.Context, userId int64, videoId int64) error
	Unlike(ctx context.Context, userId int64, videoId int64) error
	Check(ctx context.Context, userId int64, videoId int64) (bool, error)
	GetFavoriteVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(ctx context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	resp = new(interaction.DouyinFavoriteActionResponse)
	err = s.FavoritePublisher.Publish(ctx, &model.Favorite{
		UserId:     req.UserId,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	if err != nil {
		klog.Error("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("favorite action error"))
		return resp, nil
	}
	liked, err := s.FavoriteRedisManager.Check(ctx, req.UserId, req.VideoId)
	if err != nil {
		klog.Error("check like by redis err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("check like by redis error"))
		return resp, nil
	}
	if req.ActionType == consts.Like {
		if !liked {
			if err = s.FavoriteRedisManager.Like(ctx, req.UserId, req.VideoId); err != nil {
				klog.Error("like by redis error", err)
				resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("like by redis error"))
				return resp, nil
			}
		}
	} else if req.ActionType == consts.UnLike {
		if liked {
			if err = s.FavoriteRedisManager.Unlike(ctx, req.UserId, req.VideoId); err != nil {
				klog.Error("unlike by redis err", err)
				resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("unlike by redis error"))
				return resp, nil
			}
		}
	} else {
		resp.BaseResp = tools.BuildBaseResp(errno.ParamsEr.WithMessage("invalid action type"))
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
	resp.VideoIdList, err = s.FavoriteRedisManager.GetFavoriteVideoIdListByUserId(ctx, req.UserId)
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
	comment := &model.Comment{
		ID:          req.CommentId,
		UserId:      req.UserId,
		VideoId:     req.VideoId,
		ActionType:  req.ActionType,
		CommentText: req.CommentText,
		CreateDate:  time.Now(),
	}
	if req.ActionType == consts.ValidComment {
		sf, err := snowflake.NewNode(consts.CommentSnowflakeNode)
		if err != nil {
			klog.Errorf("generate id failed: %s", err.Error())
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("generate id failed"))
		}
		comment.ID = sf.Generate().Int64()
	}
	err = s.CommentPublisher.Publish(ctx, comment)
	if err != nil {
		klog.Error("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("comment action error"))
		return resp, nil
	}
	if req.ActionType == consts.ValidComment {
		err = s.CommentRedisManager.CreateComment(ctx, comment)
		if err != nil {
			klog.Errorf("create comment by redis error", err)
		}
	} else if req.ActionType == consts.InvalidComment {
		err = s.CommentRedisManager.DeleteComment(ctx, req.CommentId)
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
func (s *InteractionServerImpl) GetCommentList(ctx context.Context, req *interaction.DouyinGetCommentListRequest) (resp *interaction.DouyinGetCommentListResponse, err error) {
	resp = new(interaction.DouyinGetCommentListResponse)
	list, err := s.CommentRedisManager.GetCommentListByVideoId(ctx, req.VideoId)
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
func (s *InteractionServerImpl) GetInteractInfo(ctx context.Context, req *interaction.DouyinGetInteractInfoRequest) (resp *interaction.DouyinGetInteractInfoResponse, err error) {
	resp = new(interaction.DouyinGetInteractInfoResponse)
	if resp.InteractInfo, err = s.getInteractInfo(ctx, req.VideoId, req.ViewerId); err != nil {
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
		info, err := s.getInteractInfo(ctx, vid, req.ViewerId)
		if err != nil {
			klog.Error("get interact info err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
			return resp, nil
		}
		resp.InteractInfoList = append(resp.InteractInfoList, info)
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

func (s *InteractionServerImpl) getInteractInfo(ctx context.Context, videoId int64, viewerId int64) (info *base.InteractInfo, err error) {
	info = new(base.InteractInfo)
	if info.CommentCount, err = s.CommentRedisManager.CommentCountByVideoId(ctx, videoId); err != nil {
		klog.Error("get comment count by redis err", err)
		if info.CommentCount, err = dao.CommentCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	if info.FavoriteCount, err = s.FavoriteRedisManager.FavoriteCountByVideoId(ctx, videoId); err != nil {
		klog.Error("get favorite count by redis err", err)
		if info.FavoriteCount, err = dao.FavoriteCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	if info.IsFavorite, err = s.FavoriteRedisManager.Check(ctx, viewerId, videoId); err != nil {
		klog.Error("check like by redis err", err)
		fav, err := dao.GetFavoriteInfo(viewerId, videoId)
		if err != nil {
			klog.Error("get favorite info err", err)
			return nil, err
		}
		if fav == nil {
			info.IsFavorite = false
		} else {
			if fav.ActionType == consts.IsLike {
				info.IsFavorite = true
			} else {
				info.IsFavorite = false
			}
		}
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
