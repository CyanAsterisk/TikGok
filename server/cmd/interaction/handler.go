package main

import (
	"context"
	interaction "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
)

// InteractionServerImpl implements the last service interface defined in the IDL.
type InteractionServerImpl struct{}

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
	// TODO: Your code here...
	return
}

// CommentList implements the InteractionServerImpl interface.
func (s *InteractionServerImpl) CommentList(ctx context.Context, req *interaction.DouyinCommentListRequest) (resp *interaction.DouyinCommentListResponse, err error) {
	// TODO: Your code here...
	return
}
