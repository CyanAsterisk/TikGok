// Code generated by Kitex v0.4.4. DO NOT EDIT.

package interactionserver

import (
	"context"
	interaction "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return interactionServerServiceInfo
}

var interactionServerServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "InteractionServer"
	handlerType := (*interaction.InteractionServer)(nil)
	methods := map[string]kitex.MethodInfo{
		"Favorite":     kitex.NewMethodInfo(favoriteHandler, newInteractionServerFavoriteArgs, newInteractionServerFavoriteResult, false),
		"FavoriteList": kitex.NewMethodInfo(favoriteListHandler, newInteractionServerFavoriteListArgs, newInteractionServerFavoriteListResult, false),
		"Comment":      kitex.NewMethodInfo(commentHandler, newInteractionServerCommentArgs, newInteractionServerCommentResult, false),
		"CommentList":  kitex.NewMethodInfo(commentListHandler, newInteractionServerCommentListArgs, newInteractionServerCommentListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "interaction",
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

func favoriteHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*interaction.InteractionServerFavoriteArgs)
	realResult := result.(*interaction.InteractionServerFavoriteResult)
	success, err := handler.(interaction.InteractionServer).Favorite(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newInteractionServerFavoriteArgs() interface{} {
	return interaction.NewInteractionServerFavoriteArgs()
}

func newInteractionServerFavoriteResult() interface{} {
	return interaction.NewInteractionServerFavoriteResult()
}

func favoriteListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*interaction.InteractionServerFavoriteListArgs)
	realResult := result.(*interaction.InteractionServerFavoriteListResult)
	success, err := handler.(interaction.InteractionServer).FavoriteList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newInteractionServerFavoriteListArgs() interface{} {
	return interaction.NewInteractionServerFavoriteListArgs()
}

func newInteractionServerFavoriteListResult() interface{} {
	return interaction.NewInteractionServerFavoriteListResult()
}

func commentHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*interaction.InteractionServerCommentArgs)
	realResult := result.(*interaction.InteractionServerCommentResult)
	success, err := handler.(interaction.InteractionServer).Comment(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newInteractionServerCommentArgs() interface{} {
	return interaction.NewInteractionServerCommentArgs()
}

func newInteractionServerCommentResult() interface{} {
	return interaction.NewInteractionServerCommentResult()
}

func commentListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*interaction.InteractionServerCommentListArgs)
	realResult := result.(*interaction.InteractionServerCommentListResult)
	success, err := handler.(interaction.InteractionServer).CommentList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newInteractionServerCommentListArgs() interface{} {
	return interaction.NewInteractionServerCommentListArgs()
}

func newInteractionServerCommentListResult() interface{} {
	return interaction.NewInteractionServerCommentListResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) Favorite(ctx context.Context, req *interaction.DouyinFavoriteActionRequest) (r *interaction.DouyinFavoriteActionResponse, err error) {
	var _args interaction.InteractionServerFavoriteArgs
	_args.Req = req
	var _result interaction.InteractionServerFavoriteResult
	if err = p.c.Call(ctx, "Favorite", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) FavoriteList(ctx context.Context, req *interaction.DouyinFavoriteListRequest) (r *interaction.DouyinFavoriteListResponse, err error) {
	var _args interaction.InteractionServerFavoriteListArgs
	_args.Req = req
	var _result interaction.InteractionServerFavoriteListResult
	if err = p.c.Call(ctx, "FavoriteList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) Comment(ctx context.Context, req *interaction.DouyinCommentActionRequest) (r *interaction.DouyinCommentActionResponse, err error) {
	var _args interaction.InteractionServerCommentArgs
	_args.Req = req
	var _result interaction.InteractionServerCommentResult
	if err = p.c.Call(ctx, "Comment", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CommentList(ctx context.Context, req *interaction.DouyinCommentListRequest) (r *interaction.DouyinCommentListResponse, err error) {
	var _args interaction.InteractionServerCommentListArgs
	_args.Req = req
	var _result interaction.InteractionServerCommentListResult
	if err = p.c.Call(ctx, "CommentList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}