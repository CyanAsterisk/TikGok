package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"time"

	models "github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/middleware"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	jwt *middleware.JWT
	SocialManager
}

// SocialManager defines the Anti Corruption Layer
// for get social logic.
type SocialManager interface {
	GetSocialInfo(ctx context.Context, viewerId int64, ownerId int64) (*base.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error)
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(_ context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp = new(user.DouyinUserRegisterResponse)

	var usr model.User
	usr.Username = req.Username
	usr.Password = pkg.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) // Encrypt password with md5.

	if err = dao.CreateUser(&usr); err != nil {
		if err == dao.ErrUserExist {
			resp.BaseResp = tools.BuildBaseResp(errno.UserAlreadyExistErr)
		} else {
			klog.Error("create user error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("create user error"))
		}
		return resp, nil
	}

	resp.UserId = usr.ID
	resp.Token, err = s.jwt.CreateToken(models.CustomClaims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.ThirtyDays,
			Issuer:    consts.JWTIssuer,
		},
	})
	if err != nil {
		klog.Error("create token err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("create token error"))
		return resp, nil
	}

	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(_ context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	resp = new(user.DouyinUserLoginResponse)

	usr, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		if err == dao.ErrNoSuchUser {
			resp.BaseResp = tools.BuildBaseResp(errno.UserNotFoundErr)
		} else {
			klog.Errorf("get user by name err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("get user by name err"))
		}
		return resp, nil
	}

	if usr.Password != pkg.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) {
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("wrong password"))
		return resp, nil
	}

	resp.UserId = usr.ID
	resp.Token, err = s.jwt.CreateToken(models.CustomClaims{
		ID: usr.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.ThirtyDays,
			Issuer:    consts.JWTIssuer,
		},
	})
	if err != nil {
		klog.Error("create token err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr)
		return resp, nil
	}

	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.DouyinGetUserRequest) (resp *user.DouyinGetUserResponse, err error) {
	resp = new(user.DouyinGetUserResponse)

	usr, err := dao.GetUserById(req.OwnerId)
	if err != nil {
		if err == dao.ErrNoSuchUser {
			resp.BaseResp = tools.BuildBaseResp(errno.UserNotFoundErr)
			return resp, nil
		}
		klog.Error("get user by id failed", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr)
		return resp, nil
	}
	info, err := s.SocialManager.GetSocialInfo(ctx, req.ViewerId, req.OwnerId)

	resp.User = pkg.PackUser(usr, info)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	users, err := dao.BatchGetUserById(req.OwnerIdList)
	if err != nil {
		klog.Error("batch get user by id err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCSocialityErr.WithMessage("batch get user by id err"))
		return resp, nil
	}
	infoList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, req.OwnerIdList)

	resp.UserList = pkg.PackUsers(users, infoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.DouyinGetRelationFollowListRequest) (resp *user.DouyinGetRelationFollowListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.DouyinGetRelationFollowerListRequest) (resp *user.DouyinGetRelationFollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.DouyinGetRelationFriendListRequest) (resp *user.DouyinGetRelationFriendListResponse, err error) {
	// TODO: Your code here...
	return
}
