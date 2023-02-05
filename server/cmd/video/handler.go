package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	UserManager
	InteractionManager
	Publisher
	Subscriber
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
	GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error)
	BatchCheckFavorite(ctx context.Context, userId int64, videoIdList []int64) ([]bool, error)
	BatchGetCommentCount(ctx context.Context, videoIdList []int64) ([]int64, error)
	BatchGetFavoriteCount(ctx context.Context, videoId []int64) ([]int64, error)
}

// Publisher defines the publisher video interface.
type Publisher interface {
	Publish(context.Context, *video.DouyinPublishActionRequest) error
}

// Subscriber defines a video publish subscriber.
type Subscriber interface {
	Subscribe(context.Context) (ch chan *video.DouyinPublishActionRequest, cleanUp func(), err error)
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
		resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("get videos error"))
		return
	}

	resp.VideoList, err = s.packVideos(ctx, vs, req.ViewerId)
	if err != nil {
		klog.Error("pack videos err", err.Error())
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr.WithMessage("pack videos err"))
		return
	}
	if len(vs) > 0 {
		resp.NextTime = vs[len(vs)-1].CreatedAt.UnixNano() / 1e6
	} else {
		resp.NextTime = time.Now().UnixNano() / 1e6
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.DouyinPublishActionRequest) (resp *video.DouyinPublishActionResponse, err error) {
	resp = new(video.DouyinPublishActionResponse)
	err = s.Publish(ctx, req)
	if err != nil {
		klog.Errorf("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("publish video action error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil

	//vid := model.Video{
	//	Model: gorm.Model{
	//		CreatedAt: time.Now(),
	//		UpdatedAt: time.Now(),
	//	},
	//	Uid:      req.UserId,
	//	PlayUrl:  req.PlayUrl,
	//	CoverUrl: req.CoverUrl,
	//	Title:    req.Title,
	//}
	//err = dao.CreateVideo(&vid)
	//if err != nil {
	//	klog.Errorf("create video err", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("create video err"))
	//	return
	//}
	//
	//resp.BaseResp = tools.BuildBaseResp(nil)
	//return
}

// GetPublishedVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishedVideoList(ctx context.Context, req *video.DouyinGetPublishedListRequest) (resp *video.DouyinGetPublishedListResponse, err error) {
	resp = new(video.DouyinGetPublishedListResponse)
	vs, err := dao.GetVideosByUserId(req.OwnerId)
	if err != nil {
		klog.Error("get published video list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("get published video list err"))
		return
	}
	resp.VideoList, err = s.packVideos(ctx, vs, req.ViewerId)
	if err != nil {
		klog.Error("pack videos err", err.Error())
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr.WithMessage("pack videos err"))
		return
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// GetFavoriteVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetFavoriteVideoList(ctx context.Context, req *video.DouyinGetFavoriteListRequest) (resp *video.DouyinGetFavoriteListResponse, err error) {
	resp = new(video.DouyinGetFavoriteListResponse)

	idList, err := s.InteractionManager.GetFavoriteVideoIdList(ctx, req.OwnerId)
	if err != nil {
		klog.Error("get favorite video id list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCInteractionErr)
		return resp, nil
	}

	videos := make([]*model.Video, 0)
	for _, vid := range idList {
		v, err := dao.GetVideoByVideoId(vid)
		if err != nil {
			klog.Error("get video by id err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
			return resp, nil
		}
		videos = append(videos, v)
	}
	resp.VideoList, err = s.packVideos(ctx, videos, req.ViewerId)
	if err != nil {
		klog.Error("pack video err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

func (s *VideoServiceImpl) packVideos(ctx context.Context, mvs []*model.Video, uid int64) (bvs []*base.Video, err error) {

	return nil, nil
}
