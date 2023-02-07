package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
)

// CreateVideo creates a new video record.
func CreateVideo(video *model.Video) error {
	return global.DB.Create(&video).Error
}

// GetVideoListByLatestTime gets videos for feed.
func GetVideoListByLatestTime(latestTime int64) ([]*model.Video, error) {
	videos := make([]*model.Video, consts.VideosLimit)
	if err := global.DB.Where("create_date < ?", latestTime).
		Order("create_date desc").
		Limit(consts.VideosLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideoListByAuthorId gets videos by userId of author.
func GetVideoListByAuthorId(AuthorId int64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	if err := global.DB.Where(&model.Video{AuthorId: AuthorId}).Order("create_date desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// GetVideoByVideoId gets video by videoId
func GetVideoByVideoId(vid int64) (*model.Video, error) {
	var video *model.Video
	if err := global.DB.Model(model.Video{}).
		Where("id = ?", vid).First(&video).Error; err != nil {
		return nil, err
	}
	return video, nil
}

// BatchGetVideoByVideoId gets video list by videoId list.
func BatchGetVideoByVideoId(vidList []int64) ([]*model.Video, error) {
	if vidList == nil {
		return nil, nil
	}
	vl := make([]*model.Video, len(vidList))
	for _, vid := range vidList {
		v, err := GetVideoByVideoId(vid)
		if err != nil {
			return nil, err
		}
		vl = append(vl, v)
	}
	return vl, nil
}
