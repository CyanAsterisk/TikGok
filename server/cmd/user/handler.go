package main

import (
	"context"
	"sync"
	"time"

	models "github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/config"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/middleware"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	Jwt *middleware.JWT
	SocialManager
	ChatManager
	RedisManager
	InteractionManager
	Dao *dao.User
}

// SocialManager defines the Anti Corruption Layer
// for get social logic.
type SocialManager interface {
	GetRelationList(ctx context.Context, viewerId, ownerId int64, option int8) ([]int64, error)
	GetSocialInfo(ctx context.Context, viewerId, ownerId int64) (*base.SocialInfo, error)
	BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error)
}

type InteractionManager interface {
	GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error)
	BatchGetInteractInfo(ctx context.Context, userIdList []int64) ([]*base.UserInteractInfo, error)
}

// ChatManager defines the Anti Corruption Layer
// for get chat logic.
type ChatManager interface {
	BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*base.LatestMsg, error)
}

// RedisManager defines the redis interface.
type RedisManager interface {
	GetUserById(ctx context.Context, uid int64) (*model.User, error)
	BatchGetUserById(ctx context.Context, uidList []int64) ([]*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp = new(user.DouyinUserRegisterResponse)

	sf, err := snowflake.NewNode(consts.UserSnowflakeNode)
	if err != nil {
		klog.Errorf("generate user id failed: %s", err.Error())
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("generate user id failed"))
		return resp, nil
	}
	usr := &model.User{
		ID:       sf.Generate().Int64(),
		Username: req.Username,
		Password: pkg.Md5Crypt(req.Password, config.GlobalServerConfig.MysqlInfo.Salt), // Encrypt password with md5.
		// TODO: Add logic to set avatar backgroundImage and signature
		Avatar:          "https://w.wallhaven.cc/full/y8/wallhaven-y8lqo7.jpg",
		BackGroundImage: "https://w.wallhaven.cc/full/zy/wallhaven-zyxvqy.jpg",
		Signature:       "default signature",
	}
	if err = s.Dao.CreateUser(usr); err != nil {
		if err == dao.ErrUserExist {
			resp.BaseResp = tools.BuildBaseResp(errno.UserAlreadyExistErr)
		} else {
			klog.Error("create user error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("create user error"))
		}
		return resp, nil
	}
	if err = s.RedisManager.CreateUser(ctx, usr); err != nil {
		klog.Errorf("create user error by redis error")
	}

	resp.UserId = usr.ID
	resp.Token, err = s.Jwt.CreateToken(models.CustomClaims{
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

	usr, err := s.Dao.GetUserByUsername(req.Username)
	if err != nil {
		if err == dao.ErrNoSuchUser {
			resp.BaseResp = tools.BuildBaseResp(errno.UserNotFoundErr)
		} else {
			klog.Errorf("get user by name err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("get user by name err"))
		}
		return resp, nil
	}

	if usr.Password != pkg.Md5Crypt(req.Password, config.GlobalServerConfig.MysqlInfo.Salt) {
		resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("wrong password"))
		return resp, nil
	}

	resp.UserId = usr.ID
	resp.Token, err = s.Jwt.CreateToken(models.CustomClaims{
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

	usr, err := s.RedisManager.GetUserById(ctx, req.OwnerId)
	if err != nil {
		klog.Error("get user by redis err", err)
		if usr, err = s.Dao.GetUserById(req.OwnerId); err != nil {
			if err == dao.ErrNoSuchUser {
				resp.BaseResp = tools.BuildBaseResp(errno.UserNotFoundErr)
			} else {
				klog.Error("get user by id failed", err)
				resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr)
			}
			return resp, nil
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var socialInfo *base.SocialInfo
	go func() {
		defer wg.Done()
		socialInfo, err = s.SocialManager.GetSocialInfo(ctx, req.ViewerId, req.OwnerId)
		if err != nil {
			klog.Error("get user social info err", err)
		}
	}()
	var interactInfo *base.UserInteractInfo
	go func() {
		defer wg.Done()
		interactInfo, err = s.InteractionManager.GetInteractInfo(ctx, req.OwnerId)
		if err != nil {
			klog.Error("get user interact info err", err)
		}
	}()
	wg.Wait()
	if err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.User = pkg.PackUser(usr, socialInfo, interactInfo)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) BatchGetUserInfo(ctx context.Context, req *user.DouyinBatchGetUserRequest) (resp *user.DouyinBatchGetUserResonse, err error) {
	resp = new(user.DouyinBatchGetUserResonse)
	userList, err := s.RedisManager.BatchGetUserById(ctx, req.OwnerIdList)
	if err != nil {
		klog.Error("batch get user by redis err", err)
		if userList, err = s.Dao.BatchGetUserById(req.OwnerIdList); err != nil {
			klog.Error("batch get user by id err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("batch get user by id err"))
			return resp, nil
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)

	var socialInfoList []*base.SocialInfo
	go func() {
		defer wg.Done()
		socialInfoList, err = s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, req.OwnerIdList)
		if err != nil {
			klog.Error("batch get social info error", err)
		}
	}()

	var interactInfoList []*base.UserInteractInfo
	go func() {
		defer wg.Done()
		interactInfoList, err = s.InteractionManager.BatchGetInteractInfo(ctx, req.OwnerIdList)
		if err != nil {
			klog.Error("batch get interact info error", err)
		}
	}()
	wg.Wait()
	if err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.UserList = pkg.PackUsers(userList, socialInfoList, interactInfoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetFollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowList(ctx context.Context, req *user.DouyinGetRelationFollowListRequest) (resp *user.DouyinGetRelationFollowListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FollowList)
	if err != nil {
		klog.Error("get follow list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follow list err"))
		return resp, nil
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdList)
	if err != nil {
		klog.Error("batch get user by redis err", err)
		if userList, err = s.Dao.BatchGetUserById(userIdList); err != nil {
			klog.Error("batch get user by id err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("batch get user by id err"))
			return resp, nil
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)

	var socialInfoList []*base.SocialInfo
	go func() {
		defer wg.Done()
		socialInfoList, err = s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
		if err != nil {
			klog.Error("batch get social info error", err)
		}
	}()

	var interactInfoList []*base.UserInteractInfo
	go func() {
		defer wg.Done()
		interactInfoList, err = s.InteractionManager.BatchGetInteractInfo(ctx, userIdList)
		if err != nil {
			klog.Error("batch get interact info error", err)
		}
	}()
	wg.Wait()
	if err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}

	resp.UserList = pkg.PackUsers(userList, socialInfoList, interactInfoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetFollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFollowerList(ctx context.Context, req *user.DouyinGetRelationFollowerListRequest) (resp *user.DouyinGetRelationFollowerListResponse, err error) {
	resp = new(user.DouyinGetRelationFollowerListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FollowerList)
	if err != nil {
		klog.Error("get follower list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower list err"))
		return resp, nil
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdList)
	if err != nil {
		klog.Error("batch get user by redis err", err)
		if userList, err = s.Dao.BatchGetUserById(userIdList); err != nil {
			klog.Error("batch get user by id err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("batch get user by id err"))
			return resp, nil
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)

	var socialInfoList []*base.SocialInfo
	go func() {
		defer wg.Done()
		socialInfoList, err = s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
		if err != nil {
			klog.Error("batch get social info error", err)
		}
	}()

	var interactInfoList []*base.UserInteractInfo
	go func() {
		defer wg.Done()
		interactInfoList, err = s.InteractionManager.BatchGetInteractInfo(ctx, userIdList)
		if err != nil {
			klog.Error("batch get interact info error", err)
		}
	}()
	wg.Wait()
	if err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.UserList = pkg.PackUsers(userList, socialInfoList, interactInfoList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, req *user.DouyinGetRelationFriendListRequest) (resp *user.DouyinGetRelationFriendListResponse, err error) {
	resp = new(user.DouyinGetRelationFriendListResponse)
	userIdList, err := s.SocialManager.GetRelationList(ctx, req.ViewerId, req.OwnerId, consts.FriendsList)
	if err != nil {
		klog.Error("get friend list err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get friend list err"))
		return resp, nil
	}
	userList, err := s.RedisManager.BatchGetUserById(ctx, userIdList)
	if err != nil {
		klog.Error("batch get user by redis err", err)
		if userList, err = s.Dao.BatchGetUserById(userIdList); err != nil {
			klog.Error("batch get user by id err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.UserServerErr.WithMessage("batch get user by id err"))
			return resp, nil
		}
	}
	var wg sync.WaitGroup
	wg.Add(3)

	var socialInfoList []*base.SocialInfo
	go func() {
		defer wg.Done()
		socialInfoList, err = s.SocialManager.BatchGetSocialInfo(ctx, req.ViewerId, userIdList)
		if err != nil {
			klog.Error("batch get social info error", err)
		}
	}()

	var interactInfoList []*base.UserInteractInfo
	go func() {
		defer wg.Done()
		interactInfoList, err = s.InteractionManager.BatchGetInteractInfo(ctx, userIdList)
		if err != nil {
			klog.Error("batch get interact info error", err)
		}
	}()
	var msgList []*base.LatestMsg
	go func() {
		defer wg.Done()
		msgList, err = s.ChatManager.BatchGetLatestMessage(ctx, req.ViewerId, userIdList)
		if err != nil {
			klog.Error("batch get user latest message list err", err)
		}
	}()
	wg.Wait()
	if err != nil {
		resp.BaseResp = tools.BuildBaseResp(errno.ServiceErr)
		return resp, nil
	}
	resp.UserList = pkg.PackFriendUsers(userList, socialInfoList, interactInfoList, msgList)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
