package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/pkg"
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
	RedisManager
}

// UserManager defines the Anti Corruption Layer
// for get user logic.
type UserManager interface {
	BatchGetUser(ctx context.Context, list []int64, uid int64) ([]*base.User, error)
	GetUser(ctx context.Context, UserId, toUserId int64) (*base.User, error)
}

// InteractionManager defines the Anti Corruption Layer
// for get interaction logic.
type InteractionManager interface {
	GetFavoriteVideoIdList(ctx context.Context, userId int64) ([]int64, error)
	BatchGetInteractInfo(ctx context.Context, videoIdList []int64, viewerId int64) ([]*base.InteractInfo, error)
}

// Publisher defines the publisher video interface.
type Publisher interface {
	Publish(context.Context, *video.DouyinPublishActionRequest) error
}

// Subscriber defines a video publish subscriber.
type Subscriber interface {
	Subscribe(context.Context) (ch chan *video.DouyinPublishActionRequest, cleanUp func(), err error)
}

// RedisManager defines the redis interface.
type RedisManager interface {
	CreateVideo(video *model.Video) error
	GetVideosByLatestTime(latestTime int64) ([]*model.Video, error)
	GetVideosByUserId(uid int64) ([]*model.Video, error)
	GetVideoByVideoId(vid int64) (*model.Video, error)
	BatchGetVideoByVideoId(vidList []int64) ([]*model.Video, error)
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.DouyinFeedRequest) (resp *video.DouyinFeedResponse, err error) {
	resp = new(video.DouyinFeedResponse)
	if req.LatestTime <= 0 {
		req.LatestTime = time.Now().UnixNano() / 1e6
	}
	vs, err := s.RedisManager.GetVideosByLatestTime(req.LatestTime)
	if err != nil {
		klog.Error("get videos by latest time err", err)
		vs, err = dao.GetVideosByLatestTime(req.LatestTime)
		if err != nil {
			klog.Error("get videos by latest time err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("get videos error"))
			return resp, nil
		}
	}
	resp.VideoList, err = s.fillVideoList(ctx, vs, req.ViewerId)
	if err != nil {
		klog.Error("fill video list err", err.Error())
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr.WithMessage("fill video list err"))
		return resp, nil
	}
	if len(vs) > 0 {
		resp.NextTime = vs[len(vs)-1].UpdatedAt.UnixNano() / 1e6
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
	v := &model.Video{
		Uid:       req.UserId,
		PlayUrl:   req.PlayUrl,
		CoverUrl:  req.CoverUrl,
		Title:     req.Title,
		UpdatedAt: time.Now(),
	}

	if err = s.RedisManager.CreateVideo(v); err != nil {
		err = dao.CreateVideo(v)
		if err != nil {
			klog.Error("create video by redis err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("create video err"))
			return resp, nil
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetPublishedVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishedVideoList(ctx context.Context, req *video.DouyinGetPublishedListRequest) (resp *video.DouyinGetPublishedListResponse, err error) {
	resp = new(video.DouyinGetPublishedListResponse)

	vs, err := s.RedisManager.GetVideosByUserId(req.OwnerId)
	if err != nil {
		klog.Error("get published video by author id err", err)
		vs, err = dao.GetVideosByUserId(req.OwnerId)
		if err != nil {
			klog.Error("get published video list err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("get published video list err"))
			return resp, nil
		}
	}
	resp.VideoList, err = s.fillVideoList(ctx, vs, req.ViewerId)
	if err != nil {
		klog.Error("fill video list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr.WithMessage("fill video list err"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
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

	videoList, err := s.RedisManager.BatchGetVideoByVideoId(idList)
	if err != nil {
		klog.Error("batch get video list by if from redis err", err)
		videoList, err = dao.BatchGetVideoByVideoId(idList)
		if err != nil {
			klog.Error("batch get video list by video id list err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr)
			return resp, nil
		}
	}

	resp.VideoList, err = s.fillVideoList(ctx, videoList, req.ViewerId)
	if err != nil {
		klog.Error("fill video list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

func (s *VideoServiceImpl) fillVideoList(ctx context.Context, videoList []*model.Video, viewerId int64) ([]*base.Video, error) {
	if videoList == nil {
		return nil, nil
	}
	videoIdList := make([]int64, len(videoList))
	authorIdList := make([]int64, len(videoList))
	for _, v := range videoList {
		videoIdList = append(videoIdList, v.ID)
		authorIdList = append(authorIdList, v.Uid)
	}
	authorList, err := s.UserManager.BatchGetUser(ctx, authorIdList, viewerId)
	if err != nil {
		return nil, err
	}
	InfoList, err := s.InteractionManager.BatchGetInteractInfo(ctx, authorIdList, viewerId)
	if err != nil {
		return nil, err
	}
	return pkg.PackVideos(videoList, authorList, InfoList), nil
}
