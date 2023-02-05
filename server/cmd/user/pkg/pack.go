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
func PackUsers(users []*model.User, infoList []*base.SocialInfo) []*base.User {
	if users == nil {
		return nil
	}
	n := len(users)
	res := make([]*base.User, n)
	for i := 0; i < n; i++ {
		res = append(res, PackUser(users[i], infoList[i]))
	}
	return res
}
