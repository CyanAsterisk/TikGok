package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
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
