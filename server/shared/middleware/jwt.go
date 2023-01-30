package middleware

import (
	"context"
	"errors"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"strings"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	TokenNotFound    = errors.New("no token")
)

func JWTAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Request.Header.Get(consts.AuthorizationKey)
		if token == "" {
			c.JSON(http.StatusOK, utils.H{
				"status_code": int32(errno.ParamsEr.ErrCode),
				"status_msg":  TokenNotFound.Error(),
			})
			c.Abort()
			return
		}
		token = strings.Split(token, " ")[1]
		j := NewJWT()
		// Parse the information contained in the token
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, utils.H{
				"status_code": int32(errno.ParamsEr.ErrCode),
				"status_msg":  TokenInvalid.Error(),
			})
			c.Abort()
			return
		}
		c.Set(consts.Claims, claims)
		c.Set(consts.AccountID, claims.ID)
		c.Next(ctx)
	}
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

// CreateToken to create a token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken to parse a token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}

// RefreshToken to refresh a token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(consts.TokenRefreshTime).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
