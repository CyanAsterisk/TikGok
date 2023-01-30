package tools

import (
	hbase "github.com/CyanAsterisk/TikGok/server/cmd/api/biz/model/base"
	kbase "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

func User(u *kbase.User) *hbase.User {
	return &hbase.User{
		ID:            u.Id,
		Name:          u.Name,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      u.IsFollow,
	}
}

func Comment(c *kbase.Comment) *hbase.Comment {
	return &hbase.Comment{
		ID:         c.Id,
		User:       User(c.User),
		Content:    c.Content,
		CreateDate: c.CreateDate,
	}
}

func Video(v *kbase.Video) *hbase.Video {
	return &hbase.Video{
		ID:            v.Id,
		Author:        User(v.Author),
		PlayURL:       v.PlayUrl,
		CoverURL:      v.CoverUrl,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    v.IsFavorite,
		Title:         v.Title,
	}
}

func Videos(videos []*kbase.Video) []*hbase.Video {
	vs := make([]*hbase.Video, 0)
	for _, video := range videos {
		if v := Video(video); v != nil {
			vs = append(vs, v)
		}
	}
	return vs
}

func Comments(comments []*kbase.Comment) []*hbase.Comment {
	cs := make([]*hbase.Comment, 0)
	for _, comment := range comments {
		if c := Comment(comment); c != nil {
			cs = append(cs, c)
		}
	}
	return cs
}

func Users(users []*kbase.User) []*hbase.User {
	us := make([]*hbase.User, 0)
	for _, ku := range users {
		if hu := User(ku); hu != nil {
			us = append(us, hu)
		}
	}
	return us
}
