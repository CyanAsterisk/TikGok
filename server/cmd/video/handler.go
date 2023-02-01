package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	video "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	UserManager
	InteractionManager
}

// UserManager defines the Anti Corruption Layer
// for get user logic.
type UserManager interface {
	GetUsers(ctx context.Context, list []int64, uid int64) ([]*base.User, error)
	GetUser(ctx context.Context, UserId, toUserId int64) (*base.User, error)
}

// InteractionManager defines the Anti Corruption Layer
// for get interaction logic.
type InteractionManager interface {
	GetCommentCount(ctx context.Context, videoId int64) (int64, error)
	CheckFavorite(ctx context.Context, userId int64, videoId int64) (bool, error)
	GetFavoriteCount(ctx context.Context, videoId int64) (int64, error)
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.DouyinFeedRequest) (resp *video.DouyinFeedResponse, err error) {
	resp = new(video.DouyinFeedResponse)
	if req.LatestTime <= 0 {
		req.LatestTime = time.Now().UnixNano() / 1e6
	}
	vs, err := dao.GetVideosByLatestTime(req.LatestTime)
	if err != nil {
		klog.Error("get videos by latest time err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("get videos error"))
		return
	}

	resp.VideoList, err = s.packVideos(ctx, vs, req.UserId)
	if err != nil {
		klog.Error("pack videos err", err.Error())
		resp.BaseResp = sTools.BuildBaseResp(errno.ServiceErr.WithMessage("pack videos err"))
		return
	}
	if len(vs) > 0 {
		resp.NextTime = vs[len(vs)-1].CreatedAt.UnixNano() / 1e6
	} else {
		resp.NextTime = time.Now().UnixNano() / 1e6
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(_ context.Context, req *video.DouyinPublishActionRequest) (resp *video.DouyinPublishActionResponse, err error) {
	resp = new(video.DouyinPublishActionResponse)
	vid := model.Video{
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Uid:      req.UserId,
		PlayUrl:  req.PlayUrl,
		CoverUrl: req.CoverUrl,
		Title:    req.Title,
	}
	err = dao.CreateVideo(&vid)
	if err != nil {
		klog.Errorf("create video err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("create video err"))
		return
	}

	resp.BaseResp = sTools.BuildBaseResp(nil)
	return
}

// VideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoList(ctx context.Context, req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error) {
	resp = new(video.DouyinPublishListResponse)
	vs, err := dao.GetVideosByUserId(req.UserId)
	if err != nil {
		klog.Error("get published video list err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("get published video list err"))
		return
	}
	resp.VideoList, err = s.packVideos(ctx, vs, req.UserId)
	if err != nil {
		klog.Error("pack videos err", err.Error())
		resp.BaseResp = sTools.BuildBaseResp(errno.ServiceErr.WithMessage("pack videos err"))
		return
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return
}

// GetVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideo(ctx context.Context, req *video.DouyinGetVideoRequest) (resp *video.DouyinGetVideoResponse, err error) {
	resp = new(video.DouyinGetVideoResponse)
	v, err := dao.GetVideoByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get video err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("get video err"))
		return
	}
	resp.Video, err = s.packVideo(ctx, v, 0)
	if err != nil {
		klog.Error("pack video err", err.Error())
		resp.BaseResp = sTools.BuildBaseResp(errno.ServiceErr.WithMessage("pack video err"))
		return
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return
}

func (s *VideoServiceImpl) packVideo(ctx context.Context, mv *model.Video, uid int64) (bv *base.Video, err error) {
	bv = tools.Video(mv)
	if bv.Author, err = s.GetUser(ctx, uid, mv.Uid); err != nil {
		return nil, err
	}
	if bv.IsFavorite, err = s.CheckFavorite(ctx, uid, bv.Id); err != nil {
		return nil, err
	}
	if bv.CommentCount, err = s.GetCommentCount(ctx, bv.Id); err != nil {
		return nil, err
	}
	if bv.FavoriteCount, err = s.GetFavoriteCount(ctx, bv.Id); err != nil {
		return nil, err
	}
	return
}

func (s *VideoServiceImpl) packVideos(ctx context.Context, mvs []*model.Video, uid int64) (bvs []*base.Video, err error) {
	bvs = make([]*base.Video, 0)
	for _, mv := range mvs {
		bv, err := s.packVideo(ctx, mv, uid)
		if err != nil {
			return nil, err
		}
		bvs = append(bvs, bv)
	}
	return bvs, nil
}
