package dao

import (
	"errors"
	"sync"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"gorm.io/gorm"
)

var (
	ErrNoSuchUser = errors.New("no such user")
	ErrUserExist  = errors.New("user already exist")
)

type User struct {
	db *gorm.DB
}

// NewUser create a user dao.
func NewUser(db *gorm.DB) *User {
	m := db.Migrator()
	if !m.HasTable(&model.User{}) {
		err := m.CreateTable(&model.User{})
		if err != nil {
			panic(err)
		}
	}
	return &User{
		db: db,
	}
}

// GetUserByUsername get user by username
func (u *User) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := u.db.Model(&model.User{}).
		Where(&model.User{Username: username}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoSuchUser
	}
	return &user, err
}

// GetUserById get user by userid.
func (u *User) GetUserById(uid int64) (*model.User, error) {
	var user model.User
	err := u.db.Model(&model.User{}).
		Where(&model.User{ID: uid}).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, ErrNoSuchUser
	}
	return &user, err
}

// BatchGetUserById get users by userid.
func (u *User) BatchGetUserById(uids []int64) ([]*model.User, error) {
	if uids == nil {
		return nil, nil
	}

	length := len(uids)
	users := make([]*model.User, length)
	var wg sync.WaitGroup
	wg.Add(length)
	for i := range uids {
		go func(idx int) {
			defer wg.Done()
			var user model.User
			u.db.Model(&model.User{}).
				Where(&model.User{ID: uids[idx]}).First(&user)
			users[idx] = &user
		}(i)
	}
	wg.Wait()
	return users, nil
}

// CreateUser creates a user.
func (u *User) CreateUser(user *model.User) error {
	err := u.db.Model(&model.User{}).
		Where(&model.User{Username: user.Username}).First(&model.User{}).Error
	if err == nil {
		return ErrUserExist
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return u.db.Model(&model.User{}).Create(user).Error
}

// DeleteUserById delete a user by id.
func (u *User) DeleteUserById(userId int64) error {
	return u.db.Model(&model.User{}).Delete(&model.User{ID: userId}).Error
}
