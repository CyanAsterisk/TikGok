package main

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/dao"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// SocialityServiceImpl implements the last service interface defined in the IDL.
type SocialityServiceImpl struct {
	UserManager
	RedisManager
	Publisher
	Subscriber
}

// UserManager defines the Anti Corruption Layer
// for get user logic.
type UserManager interface {
	GetUsers(ctx context.Context, list []int64, uid int64) ([]*base.User, error)
	GetFriendUsers(ctx context.Context, list []int64, uid int64) ([]*base.FriendUser, error)
}

// RedisManager defines the redis interface.
type RedisManager interface {
	Action(context.Context, *sociality.DouyinRelationActionRequest) error
	Check(ctx context.Context, uid, toUid int64) (bool, error)
	Count(ctx context.Context, uid int64, option int8) (int64, error)
	List(ctx context.Context, uid int64, option int8) ([]int64, error)
}

// Publisher defines the publisher interface.
type Publisher interface {
	Publish(context.Context, *sociality.DouyinRelationActionRequest) error
}

// Subscriber defines a car update subscriber.
type Subscriber interface {
	Subscribe(context.Context) (ch chan *sociality.DouyinRelationActionRequest, cleanUp func(), err error)
}

// Action implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) Action(ctx context.Context, req *sociality.DouyinRelationActionRequest) (resp *sociality.DouyinRelationActionResponse, err error) {
	resp = new(sociality.DouyinRelationActionResponse)
	err = s.Publisher.Publish(ctx, req)
	if err != nil {
		klog.Error("action publish error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
		return resp, nil
	}
	err = s.RedisManager.Action(ctx, req)
	if err != nil {
		klog.Error("redis action error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetRelationIdList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetRelationIdList(ctx context.Context, req *sociality.DouyinGetRelationIdListRequest) (resp *sociality.DouyinGetRelationIdListResponse, err error) {
	resp = new(sociality.DouyinGetRelationIdListResponse)
	resp.UserIdList, err = s.RedisManager.List(ctx, req.OwnerId, req.Option)
	if err != nil {
		klog.Error("get id list by redis error", err)
		if req.Option == consts.FollowList {
			resp.UserIdList, err = dao.GetFollowIdList(req.OwnerId)
		} else if req.Option == consts.FollowerList {
			resp.UserIdList, err = dao.GetFollowerIdList(req.OwnerId)
		} else if req.Option == consts.FriendsList {
			resp.UserIdList, err = dao.GetFriendsList(req.OwnerId)
		} else {
			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("wrong option"))
			return resp, nil
		}
		if err != nil {
			klog.Error("get relation id list error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get id list error"))
			return resp, nil
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetSocialInfo(ctx context.Context, req *sociality.DouyinGetSocialInfoRequest) (resp *sociality.DouyinGetSocialInfoResponse, err error) {
	resp = new(sociality.DouyinGetSocialInfoResponse)
	if resp.SocialInfo, err = s.getSocialInfo(ctx, req.ViewerId, req.OwnerId); err != nil {
		klog.Error("get social info err", err)
		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get social info err"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) BatchGetSocialInfo(ctx context.Context, req *sociality.DouyinBatchGetSocialInfoRequest) (resp *sociality.DouyinBatchGetSocialInfoResponse, err error) {
	resp = new(sociality.DouyinBatchGetSocialInfoResponse)
	for _, oid := range req.OwnerIdList {
		info, err := s.getSocialInfo(ctx, oid, req.ViewerId)
		if err != nil {
			klog.Error("get social info err", err)
			resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr)
			return resp, nil
		}
		resp.SocialInfoList = append(resp.SocialInfoList, info)
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

func (s *SocialityServiceImpl) getSocialInfo(ctx context.Context, viewerId int64, ownerId int64) (info *base.SocialInfo, err error) {
	info = new(base.SocialInfo)
	if info.FollowCount, err = s.RedisManager.Count(ctx, ownerId, consts.FollowCount); err != nil {
		klog.Error("get follow count by redis err", err)
		info.FollowCount, err = dao.GetFollowingNumsByUserId(ownerId)
		if err != nil {
			return nil, err
		}
	}
	if info.FollowerCount, err = s.RedisManager.Count(ctx, ownerId, consts.FollowerCount); err != nil {
		klog.Error("get follower count by redis err", err)
		info.FollowCount, err = dao.GetFollowerNumsByUserId(ownerId)
		if err != nil {
			return nil, err
		}
	}
	if info.IsFollow, err = s.RedisManager.Check(ctx, viewerId, ownerId); err != nil {
		klog.Error("check follow by redis err", err)
		record, err := dao.FindRecord(ownerId, viewerId)
		if err != nil {
			return nil, err
		}
		if record != nil && record.ActionType == consts.IsFollow {
			info.IsFollow = true
		} else {
			info.IsFollow = false
		}
	}
	return info, err
}
