package pkg

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
	res, err := m.UserService.BatchGetUserInfo(ctx, &user.DouyinBatchGetUserRequest{
		ViewerId: viewerId,
		OwnerIds: list,
	})
	if err != nil {
		return nil, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.UserServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Users, nil
}

// GetUser gets user info.
func (m *UserManager) GetUser(ctx context.Context, viewerId, ownerId int64) (*base.User, error) {
	resp, err := m.UserService.GetUserInfo(ctx, &user.DouyinGetUserRequest{ViewerId: viewerId, OwnerId: ownerId})
	if err != nil {
		return nil, err
	}
	if int64(resp.BaseResp.StatusCode) != errno.Success.ErrCode {
		return nil, errno.UserServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.User, nil
}
