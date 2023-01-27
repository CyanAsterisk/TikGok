package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

// GetFavoriteUserList gets favorite user list by videoId.
func GetFavoriteUserList(videoId int64) ([]int64, error) {
	var userList []int64
	err := global.DB.Model(model.Favorite{}).
		Where(&model.Favorite{VideoId: videoId, ActionType: consts.IsLike}).Pluck("user_id", &userList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// CreateFavorite creates a favorite record.
func CreateFavorite(fav *model.Favorite) error {
	err := global.DB.Model(model.Favorite{}).
		Create(&fav).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFavorite updates favorite status.
func UpdateFavorite(userId, videoId int64, actionType int8) error {
	err := global.DB.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, VideoId: videoId}).Update("action_type", actionType).Error
	if err != nil {
		return err
	}
	return nil
}

// GetFavoriteInfo get favorite info.
func GetFavoriteInfo(userId, videoId int64) (*model.Favorite, error) {
	var fvInfo *model.Favorite
	err := global.DB.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, VideoId: videoId}).First(&fvInfo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return fvInfo, nil
}

// GetFavoriteVideoIdListByUserId gets a user's favorite video list
func GetFavoriteVideoIdListByUserId(userId int64) ([]int64, error) {
	var videoList []int64
	err := global.DB.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, ActionType: consts.IsLike}).Pluck("video_id", &videoList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return videoList, nil
}
