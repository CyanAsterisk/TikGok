package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

func PackUser(user *model.User, followerCnt int64, followingCnt int64, isFollow bool) *base.User {
	if user == nil {
		return nil
	}
	return &base.User{
		Id:            user.ID,
		Name:          user.Username,
		FollowCount:   followingCnt,
		FollowerCount: followerCnt,
		IsFollow:      isFollow,
	}
}

// PackUsers packs users, please make sure len(users) == len(followCnt) == len(followingCnt) == len(isFollow).
func PackUsers(users []*model.User, followerCnt []int64, followingCnt []int64, isFollow []bool) []*base.User {
	if users == nil {
		return nil
	}
	n := len(users)
	res := make([]*base.User, n)
	for i := 0; i < n; i++ {
		res = append(res, PackUser(users[i], followerCnt[i], followingCnt[i], isFollow[i]))
	}
	return res
}
