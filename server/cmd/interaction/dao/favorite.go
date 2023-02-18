package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

type Favorite struct {
	db *gorm.DB
}

// NewFavorite create a interaction favorite dao.
func NewFavorite(db *gorm.DB) *Favorite {
	m := db.Migrator()
	if !m.HasTable(&model.Favorite{}) {
		err := m.CreateTable(&model.Favorite{})
		if err != nil {
			panic(err)
		}
	}
	return &Favorite{
		db: db,
	}
}

// GetFavoriteCountByVideoId gets the number of favorite by videoId.
func (f *Favorite) GetFavoriteCountByVideoId(videoId int64) (int64, error) {
	var count int64
	err := f.db.Model(&model.Favorite{}).
		Where(&model.Favorite{VideoId: videoId, ActionType: consts.IsLike}).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}

// GetFavoriteUserList gets favorite user list by videoId.
func (f *Favorite) GetFavoriteUserList(videoId int64) ([]int64, error) {
	var userList []int64
	err := f.db.Model(model.Favorite{}).
		Where(&model.Favorite{VideoId: videoId, ActionType: consts.IsLike}).
		Order("create_date desc").Pluck("user_id", &userList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return userList, nil
}

// CreateFavorite creates a favorite record.
func (f *Favorite) CreateFavorite(fav *model.Favorite) error {
	err := f.db.Model(&model.Favorite{}).
		Where(&model.Favorite{VideoId: fav.VideoId, UserId: fav.UserId}).First(&model.Favorite{}).Error
	if err == nil {
		return ErrRecordAlreadyExist
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return f.db.Model(model.Favorite{}).
		Create(&fav).Error
}

// UpdateFavorite updates favorite status.
func (f *Favorite) UpdateFavorite(userId, videoId int64, actionType int8) error {
	return f.db.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, VideoId: videoId}).
		Update("action_type", actionType).Error
}

// GetFavoriteInfo get favorite info.
func (f *Favorite) GetFavoriteInfo(userId, videoId int64) (*model.Favorite, error) {
	var fvInfo *model.Favorite
	err := f.db.Model(model.Favorite{}).
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
func (f *Favorite) GetFavoriteVideoIdListByUserId(userId int64) ([]int64, error) {
	var videoList []int64
	err := f.db.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, ActionType: consts.IsLike}).
		Order("create_date desc").Pluck("video_id", &videoList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

// GetFavoriteVideoCountByUserId gets a user's favorite video count.
func (f *Favorite) GetFavoriteVideoCountByUserId(userId int64) (int64, error) {
	var count int64
	err := f.db.Model(model.Favorite{}).
		Where(&model.Favorite{UserId: userId, ActionType: consts.IsLike}).
		Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}
