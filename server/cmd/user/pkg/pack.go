package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
)

func PackUser(user *model.User, info *base.SocialInfo) *base.User {
	if user == nil {
		return nil
	}
	return &base.User{
		Id:            user.ID,
		Name:          user.Username,
		FollowCount:   info.FollowCount,
		FollowerCount: info.FollowerCount,
		IsFollow:      info.IsFollow,
	}
}

// PackUsers packs users, please make sure len(users) == len(infoList).
func PackUsers(userList []*model.User, infoList []*base.SocialInfo) []*base.User {
	if userList == nil {
		return nil
	}
	res := make([]*base.User, len(userList))
	for i, u := range userList {
		res = append(res, PackUser(u, infoList[i]))
	}
	return res
}

func PackFriendUser(user *model.User, info *base.SocialInfo, msg *chat.LatestMsg) *base.FriendUser {
	if user == nil {
		return nil
	}
	return &base.FriendUser{
		Id:            user.ID,
		Name:          user.Username,
		FollowCount:   info.FollowCount,
		FollowerCount: info.FollowerCount,
		IsFollow:      info.IsFollow,
		Message:       msg.Message,
		MsgType:       msg.MsgType,
	}
}

func PackFriendUsers(userList []*model.User, infoList []*base.SocialInfo, msgList []*chat.LatestMsg) []*base.FriendUser {
	if userList == nil {
		return nil
	}
	res := make([]*base.FriendUser, len(userList))
	for i, u := range userList {
		res = append(res, PackFriendUser(u, infoList[i], msgList[i]))
	}
	return res
}
