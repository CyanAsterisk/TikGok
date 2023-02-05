package main

import (
	"context"
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
	GetFollowerCount(ctx context.Context, userId int64) (count int64, err error)
	GetFollowingCount(ctx context.Context, userId int64) (count int64, err error)
	CheckFollow(ctx context.Context, userId int64, toUserId int64) (bool, error)

	BatchGetFollowerCount(ctx context.Context, userIds []int64) (counts []int64, err error)
	BatchGetFollowingCount(ctx context.Context, userIds []int64) (counts []int64, err error)
	BatchCheckFollow(ctx context.Context, userId int64, toUserIds []int64) (checks []bool, err error)
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

	followerCount, err := s.GetFollowerCount(ctx, req.OwnerId)
	if err != nil {
		klog.Error("get followerList err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("get followerList err"))
		return resp, nil
	}
	followingCount, err := s.GetFollowingCount(ctx, req.OwnerId)
	if err != nil {
		klog.Error("get followingList err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("get followingList err"))
		return resp, nil
	}

	isFollow, err := s.CheckFollow(ctx, req.ViewerId, req.OwnerId)
	if err != nil {
		klog.Error("check follow err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("check follow err"))
		return resp, nil
	}

	resp.User = pkg.PackUser(usr, followerCount, followingCount, isFollow)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	users, err := dao.BatchGetUserById(req.OwnerIds)
	if err != nil {
		klog.Error("batch get user by id err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCSocialityErr.WithMessage("batch get user by id err"))
		return resp, nil
	}
	followerCnt, err := s.SocialManager.BatchGetFollowerCount(ctx, req.OwnerIds)
	if err != nil {
		klog.Error("batch get user follower count err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCSocialityErr.WithMessage("batch get user follower count err"))
		return resp, nil
	}
	followingCnt, err := s.SocialManager.BatchGetFollowingCount(ctx, req.OwnerIds)
	if err != nil {
		klog.Error("batch get user following count err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCSocialityErr.WithMessage("batch get user following count err"))
		return resp, nil
	}

	isFollow, err := s.SocialManager.BatchCheckFollow(ctx, req.ViewerId, req.OwnerIds)
	if err != nil {
		klog.Error("batch check follow err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.RPCSocialityErr.WithMessage("batch check follow err"))
		return resp, nil
	}

	resp.Users = pkg.PackUsers(users, followerCnt, followingCnt, isFollow)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}
