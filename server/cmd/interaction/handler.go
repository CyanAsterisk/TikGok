package main

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/pack"
	interaction "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/codes"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct {
	CommentManager
}

// CommentManager manager comment status.
type CommentManager interface {
	GetResp(req *interaction.DouyinCommentActionRequest) (comment *model.Comment, err error)
}

// Favorite implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Favorite(ctx context.Context, req *interaction.DouyinFavoriteActionRequest) (resp *interaction.DouyinFavoriteActionResponse, err error) {
	// TODO: Your code here...
	return
}

// FavoriteList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) FavoriteList(ctx context.Context, req *interaction.DouyinFavoriteListRequest) (resp *interaction.DouyinFavoriteListResponse, err error) {
	// TODO: Your code here...
	return
}

// Comment implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) Comment(ctx context.Context, req *interaction.DouyinCommentActionRequest) (resp *interaction.DouyinCommentActionResponse, err error) {
	resp = new(interaction.DouyinCommentActionResponse)
	cmt, err := s.GetResp(req)
	if err != nil {
		klog.Error("comment uses get response error", err)
		resp.BaseResp = pack.BuildBaseResp(status.Err(codes.Internal, "comment error"))
		return resp, nil
	}
	resp.Comment = pack.Comment(cmt)
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}

// CommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentList(ctx context.Context, req *interaction.DouyinCommentListRequest) (resp *interaction.DouyinCommentListResponse, err error) {
	resp = new(interaction.DouyinCommentListResponse)
	list, err := dao.GetCommentListByVideoId(req.VideoId)
	if err != nil {
		klog.Error("get comment list by video id error", err)
		resp.BaseResp = pack.BuildBaseResp(status.Err(codes.Internal, "get comment list error"))
		return resp, nil
	}
	resp.CommentList = pack.Comments(list)
	resp.BaseResp = pack.BuildBaseResp(nil)
	return resp, nil
}
