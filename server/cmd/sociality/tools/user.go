package tools

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user/userservice"
)

type UserManager struct {
	UserService userservice.Client
}

// GetUsers gets users info by list.
func (m *UserManager) GetUsers(ctx context.Context, list []int64) ([]*base.User, error) {
	var users []*base.User
	for _, uid := range list {
		u, err := m.UserService.GetUserInfo(ctx, &user.DouyinUserRequest{UserId: uid})
		if err != nil {
			return nil, err
		}
		users = append(users, u.User)
	}
	return users, nil
}
