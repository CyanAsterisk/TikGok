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
	VideoManager

	CommentPublisher
	FavoritePublisher

	CommentRedisManager
	FavoriteRedisManager

	CommentDao  *dao.Comment
	FavoriteDao *dao.Favorite
}

type VideoManager interface {
	GetPublishedVideoIdList(ctx context.Context, userId int64) ([]int64, error)
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
	GetFavoriteCountByVideoId(ctx context.Context, videoId int64) (int64, error)
	Like(ctx context.Context, userId, videoId, time int64) error
	Unlike(ctx context.Context, userId, videoId int64) error
	Check(ctx context.Context, userId, videoId int64) (bool, error)
	GetFavoriteVideoIdListByUserId(ctx context.Context, userId int64) ([]int64, error)
	GetFavoriteVideoCountByUserId(ctx context.Context, userId int64) (int64, error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(ctx context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	resp = new(interaction.DouyinFavoriteActionResponse)
	fav := &model.Favorite{
		UserId:     req.UserId,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
		CreateDate: time.Now().UnixNano(),
	}
	err = s.FavoritePublisher.Publish(ctx, fav)
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
			if err = s.FavoriteRedisManager.Like(ctx, fav.UserId, fav.VideoId, fav.CreateDate); err != nil {
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
}

// GetFavoriteVideoIdList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetFavoriteVideoIdList(ctx context.Context, req *interaction.DouyinGetFavoriteVideoIdListRequest) (resp *interaction.DouyinGetFavoriteVideoIdListResponse, err error) {
	resp = new(interaction.DouyinGetFavoriteVideoIdListResponse)
	resp.VideoIdList, err = s.FavoriteRedisManager.GetFavoriteVideoIdListByUserId(ctx, req.UserId)
	if err != nil {
		klog.Error("get videoIdList by redis err", err)
		resp.VideoIdList, err = s.FavoriteDao.GetFavoriteVideoIdListByUserId(req.UserId)
		if err != nil {
			klog.Error("get user favorite video list error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get user favorite video id list error"))
			return resp, nil
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
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
		CreateDate:  time.Now().UnixNano(),
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
}

// GetCommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetCommentList(ctx context.Context, req *interaction.DouyinGetCommentListRequest) (resp *interaction.DouyinGetCommentListResponse, err error) {
	resp = new(interaction.DouyinGetCommentListResponse)
	list, err := s.CommentRedisManager.GetCommentListByVideoId(ctx, req.VideoId)
	if err != nil {
		klog.Error("get comment list by redis err", err)
		list, err = s.CommentDao.GetCommentListByVideoId(req.VideoId)
		if err != nil {
			klog.Error("get comment list by video id error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get comment list error"))
			return resp, nil
		}
	}
	resp.CommentList = pkg.Comments(list)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetVideoInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetVideoInteractInfo(ctx context.Context, req *interaction.DouyinGetVideoInteractInfoRequest) (resp *interaction.DouyinGetVideoInteractInfoResponse, err error) {
	resp = new(interaction.DouyinGetVideoInteractInfoResponse)
	if resp.InteractInfo, err = s.getVideoInteractInfo(ctx, req.VideoId, req.ViewerId); err != nil {
		klog.Error("get video interact info err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetVideoInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) BatchGetVideoInteractInfo(ctx context.Context, req *interaction.DouyinBatchGetVideoInteractInfoRequest) (resp *interaction.DouyinBatchGetVideoInteractInfoResponse, err error) {
	resp = new(interaction.DouyinBatchGetVideoInteractInfoResponse)
	for _, vid := range req.VideoIdList {
		info, err := s.getVideoInteractInfo(ctx, vid, req.ViewerId)
		if err != nil {
			klog.Error("get video interact info err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
			return resp, nil
		}
		resp.InteractInfoList = append(resp.InteractInfoList, info)
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

func (s *InteractionServerImpl) getVideoInteractInfo(ctx context.Context, videoId, viewerId int64) (info *base.VideoInteractInfo, err error) {
	info = new(base.VideoInteractInfo)
	if info.CommentCount, err = s.CommentRedisManager.CommentCountByVideoId(ctx, videoId); err != nil {
		klog.Error("get comment count by redis err", err)
		if info.CommentCount, err = s.CommentDao.CommentCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	if info.FavoriteCount, err = s.FavoriteRedisManager.GetFavoriteCountByVideoId(ctx, videoId); err != nil {
		klog.Error("get favorite count by redis err", err)
		if info.FavoriteCount, err = s.FavoriteDao.GetFavoriteCountByVideoId(videoId); err != nil {
			return nil, err
		}
	}
	if info.IsFavorite, err = s.FavoriteRedisManager.Check(ctx, viewerId, videoId); err != nil {
		klog.Error("check like by redis err", err)
		fav, err := s.FavoriteDao.GetFavoriteInfo(viewerId, videoId)
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
}

// GetUserInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) GetUserInteractInfo(ctx context.Context, req *interaction.DouyinGetUserInteractInfoRequest) (resp *interaction.DouyinGetUserInteractInfoResponse, err error) {
	resp = new(interaction.DouyinGetUserInteractInfoResponse)
	resp.InteractInfo, err = s.getUserInteractInfo(ctx, req.UserId)
	if err != nil {
		klog.Error("get user interact info err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get user interact info err"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetUserInteractInfo implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) BatchGetUserInteractInfo(ctx context.Context, req *interaction.DouyinBatchGetUserInteractInfoRequest) (resp *interaction.DouyinBatchGetUserInteractInfoResponse, err error) {
	resp = new(interaction.DouyinBatchGetUserInteractInfoResponse)
	for _, uid := range req.UserIdList {
		info, err := s.getUserInteractInfo(ctx, uid)
		if err != nil {
			klog.Error("get user interact info err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
			return resp, nil
		}
		resp.InteractInfoList = append(resp.InteractInfoList, info)
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

func (s *InteractionServerImpl) getUserInteractInfo(ctx context.Context, userId int64) (info *base.UserInteractInfo, err error) {
	info = new(base.UserInteractInfo)
	videoIdList, err := s.VideoManager.GetPublishedVideoIdList(ctx, userId)
	if err != nil {
		return nil, err
	}
	info.WorkCount = int64(len(videoIdList))
	for _, vid := range videoIdList {
		count, err := s.FavoriteDao.GetFavoriteCountByVideoId(vid)
		if err != nil {
			return nil, err
		}
		info.TotalFavorited += count
	}

	info.FavoriteCount, err = s.FavoriteDao.GetFavoriteVideoCountByUserId(userId)
	if err != nil {
		return nil, err
	}
	return info, nil
}
