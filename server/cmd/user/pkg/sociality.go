package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality/socialityservice"
)

type SocialManager struct {
	SocialService socialityservice.Client
}

func (s *SocialManager) GetSocialInfo(ctx context.Context, viewerId, ownerId int64) (*base.SocialInfo, error) {
	resp, err := s.SocialService.GetSocialInfo(ctx, &sociality.DouyinGetSocialInfoRequest{
		ViewerId: viewerId,
		OwnerId:  ownerId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.SocialInfo, nil
}

func (s *SocialManager) BatchGetSocialInfo(ctx context.Context, viewerId int64, ownerIdList []int64) ([]*base.SocialInfo, error) {
	resp, err := s.SocialService.BatchGetSocialInfo(ctx, &sociality.DouyinBatchGetSocialInfoRequest{
		ViewerId:    viewerId,
		OwnerIdList: ownerIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.SocialInfoList, nil
}

func (s *SocialManager) GetRelationList(ctx context.Context, viewerId, ownerId int64, option int8) ([]int64, error) {
	resp, err := s.SocialService.GetRelationIdList(ctx, &sociality.DouyinGetRelationIdListRequest{
		ViewerId: viewerId,
		OwnerId:  ownerId,
		Option:   option,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.UserIdList, nil
}
