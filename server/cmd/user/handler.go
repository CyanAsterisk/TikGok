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
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
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
	ChatManager
}

// SocialManager defines the Anti Corruption Layer
// for get social logic.
type SocialManager interface {
	GetRelationList(ctx context.Context, viewerId int64, ownerId int64, option int8) ([]int64, error)
	GetSocialInfo(ctx context.Context, viewerId int64, ownerId int64) (*base.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error)
}

// ChatManager defines the Anti Corruption Layer
// for get chat logic.
type ChatManager interface {
	BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*base.LatestMsg, error)
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
	if err != nil {
		klog.Error("get user social info err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get user social info err"))
		return resp, nil
	}
	resp.User = pkg.PackUser(usr, info)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	resp = new(user.DouyinBatchGetUserResonse)
	userList, err := dao.BatchGetUserById(req.OwnerIdList)
	if err != nil {
		klog.Error("batch get user by id err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get user by id err"))
		return resp, nil
	}
	infoList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, req.OwnerIdList)

	resp.UserList = pkg.PackUsers(userList, infoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.DouyinGetRelationFollowListRequest) (resp *user.DouyinGetRelationFollowListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.ViewerId, consts.FollowList)
	if err != nil {
		klog.Error("get follow list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follow list err"))
		return resp, nil
	}
	userList, err := dao.BatchGetUserById(userIdList)
	if err != nil {
		klog.Error("batch get user by id err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get user by id err"))
		return resp, nil
	}

	infoList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
	if err != nil {
		klog.Error("batch get user info list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get user info list err"))
		return resp, nil
	}

	resp.UserList = pkg.PackUsers(userList, infoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.DouyinGetRelationFollowerListRequest) (resp *user.DouyinGetRelationFollowerListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowerListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.ViewerId, consts.FollowerList)
	if err != nil {
		klog.Error("get follower list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower list err"))
		return resp, nil
	}
	userList, err := dao.BatchGetUserById(userIdList)
	if err != nil {
		klog.Error("batch get user info list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get user info list err"))
		return resp, nil
	}

	infoList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
	if err != nil {
		klog.Error("batch get user info list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get user info list err"))
		return resp, nil
	}

	resp.UserList = pkg.PackUsers(userList, infoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.DouyinGetRelationFriendListRequest) (resp *user.DouyinGetRelationFriendListResponse, err error) {
	resp = new(user.DouyinGetRelationFriendListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.ViewerId, consts.FriendsList)
	if err != nil {
		klog.Error("get friend list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get friend list err"))
		return resp, nil
	}
	userList, err := dao.BatchGetUserById(userIdList)
	if err != nil {
		klog.Error("batch get user  list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("batch get user list err"))
		return resp, nil
	}

	infoList, err := s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
	if err != nil {
		klog.Error("batch get social info list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("batch get social info list err"))
		return resp, nil
	}
	msgList, err := s.ChatManager.BatchGetLatestMessage(ctx, req.ViewerId, userIdList)
	if err != nil {
		klog.Error("batch get user latest message list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("batch get user latest message list err"))
		return resp, nil
	}

	resp.UserList = pkg.PackFriendUsers(userList, infoList, msgList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
