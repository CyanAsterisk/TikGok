package dao

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"gorm.io/gorm"
)

var (
	ErrNoSuchUser = errors.New("no such user")
	ErrUserExist  = errors.New("user already exist")
)

// GetUserByUsername get user by username
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).
		Where(&model.User{Username: username}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoSuchUser
	}
	return &user, err
}

// GetUserById get user by username
func GetUserById(uid int64) (*model.User, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).
		Where(&model.User{ID: uid}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoSuchUser
	}
	return &user, err
}

// CreateUser creates a user.
func CreateUser(user *model.User) error {
	err := global.DB.Model(&model.User{}).
		Where(&model.User{Username: user.Username}).First(&user).Error
	if err == nil {
		return ErrUserExist
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return global.DB.Model(&model.User{}).Create(user).Error
}
