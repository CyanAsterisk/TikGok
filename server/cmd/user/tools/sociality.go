package tools

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
)

type SocialManager struct {
	SocialService socialityservice.Client
}

func (s *SocialManager) GetFollowerCount(ctx context.Context, userId int64) (count int64, err error) {
	res, err := s.SocialService.GetFollowerCount(ctx, &sociality.DouyinGetFollowerCountRequest{UserId: userId})
	if err != nil {
		return 0, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.SocialityServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Count, nil
}

func (s *SocialManager) GetFollowingCount(ctx context.Context, userId int64) (count int64, err error) {
	res, err := s.SocialService.GetFollowingCount(ctx, &sociality.DouyinGetFollowingCountRequest{UserId: userId})
	if err != nil {
		return 0, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.SocialityServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Count, nil
}

func (s *SocialManager) CheckFollow(ctx context.Context, userId int64, toUserId int64) (bool, error) {
	res, err := s.SocialService.CheckFollow(ctx, &sociality.DouyinCheckFollowRequest{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		return false, err
	}
	if res.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return false, errno.SocialityServerErr.WithMessage(res.BaseResp.StatusMsg)
	}
	return res.Check, nil
}
