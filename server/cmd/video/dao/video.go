package dao

import (
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/video/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
)

// CreateVideo creates a new video record.
func CreateVideo(video *model.Video) error {
	return global.DB.Create(&video).Error
}

// GetVideosByLatestTime gets videos for feed.
func GetVideosByLatestTime(latestTime int64) ([]*model.Video, error) {
	videos := make([]*model.Video, consts.VideosLimit)
	if latestTime <= 0 {
		latestTime = time.Now().UnixNano() / 1e6
	}
	if err := global.DB.Where("updated_at < ?", time.Unix(0, latestTime*1e6).Local()).
		Order("updated_at desc").
		Limit(consts.VideosLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideosByUserId gets videos by userId
func GetVideosByUserId(uid int64) ([]*model.Video, error) {
	res := make([]*model.Video, 0)
	if err := global.DB.Where(&model.Video{Uid: uid}).Order("updated_at desc").Find(&res).Error; err != nil {
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
