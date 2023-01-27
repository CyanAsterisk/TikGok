package tools

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// User model to idl
func User(u *model.User) *base.User {
	if u == nil {
		return nil
	}
	return &base.User{
		Id:   u.ID,
		Name: u.Username,
	}
}

// Users pack list of user info
func Users(users []*model.User) []*base.User {
	res := make([]*base.User, 0)
	for _, mu := range users {
		if uu := User(mu); uu != nil {
			res = append(res, uu)
		}
	}
	return res
}
