package tools

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
)

type SocialManager struct {
	SocialService socialityservice.Client
}

func GetFollowerCount(ctx context.Context, userId int64) (count int64, err error) {
	//TODO: add this api to social service
	return 0, nil
}

func GetFollowingCount(ctx context.Context, userId int64) (count int64, err error) {
	//TODO: add this api to social service
	return 0, nil
}

// CheckFollow checks whether userId follower toUserId.
func CheckFollow(ctx context.Context, userId int64, toUserId int64) (bool, error) {
	//TODO: add this api to social service
	return true, nil
}
