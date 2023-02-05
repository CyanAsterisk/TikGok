package pkg

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
	resp, err := s.SocialService.GetFollowerCount(ctx, &sociality.DouyinGetFollowerCountRequest{UserId: userId})
	if err != nil {
		return 0, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Count, nil
}

func (s *SocialManager) GetFollowingCount(ctx context.Context, userId int64) (count int64, err error) {
	resp, err := s.SocialService.GetFollowingCount(ctx, &sociality.DouyinGetFollowingCountRequest{UserId: userId})
	if err != nil {
		return 0, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return 0, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Count, nil
}

func (s *SocialManager) CheckFollow(ctx context.Context, userId int64, toUserId int64) (bool, error) {
	resp, err := s.SocialService.CheckFollow(ctx, &sociality.DouyinCheckFollowRequest{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		return false, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return false, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Check, nil
}

func (s *SocialManager) BatchGetFollowerCount(ctx context.Context, userIds []int64) (counts []int64, err error) {
	resp, err := s.SocialService.BatchGetFollowerCount(ctx, &sociality.DouyinBatchGetFollowerCountRequest{UserIds: userIds})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Counts, nil
}

func (s *SocialManager) BatchGetFollowingCount(ctx context.Context, userIds []int64) (counts []int64, err error) {
	resp, err := s.SocialService.BatchGetFollowingCount(ctx, &sociality.DouyinBatchGetFollowingCountRequest{UserIds: userIds})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Counts, nil
}

func (s *SocialManager) BatchCheckFollow(ctx context.Context, userId int64, toUserIds []int64) (checks []bool, err error) {
	resp, err := s.SocialService.BatchCheckFollow(ctx, &sociality.DouyinBatchCheckFollowRequest{
		UserId:    userId,
		ToUserIds: toUserIds,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.Checks, nil
}
