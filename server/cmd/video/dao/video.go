package dao

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

type Video struct {
	db *gorm.DB
}

var (
	ErrNoSuchRecord       = errors.New("no such video record")
	ErrRecordAlreadyExist = errors.New("video record already exist")
)

// NewVideo create a video dao.
func NewVideo(db *gorm.DB) *Video {
	m := db.Migrator()
	if !m.HasTable(&model.Video{}) {
		err := m.CreateTable(&model.Video{})
		if err != nil {
			panic(err)
		}
	}
	return &Video{
		db: db,
	}
}

// CreateVideo creates a new video record.
func (v *Video) CreateVideo(video *model.Video) error {
	err := v.db.Model(&model.Video{}).
		Where(&model.Video{ID: video.ID}).First(&model.Video{}).Error
	if err == nil {
		return ErrRecordAlreadyExist
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return v.db.Create(&video).Error
}

// GetVideoListByLatestTime gets videos for feed.
func (v *Video) GetVideoListByLatestTime(latestTime int64) ([]*model.Video, error) {
	videos := make([]*model.Video, 0)
	if err := v.db.Where("create_time <= ?", latestTime).
		Order("create_time desc").
		Limit(consts.VideosLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideoListByAuthorId gets videos by userId of author.
func (v *Video) GetVideoListByAuthorId(AuthorId int64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	if err := v.db.Where(&model.Video{AuthorId: AuthorId}).Order("create_time desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// GetVideoIdListByAuthorId gets videos by userId of author.
func (v *Video) GetVideoIdListByAuthorId(AuthorId int64) ([]int64, error) {
	var list []int64
	if err := v.db.Model(model.Video{}).Where(&model.Video{AuthorId: AuthorId}).
		Pluck("id", &list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetVideoByVideoId gets video by videoId
func (v *Video) GetVideoByVideoId(vid int64) (*model.Video, error) {
	var video model.Video
	err := v.db.Model(&model.Video{}).
		Where(&model.Video{ID: vid}).First(&video).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrNoSuchRecord
		}
		return nil, err
	}
	return &video, nil
}

// BatchGetVideoByVideoId gets video list by videoId list.
func (v *Video) BatchGetVideoByVideoId(vidList []int64) ([]*model.Video, error) {
	if vidList == nil {
		return nil, nil
	}
	var vl []*model.Video
	for _, vid := range vidList {
		v, err := v.GetVideoByVideoId(vid)
		if err != nil {
			return nil, err
		}
		vl = append(vl, v)
	}
	return vl, nil
}

// DeleteVideoById delete video by id.
func (v *Video) DeleteVideoById(videoId int64) error {
	err := v.db.Model(&model.Video{}).
		Where(&model.Video{ID: videoId}).First(&model.Video{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNoSuchRecord
		} else {
			return err
		}
	}
	return v.db.Model(&model.Video{}).Delete(&model.Video{ID: videoId}).Error
}
