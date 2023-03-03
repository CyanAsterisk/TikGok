package main

import (
	"context"
	"sync"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	UserManager
	InteractionManager
	Publisher
	RedisManager
	Dao *dao.Video
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
	BatchGetVideoInteractInfo(ctx context.Context, videoIdList []int64, viewerId int64) ([]*base.VideoInteractInfo, error)
}

// Publisher defines the publisher video interface.
type Publisher interface {
	Publish(context.Context, *model.Video) error
}

// RedisManager defines the redis interface.
type RedisManager interface {
	CreateVideo(ctx context.Context, video *model.Video) error
	GetVideoListByLatestTime(ctx context.Context, latestTime int64) ([]*model.Video, error)
	GetVideoListByAuthorId(ctx context.Context, authorId int64) ([]*model.Video, error)
	GetVideoIdListByAuthorId(ctx context.Context, authorId int64) ([]int64, error)
	GetVideoByVideoId(ctx context.Context, videoId int64) (*model.Video, error)
	BatchGetVideoByVideoId(ctx context.Context, videoIdList []int64) ([]*model.Video, error)
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.DouyinFeedRequest) (resp *video.DouyinFeedResponse, err error) {
	resp = new(video.DouyinFeedResponse)
	if req.LatestTime <= 0 {
		req.LatestTime = time.Now().UnixNano()
	}
	vs, err := s.RedisManager.GetVideoListByLatestTime(ctx, req.LatestTime)
	if err != nil {
		klog.Error("get videos by latest time err", err)
		vs, err = s.Dao.GetVideoListByLatestTime(req.LatestTime)
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
		resp.NextTime = vs[len(vs)-1].CreateTime
	} else {
		resp.NextTime = time.Now().UnixNano()
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.DouyinPublishActionRequest) (resp *video.DouyinPublishActionResponse, err error) {
	resp = new(video.DouyinPublishActionResponse)
	sf, err := snowflake.NewNode(consts.VideoSnowflakeNode)
	if err != nil {
		klog.Errorf("create snowflake node err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr)
		return resp, nil
	}
	videoRecord := &model.Video{
		ID:         sf.Generate().Int64(),
		AuthorId:   req.UserId,
		PlayUrl:    req.PlayUrl,
		CoverUrl:   req.CoverUrl,
		Title:      req.Title,
		CreateTime: time.Now().UnixNano(),
	}
	err = s.Publish(ctx, videoRecord)
	if err != nil {
		klog.Errorf("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("publish video action error"))
		return resp, nil
	}
	if err = s.RedisManager.CreateVideo(ctx, videoRecord); err != nil {
		klog.Error("create video by redis err", err)
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetPublishedVideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishedVideoList(ctx context.Context, req *video.DouyinGetPublishedListRequest) (resp *video.DouyinGetPublishedListResponse, err error) {
	resp = new(video.DouyinGetPublishedListResponse)

	vs, err := s.RedisManager.GetVideoListByAuthorId(ctx, req.OwnerId)
	if err != nil {
		klog.Error("get published video by author id err", err)
		vs, err = s.Dao.GetVideoListByAuthorId(req.OwnerId)
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

// GetPublishedVideoIdList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishedVideoIdList(ctx context.Context, req *video.DouyinGetPublishedVideoIdListRequest) (resp *video.DouyinGetPublishedVideoIdListResponse, err error) {
	resp = new(video.DouyinGetPublishedVideoIdListResponse)
	if resp.VideoIdList, err = s.RedisManager.GetVideoIdListByAuthorId(ctx, req.UserId); err != nil {
		klog.Error("get published video id list by author id err", err)
		resp.VideoIdList, err = s.Dao.GetVideoIdListByAuthorId(req.UserId)
		if err != nil {
			klog.Error("get published video id list list err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.VideoServerErr.WithMessage("get published video id list err"))
			return resp, nil
		}
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

	videoList, err := s.RedisManager.BatchGetVideoByVideoId(ctx, idList)
	if err != nil {
		klog.Error("batch get video list by if from redis err", err)
		videoList, err = s.Dao.BatchGetVideoByVideoId(idList)
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
	videoIdList := make([]int64, 0)
	authorIdList := make([]int64, 0)
	for _, v := range videoList {
		videoIdList = append(videoIdList, v.ID)
		authorIdList = append(authorIdList, v.AuthorId)
	}
	var wg sync.WaitGroup
	wg.Add(2)

	var err error
	var authorList []*base.User
	go func() {
		defer wg.Done()
		authorList, err = s.UserManager.BatchGetUser(ctx, authorIdList, viewerId)
	}()

	var infoList []*base.VideoInteractInfo
	go func() {
		defer wg.Done()
		infoList, err = s.InteractionManager.BatchGetVideoInteractInfo(ctx, videoIdList, viewerId)
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return pkg.PackVideos(videoList, authorList, infoList), nil
}
