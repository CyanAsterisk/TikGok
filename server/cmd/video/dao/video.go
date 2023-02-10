package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"gorm.io/gorm"
)

type Video struct {
	db *gorm.DB
}

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
	return v.db.Create(&video).Error
}

// GetVideoListByLatestTime gets videos for feed.
func (v *Video) GetVideoListByLatestTime(latestTime int64) ([]*model.Video, error) {
	videos := make([]*model.Video, consts.VideosLimit)
	if err := v.db.Where("create_date < ?", latestTime).
		Order("create_date desc").
		Limit(consts.VideosLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideoListByAuthorId gets videos by userId of author.
func (v *Video) GetVideoListByAuthorId(AuthorId int64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	if err := v.db.Where(&model.Video{AuthorId: AuthorId}).Order("create_date desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// GetVideoByVideoId gets video by videoId
func (v *Video) GetVideoByVideoId(vid int64) (*model.Video, error) {
	var video *model.Video
	if err := v.db.Model(model.Video{}).
		Where("id = ?", vid).First(&video).Error; err != nil {
		return nil, err
	}
	return video, nil
}

// BatchGetVideoByVideoId gets video list by videoId list.
func (v *Video) BatchGetVideoByVideoId(vidList []int64) ([]*model.Video, error) {
	if vidList == nil {
		return nil, nil
	}
	vl := make([]*model.Video, len(vidList))
	for _, vid := range vidList {
		v, err := v.GetVideoByVideoId(vid)
		if err != nil {
			return nil, err
		}
		vl = append(vl, v)
	}
	return vl, nil
}
