package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

// GetFollowerNumsByUserId gets follower nums by userId.
func GetFollowerNumsByUserId(userId int64) (int64, error) {
	var num int64
	err := global.DB.Model(&model.Follow{}).
		Where(&model.Follow{UserId: userId, ActionType: consts.IsFollow}).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

// GetFollowingNumsByUserId gets following nums by userId.
func GetFollowingNumsByUserId(userId int64) (int64, error) {
	var num int64
	err := global.DB.Model(&model.Follow{}).
		Where(&model.Follow{FollowerId: userId, ActionType: consts.IsFollow}).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}

// CreateFollow creates a follow record.
func CreateFollow(follow *model.Follow) error {
	err := global.DB.Model(model.Follow{}).
		Create(&follow).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFollow to update follow status.
func UpdateFollow(userId, followId int64, actionType int8) error {
	err := global.DB.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, FollowerId: followId}).Update("action_type", actionType).Error
	if err != nil {
		return err
	}
	return nil
}

// FindRecord to find if there's a record between user and another user.
func FindRecord(userId, followId int64) (*model.Follow, error) {
	var follow *model.Follow
	err := global.DB.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, FollowerId: followId}).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return follow, nil
}

// GetFollowerIdList gets followerId list.
func GetFollowerIdList(userId int64) ([]int64, error) {
	var list []int64
	err := global.DB.Model(model.Follow{}).
		Where(&model.Follow{UserId: userId, ActionType: consts.IsFollow}).Pluck("follower_id", &list).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetFollowingIdList gets followingId list.
func GetFollowingIdList(userId int64) ([]int64, error) {
	var list []int64
	err := global.DB.Model(model.Follow{}).
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
func GetFriendsList(userId int64) ([]int64, error) {
	var followingList []int64
	var list []int64
	err := global.DB.Model(model.Follow{}).
		Where(&model.Follow{FollowerId: userId, ActionType: consts.IsFollow}).Pluck("user_id", &list).
		FindInBatches(&followingList, 100, func(tx *gorm.DB, batch int) error {
			for _, c := range followingList {
				_, err := FindRecord(userId, c)
				if err != nil {
					return err
				}
				list = append(list, c)
			}
			tx.Save(&list)
			return nil
		}).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
