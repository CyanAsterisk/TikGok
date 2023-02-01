package tools

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
)

type UserManager struct {
	UserService userservice.Client
}

// GetUsers gets users info by list.
func (m *UserManager) GetUsers(ctx context.Context, list []int64, viewerId int64) ([]*base.User, error) {
	var users []*base.User
	for _, oid := range list {
		u, err := m.UserService.GetUserInfo(ctx, &user.DouyinUserRequest{ViewerId: viewerId, OwnerId: oid})
		if err != nil {
			return nil, err
		}
		users = append(users, u.User)
	}
	return users, nil
}

// GetUser gets user info.
func (m *UserManager) GetUser(ctx context.Context, viewerId, ownerId int64) (*base.User, error) {
	resp, err := m.UserService.GetUserInfo(ctx, &user.DouyinUserRequest{ViewerId: viewerId, OwnerId: ownerId})
	if err != nil {
		return nil, err
	}
	if int64(resp.BaseResp.StatusCode) != errno.Success.ErrCode {
		return nil, errno.UserServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.User, nil
}
