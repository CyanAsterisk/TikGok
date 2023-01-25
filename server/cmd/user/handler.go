package main

import (
	"context"
	"time"

	models "github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/tools"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/errno"
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
	var usr model.User
	result := global.DB.Where(&model.User{Username: req.Username}).First(&usr)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "Account already exists")
	}
	usr.Username = req.Username
	usr.Password = tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) // Encrypt password with md5.
	result = global.DB.Create(&usr)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	token, err := s.jwt.CreateToken(models.CustomClaims{
		ID: usr.ID,
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
		StatusCode: int32(errno.Err_Success),
		UserId:     usr.ID,
		Token:      token,
	}, nil
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(_ context.Context, req *user.DouyinUserLoginRequest) (*user.DouyinUserLoginResponse, error) {
	var usr model.User
	result := global.DB.Where(&model.User{Username: req.Username}).First(&usr)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "no such user")
	}

	if usr.Password != tools.Md5Crypt(req.Password, global.ServerConfig.MysqlInfo.Salt) {
		return nil, status.Errorf(codes.PermissionDenied, "wrong password")
	}
	token, err := s.jwt.CreateToken(models.CustomClaims{
		ID: usr.ID,
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
		StatusCode: int32(errno.Err_Success),
		UserId:     usr.ID,
		Token:      token,
	}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.DouyinUserRequest) (*user.DouyinUserResponse, error) {
	_, err := s.jwt.ParseToken(req.Token)
	if err != nil {
		if err == middleware.TokenExpired {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
	}
	var usr model.User
	result := global.DB.Where(&model.User{ID: req.UserId}).First(&usr)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}
	return &user.DouyinUserResponse{
		StatusCode: int32(errno.Err_Success),
		StatusMsg:  "",
		User: &user.User{
			Id:   usr.ID,
			Name: usr.Username,
		},
	}, nil
}
