package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

func PackVideo(video *model.Video, author *base.User, info *base.InteractInfo) *base.Video {
	if video == nil {
		return nil
	}
	return &base.Video{
		Id:            int64(video.ID),
		Author:        author,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: info.FavoriteCount,
		CommentCount:  info.CommentCount,
		IsFavorite:    info.IsFavorite,
		Title:         video.Title,
	}
}
func PackVideos(videoList []*model.Video, authorList []*base.User, infoList []*base.InteractInfo) []*base.Video {
	if videoList == nil {
		return nil
	}
	res := make([]*base.Video, len(videoList))
	for i, v := range videoList {
		res = append(res, PackVideo(v, authorList[i], infoList[i]))
	}
	return res
}
