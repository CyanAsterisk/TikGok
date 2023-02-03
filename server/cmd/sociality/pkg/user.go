package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
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
	var users []*base.User
	for _, tid := range list {
		u, err := m.UserService.GetUserInfo(ctx, &user.DouyinUserRequest{ViewerId: viewerId, OwnerId: tid})
		if err != nil {
			return nil, err
		}
		users = append(users, u.User)
	}
	return users, nil
}

func (m *UserManager) GetFriendUsers(ctx context.Context, list []int64, viewerId int64) ([]*base.FriendUser, error) {
	var fUsers []*base.FriendUser
	for _, oid := range list {
		u, err := m.UserService.GetUserInfo(ctx, &user.DouyinUserRequest{ViewerId: viewerId, OwnerId: oid})
		if err != nil {
			return nil, err
		}
		c, err := m.ChatService.LatestMessage(ctx, &chat.DouyinMessageLatestRequest{UserId: viewerId, ToUserId: oid})
		if err != nil {
			return nil, err
		}
		fu := &base.FriendUser{
			Id:            u.User.Id,
			Name:          u.User.Name,
			FollowCount:   u.User.FollowCount,
			FollowerCount: u.User.FollowerCount,
			IsFollow:      u.User.IsFollow,
			Message:       c.Message,
			MsgType:       c.MsgType,
		}
		fUsers = append(fUsers, fu)
	}
	return fUsers, nil
}
