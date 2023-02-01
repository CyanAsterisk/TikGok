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

func FUser(fu *kbase.FriendUser) *hbase.FriendUser {
	return &hbase.FriendUser{
		ID:            fu.Id,
		Name:          fu.Name,
		FollowCount:   fu.FollowCount,
		FollowerCount: fu.FollowerCount,
		IsFollow:      fu.IsFollow,
		Message:       fu.Message,
		MsgType:       fu.MsgType,
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

func Message(m *kbase.Message) *hbase.Message {
	return &hbase.Message{
		ID:         m.Id,
		ToUserID:   m.ToUserId,
		FromUserID: m.Id,
		Content:    m.Content,
		CreateTime: m.CreateTime,
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

func FUsers(kfus []*kbase.FriendUser) []*hbase.FriendUser {
	hfus := make([]*hbase.FriendUser, 0)
	for _, kfu := range kfus {
		if hfu := FUser(kfu); hfu != nil {
			hfus = append(hfus, hfu)
		}
	}
	return hfus
}

func Messages(kms []*kbase.Message) []*hbase.Message {
	hms := make([]*hbase.Message, 0)
	for _, km := range kms {
		if hm := Message(km); hm != nil {
			hms = append(hms, hm)
		}
	}
	return hms
}
