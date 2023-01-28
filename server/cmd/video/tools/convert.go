package tools

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// Video model to idl
func Video(v *model.Video) *base.Video {
	if v == nil {
		return nil
	}
	return &base.Video{
		Id: int64(v.ID),
		Author: &base.User{
			Id:            v.Uid,
			Name:          "",
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		},
		PlayUrl:  consts.MinIOServer + v.PlayUrl,
		CoverUrl: consts.MinIOServer + v.CoverUrl,
		Title:    v.Title,
	}
}

func Videos(videos []*model.Video) []*base.Video {
	vs := make([]*base.Video, 0)
	for _, vid := range videos {
		if v := Video(vid); v != nil {
			vs = append(vs, v)
		}
	}
	return vs
}
