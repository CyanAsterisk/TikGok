package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
)

// GetUserByUsername get user by username
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).
		Where(&model.User{Username: username}).First(&user).Error
	return &user, err
}

// GetUserById get user by username
func GetUserById(uid int64) (*model.User, error) {
	var user model.User
	err := global.DB.Model(&model.User{}).
		Where(&model.User{ID: uid}).First(&user).Error
	return &user, err
}

// CreateUser creates a user.
func CreateUser(user *model.User) error {
	return global.DB.Model(&model.User{}).Create(user).Error
}
