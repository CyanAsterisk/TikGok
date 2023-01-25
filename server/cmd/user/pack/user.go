package pack

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
)

// User model to idl
func User(u *model.User) *user.User {
	if u == nil {
		return nil
	}
	return &user.User{
		Id:   u.ID,
		Name: u.Username,
	}
}

// Users pack list of user info
func Users(users []*model.User) []*user.User {
	res := make([]*user.User, 0)
	for _, mu := range users {
		if uu := User(mu); uu != nil {
			res = append(res, uu)
		}
	}
	return res
}
