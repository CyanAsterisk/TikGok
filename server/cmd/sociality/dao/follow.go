package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

type Follow struct {
	db *gorm.DB
}

// NewFollow create a social follow dao.
func NewFollow(db *gorm.DB) *Follow {
	m := db.Migrator()
	if !m.HasTable(&model.Follow{}) {
		err := m.CreateTable(&model.Follow{})
		if err != nil {
			panic(err)
		}
	}
	return &Follow{
		db: db,
	}
}

// GetFollowerNumsByUserId gets follower nums by userId.
func (f *Follow) GetFollowerNumsByUserId(userId int64) (int64, error) {
	var num int64
	err := f.db.Model(&model.Follow{}).
		Where(&model.Follow{UserId: userId, ActionType: consts.IsFollow}).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

// GetFollowNumsByUserId gets following nums by userId.
func (f *Follow) GetFollowNumsByUserId(userId int64) (int64, error) {
	var num int64
	err := f.db.Model(&model.Follow{}).
		Where(&model.Follow{FollowerId: userId, ActionType: consts.IsFollow}).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

// CreateFollow creates a follow record.
func (f *Follow) CreateFollow(follow *model.Follow) error {
	err := f.db.Model(model.Follow{}).
		Create(&follow).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFollow to update follow status.
func (f *Follow) UpdateFollow(userId, followId int64, actionType int8) error {
	err := f.db.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, FollowerId: followId}).Update("action_type", actionType).Error
	if err != nil {
		return err
	}
	return nil
}

// FindRecord to find if there's a record between user and another user.
func (f *Follow) FindRecord(userId, followerId int64) (*model.Follow, error) {
	var follow *model.Follow
	err := f.db.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, FollowerId: followerId}).First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return follow, nil
}

// GetFollowerIdList gets followerId list.
func (f *Follow) GetFollowerIdList(userId int64) ([]int64, error) {
	var list []int64
	err := f.db.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, ActionType: consts.IsFollow}).Pluck("follower_id", &list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetFollowIdList gets followingId list.
func (f *Follow) GetFollowIdList(userId int64) ([]int64, error) {
	var list []int64
	err := f.db.Model(model.Follow{}).
		Where(&model.Follow{FollowerId: userId, ActionType: consts.IsFollow}).Pluck("user_id", &list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetFriendsList gets friends list.
func (f *Follow) GetFriendsList(userId int64) ([]int64, error) {
	var followingList []int64
	var friendsList []int64
	err := f.db.Model(model.Follow{}).
		Where(&model.Follow{FollowerId: userId, ActionType: consts.IsFollow}).Pluck("user_id", &followingList).
		FindInBatches(&followingList, 100, func(tx *gorm.DB, batch int) error {
			for _, c := range followingList {
				_, err := f.FindRecord(userId, c)
				if err != nil {
					return err
				}
				friendsList = append(friendsList, c)
			}
			tx.Save(&friendsList)
			return nil
		}).Error
	if err != nil {
		return nil, err
	}
	return friendsList, nil
}
