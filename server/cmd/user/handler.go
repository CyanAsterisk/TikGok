package main

import (
	"context"
	"time"

	models "github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
	"github.com/CyanAsterisk/TikGok/server/shared/middleware"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/codes"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/status"
	"github.com/golang-jwt/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	jwt *middleware.JWT
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(_ context.Context, req *user.DouyinUserRegisterRequest) (*user.DouyinUserRegisterResponse, error) {
	var _user model.User
	result := global.DB.Where(&model.User{Username: req.Username}).First(&_user)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "Account already exists")
	}
	_user.Username = req.Username
	_user.Password = tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) // Encrypt password with md5.
	result = global.DB.Create(&_user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	token, err := s.jwt.CreateToken(models.CustomClaims{
		ID: _user.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.ThirtyDays,
			Issuer:    consts.JWTIssuer,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create token error: %s", err.Error())
	}
	return &user.DouyinUserRegisterResponse{
		UserId: _user.ID,
		Token:  token,
	}, nil

}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(_ context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	var _user model.User
	result := global.DB.Where(&model.User{Username: req.Username}).First(&_user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no such user")
	}

	if _user.Password != tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) {
		return nil, status.Errorf(codes.PermissionDenied, "wrong password")
	}
	token, err := s.jwt.CreateToken(models.CustomClaims{
		ID: _user.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.ThirtyDays,
			Issuer:    consts.JWTIssuer,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create token error: %s", err.Error())
	}
	return &user.DouyinUserLoginResponse{
		UserId: _user.ID,
		Token:  token,
	}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.DouyinUserRequest) (resp *user.DouyinUserRequest, err error) {
	// TODO: Your code here...
	return
}
