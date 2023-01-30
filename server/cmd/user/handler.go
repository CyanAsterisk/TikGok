package main

import (
	"context"
	"time"

	models "github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/middleware"
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	jwt *middleware.JWT
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(_ context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp = new(user.DouyinUserRegisterResponse)

	var usr model.User
	usr.Username = req.Username
	usr.Password = tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) // Encrypt password with md5.

	if err = dao.CreateUser(&usr); err != nil {
		if err == gorm.ErrRegistered {
			resp.BaseResp = sTools.BuildBaseResp(errno.UserAlreadyExistErr)
		} else {
			klog.Error("create user error", err)
			resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr.WithMessage("create user error"))
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
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr.WithMessage("create token error"))
		return resp, nil
	}

	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(_ context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	resp = new(user.DouyinUserLoginResponse)

	usr, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp.BaseResp = sTools.BuildBaseResp(errno.UserNotFoundErr)
		} else {
			klog.Errorf("get user by name err", err)
			resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr.WithMessage("get user by name err"))
		}
		return resp, nil
	}

	if usr.Password != tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) {
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr.WithMessage("wrong password"))
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
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr)
		return resp, nil
	}

	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	resp = new(user.DouyinUserResponse)

	usr, err := dao.GetUserById(req.UserId)
	if err != nil {
		klog.Error("get user by id failed", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.UserServerErr)
		return resp, nil
	}
	resp.User = tools.User(usr)

	res, err := global.SocialClient.FollowerList(ctx, &sociality.DouyinRelationFollowerListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		klog.Error("get followerList err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCSocialityErr)
		return resp, nil
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		errMsg := res.BaseResp.StatusMsg
		klog.Error("get followerList err", errMsg)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCSocialityErr.WithMessage(errMsg))
		return resp, nil
	}
	resp.User.FollowerCount = int64(len(res.UserList))

	for _, u := range res.UserList {
		if u.Id == req.UserId {
			resp.User.IsFollow = true
		}
	}

	response, err := global.SocialClient.FollowingList(ctx, (*sociality.DouyinRelationFollowListRequest)(&sociality.DouyinRelationFollowerListRequest{
		UserId: req.UserId,
	}))
	if err != nil {
		klog.Error("get followerList err", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCSocialityErr)
		return resp, nil
	}
	if response.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		errMsg := response.BaseResp.StatusMsg
		klog.Error("get followerList err", errMsg)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCInteractionErr.WithMessage(errMsg))
		return resp, nil
	}
	resp.User.FollowCount = int64(len(response.UserList))

	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}
