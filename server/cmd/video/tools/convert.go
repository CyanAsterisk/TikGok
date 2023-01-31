package tools

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/video/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// Video model to idl
func Video(v *model.Video) *base.Video {
	return &base.Video{
		Id:       int64(v.ID),
		PlayUrl:  v.PlayUrl,
		CoverUrl: v.CoverUrl,
		Title:    v.Title,
	}

}

// Videos model to idl
func Videos(videos []*model.Video) []*base.Video {
	vs := make([]*base.Video, 0)
	for _, vid := range videos {
		v := Video(vid)
		vs = append(vs, v)
	}
	return vs
}
