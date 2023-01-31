package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/CyanAsterisk/TikGok/server/cmd/api/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	TokenNotFound    = errors.New("no token")
)

func JWTAuth(secretKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.Query(consts.AuthorizationKey)
		if token == "" {
			token = string(c.FormValue(consts.AuthorizationKey))
			if token == "" {
				c.JSON(http.StatusOK, utils.H{
					"status_code": int32(errno.ParamsEr.ErrCode),
					"status_msg":  TokenNotFound.Error(),
				})
				c.Abort()
				return
			}
		}
		j := NewJWT(secretKey)
		// Parse the information contained in the token
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, utils.H{
				"status_code": int32(errno.ParamsEr.ErrCode),
				"status_msg":  err.Error(),
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

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SigningKey: []byte(secretKey),
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
