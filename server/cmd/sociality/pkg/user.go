package pkg

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
)

type UserManager struct {
	UserService userservice.Client
	ChatService chatservice.Client
}

// GetUsers gets users info by list.
func (m *UserManager) GetUsers(ctx context.Context, list []int64, viewerId int64) ([]*base.User, error) {
	resp, err := m.UserService.BatchGetUserInfo(ctx, &user.DouyinBatchGetUserRequest{
		ViewerId: viewerId,
		OwnerIds: list,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.UserServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Users, nil
}

func (m *UserManager) GetFriendUsers(ctx context.Context, list []int64, viewerId int64) ([]*base.FriendUser, error) {

	resp, err := m.UserService.BatchGetUserInfo(ctx, &user.DouyinBatchGetUserRequest{
		ViewerId: viewerId,
		OwnerIds: list,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.UserServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}

	res, err := m.ChatService.BatchGetLatestMessage(ctx, &chat.DouyinMessageBatchGetLatestRequest{
		UserId:    viewerId,
		ToUserIds: list,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.ChatServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}

	fUser := make([]*base.FriendUser, len(resp.Users))
	for i, u := range resp.Users {
		fu := &base.FriendUser{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      u.IsFollow,
			Message:       res.Message[i],
			MsgType:       res.MsgType[i],
		}
		fUser = append(fUser, fu)
	}
	return fUser, nil
}
