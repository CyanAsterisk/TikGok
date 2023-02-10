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

// GetUserById get user by userid.
func GetUserById(uid int64) (*model.User, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).
		Where(&model.User{ID: uid}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoSuchUser
	}
	return &user, err
}

// BatchGetUserById get users by userid.
func BatchGetUserById(uids []int64) ([]*model.User, error) {
	if uids == nil {
		return nil, nil
	}
	users := make([]*model.User, 0)
	for _, id := range uids {
		var user model.User
		err := global.DB.Model(&model.User{}).
			Where(&model.User{ID: id}).First(&user).Error
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser creates a user.
func CreateUser(user *model.User) error {
	err := global.DB.Model(&model.User{}).
		Where(&model.User{Username: user.Username}).First(&model.User{}).Error
	if err == nil {
		return ErrUserExist
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return global.DB.Model(&model.User{}).Create(user).Error
}
