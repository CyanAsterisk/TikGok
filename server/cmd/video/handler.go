package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	video "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(_ context.Context, req *video.DouyinFeedRequest) (resp *video.DouyinFeedResponse, err error) {
	resp = new(video.DouyinFeedResponse)
	if req.LatestTime <= 0 {
		req.LatestTime = time.Now().UnixNano() / 1e6
	}
	vs, err := dao.GetVideosByLatestTime(req.LatestTime)
	if err != nil {
		klog.Error("get videos by latest time err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("get videos error"))
	}

	if resp.VideoList, err = tools.Videos(vs); err != nil {
		klog.Errorf("convert err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr)
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
func (s *VideoServiceImpl) VideoList(_ context.Context, req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error) {
	resp = new(video.DouyinPublishListResponse)
	vs, err := dao.GetVideosByUserId(req.UserId)
	if err != nil {
		klog.Error("get published video list err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.VideoServerErr.WithMessage("get published video list err"))
		return
	}
	if resp.VideoList, err = tools.Videos(vs); err != nil {
		klog.Errorf("convert err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr)
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
	if resp.Video, err = tools.Video(v); err != nil {
		klog.Errorf("convert err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr)
		return
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return
}
