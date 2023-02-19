package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

func PackUser(user *model.User, socialInfo *base.SocialInfo, interInfo *base.UserInteractInfo) *base.User {
	if user == nil {
		return nil
	}
	return &base.User{
		Id:              user.ID,
		Name:            user.Username,
		FollowCount:     socialInfo.FollowCount,
		FollowerCount:   socialInfo.FollowerCount,
		IsFollow:        socialInfo.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackGroundImage,
		Signature:       user.Signature,
		TotalFavorited:  interInfo.TotalFavorited,
		WorkCount:       interInfo.WorkCount,
		FavoriteCount:   interInfo.FavoriteCount,
	}
}

// PackUsers packs users, please make sure len(users) == len(infoList).
func PackUsers(userList []*model.User, socialInfoList []*base.SocialInfo, interactInfoList []*base.UserInteractInfo) []*base.User {
	if userList == nil {
		return nil
	}
	res := make([]*base.User, 0)
	for i, u := range userList {
		res = append(res, PackUser(u, socialInfoList[i], interactInfoList[i]))
	}
	return res
}

func PackFriendUser(user *model.User, socialInfo *base.SocialInfo, interactInfo *base.UserInteractInfo, msg *base.LatestMsg) *base.FriendUser {
	if user == nil {
		return nil
	}
	return &base.FriendUser{
		Id:              user.ID,
		Name:            user.Username,
		FollowCount:     socialInfo.FollowCount,
		FollowerCount:   socialInfo.FollowerCount,
		IsFollow:        socialInfo.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackGroundImage,
		Signature:       user.Signature,
		TotalFavorited:  interactInfo.TotalFavorited,
		WorkCount:       interactInfo.WorkCount,
		FavoriteCount:   interactInfo.FavoriteCount,
		Message:         msg.Message,
		MsgType:         msg.MsgType,
	}
}

func PackFriendUsers(userList []*model.User, socialInfoList []*base.SocialInfo, interactInfoList []*base.UserInteractInfo, msgList []*base.LatestMsg) []*base.FriendUser {
	if userList == nil {
		return nil
	}
	res := make([]*base.FriendUser, 0)
	for i, u := range userList {
		res = append(res, PackFriendUser(u, socialInfoList[i], interactInfoList[i], msgList[i]))
	}
	return res
}
