// Code generated by Kitex v0.4.4. DO NOT EDIT.

package videoservice

import (
	"context"
	video "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/video"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return videoServiceServiceInfo
}

var videoServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "VideoService"
	handlerType := (*video.VideoService)(nil)
	methods := map[string]kitex.MethodInfo{
		"Feed":                    kitex.NewMethodInfo(feedHandler, newVideoServiceFeedArgs, newVideoServiceFeedResult, false),
		"PublishVideo":            kitex.NewMethodInfo(publishVideoHandler, newVideoServicePublishVideoArgs, newVideoServicePublishVideoResult, false),
		"GetPublishedVideoList":   kitex.NewMethodInfo(getPublishedVideoListHandler, newVideoServiceGetPublishedVideoListArgs, newVideoServiceGetPublishedVideoListResult, false),
		"GetFavoriteVideoList":    kitex.NewMethodInfo(getFavoriteVideoListHandler, newVideoServiceGetFavoriteVideoListArgs, newVideoServiceGetFavoriteVideoListResult, false),
		"GetPublishedVideoIdList": kitex.NewMethodInfo(getPublishedVideoIdListHandler, newVideoServiceGetPublishedVideoIdListArgs, newVideoServiceGetPublishedVideoIdListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "video",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func feedHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceFeedArgs)
	realResult := result.(*video.VideoServiceFeedResult)
	success, err := handler.(video.VideoService).Feed(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceFeedArgs() interface{} {
	return video.NewVideoServiceFeedArgs()
}

func newVideoServiceFeedResult() interface{} {
	return video.NewVideoServiceFeedResult()
}

func publishVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServicePublishVideoArgs)
	realResult := result.(*video.VideoServicePublishVideoResult)
	success, err := handler.(video.VideoService).PublishVideo(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServicePublishVideoArgs() interface{} {
	return video.NewVideoServicePublishVideoArgs()
}

func newVideoServicePublishVideoResult() interface{} {
	return video.NewVideoServicePublishVideoResult()
}

func getPublishedVideoListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetPublishedVideoListArgs)
	realResult := result.(*video.VideoServiceGetPublishedVideoListResult)
	success, err := handler.(video.VideoService).GetPublishedVideoList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceGetPublishedVideoListArgs() interface{} {
	return video.NewVideoServiceGetPublishedVideoListArgs()
}

func newVideoServiceGetPublishedVideoListResult() interface{} {
	return video.NewVideoServiceGetPublishedVideoListResult()
}

func getFavoriteVideoListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetFavoriteVideoListArgs)
	realResult := result.(*video.VideoServiceGetFavoriteVideoListResult)
	success, err := handler.(video.VideoService).GetFavoriteVideoList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceGetFavoriteVideoListArgs() interface{} {
	return video.NewVideoServiceGetFavoriteVideoListArgs()
}

func newVideoServiceGetFavoriteVideoListResult() interface{} {
	return video.NewVideoServiceGetFavoriteVideoListResult()
}

func getPublishedVideoIdListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetPublishedVideoIdListArgs)
	realResult := result.(*video.VideoServiceGetPublishedVideoIdListResult)
	success, err := handler.(video.VideoService).GetPublishedVideoIdList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceGetPublishedVideoIdListArgs() interface{} {
	return video.NewVideoServiceGetPublishedVideoIdListArgs()
}

func newVideoServiceGetPublishedVideoIdListResult() interface{} {
	return video.NewVideoServiceGetPublishedVideoIdListResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Feed(ctx context.Context, req *video.DouyinFeedRequest) (r *video.DouyinFeedResponse, err error) {
	var _args video.VideoServiceFeedArgs
	_args.Req = req
	var _result video.VideoServiceFeedResult
	if err = p.c.Call(ctx, "Feed", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PublishVideo(ctx context.Context, req *video.DouyinPublishActionRequest) (r *video.DouyinPublishActionResponse, err error) {
	var _args video.VideoServicePublishVideoArgs
	_args.Req = req
	var _result video.VideoServicePublishVideoResult
	if err = p.c.Call(ctx, "PublishVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetPublishedVideoList(ctx context.Context, req *video.DouyinGetPublishedListRequest) (r *video.DouyinGetPublishedListResponse, err error) {
	var _args video.VideoServiceGetPublishedVideoListArgs
	_args.Req = req
	var _result video.VideoServiceGetPublishedVideoListResult
	if err = p.c.Call(ctx, "GetPublishedVideoList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFavoriteVideoList(ctx context.Context, req *video.DouyinGetFavoriteListRequest) (r *video.DouyinGetFavoriteListResponse, err error) {
	var _args video.VideoServiceGetFavoriteVideoListArgs
	_args.Req = req
	var _result video.VideoServiceGetFavoriteVideoListResult
	if err = p.c.Call(ctx, "GetFavoriteVideoList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetPublishedVideoIdList(ctx context.Context, req *video.DouyinGetPublishedVideoIdListRequest) (r *video.DouyinGetPublishedVideoIdListResponse, err error) {
	var _args video.VideoServiceGetPublishedVideoIdListArgs
	_args.Req = req
	var _result video.VideoServiceGetPublishedVideoIdListResult
	if err = p.c.Call(ctx, "GetPublishedVideoIdList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
