package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// SocialityServiceImpl implements the last service interface defined in the IDL.
type SocialityServiceImpl struct {
	UserManager
}

// UserManager defines the Anti Corruption Layer
// for get user logic.
type UserManager interface {
	GetUsers(ctx context.Context, list []int64, uid int64) ([]*base.User, error)
	GetFriendUsers(ctx context.Context, list []int64, uid int64) ([]*base.FriendUser, error)
}

// Action implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) Action(_ context.Context, req *sociality.DouyinRelationActionRequest) (resp *sociality.DouyinRelationActionResponse, err error) {
	resp = new(sociality.DouyinRelationActionResponse)
	fr, err := dao.FindRecord(req.ToUserId, req.UserId)
	if err == nil && fr == nil {
		err = dao.CreateFollow(&model.Follow{
			UserId:     req.ToUserId,
			FollowerId: req.UserId,
			ActionType: req.ActionType,
		})
		if err != nil {
			klog.Error("follow action error", err)
			resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
			return resp, nil
		}
		resp.BaseResp = sTools.BuildBaseResp(nil)
		return resp, nil
	}
	if err != nil {
		klog.Error("follow error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
		return resp, nil
	}
	err = dao.UpdateFollow(req.ToUserId, req.UserId, req.ActionType)
	if err != nil {
		klog.Error("follow error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("follow action error"))
		return resp, nil
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FollowingList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) FollowingList(ctx context.Context, req *sociality.DouyinRelationFollowListRequest) (resp *sociality.DouyinRelationFollowListResponse, err error) {
	resp = new(sociality.DouyinRelationFollowListResponse)
	list, err := dao.GetFollowingIdList(req.OwnerId)
	if err != nil {
		klog.Error("get following list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get following list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetUsers(ctx, list, req.ViewerId)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCUserErr.WithMessage("get following list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FollowerList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) FollowerList(ctx context.Context, req *sociality.DouyinRelationFollowerListRequest) (resp *sociality.DouyinRelationFollowerListResponse, err error) {
	resp = new(sociality.DouyinRelationFollowerListResponse)
	list, err := dao.GetFollowerIdList(req.OwnerId)
	if err != nil {
		klog.Error("get follower list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetUsers(ctx, list, req.ViewerId)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCUserErr.WithMessage("get follower list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FriendList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) FriendList(ctx context.Context, req *sociality.DouyinRelationFriendListRequest) (resp *sociality.DouyinRelationFriendListResponse, err error) {
	resp = new(sociality.DouyinRelationFriendListResponse)
	list, err := dao.GetFriendsList(req.OwnerId)
	if err != nil {
		klog.Error("get friends list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get friends list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetFriendUsers(ctx, list, req.ViewerId)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.RPCUserErr.WithMessage("get friends list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// CheckFollow implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) CheckFollow(_ context.Context, req *sociality.DouyinCheckFollowRequest) (resp *sociality.DouyinCheckFollowResponse, err error) {
	resp = new(sociality.DouyinCheckFollowResponse)
	info, err := dao.FindRecord(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("check follow error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("check follow error"))
		return resp, nil
	}
	if info.ActionType == consts.IsFollow {
		resp.Check = true
	} else {
		resp.Check = false
	}
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// GetFollowerCount implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetFollowerCount(_ context.Context, req *sociality.DouyinGetFollowerCountRequest) (resp *sociality.DouyinGetFollowerCountResponse, err error) {
	resp = new(sociality.DouyinGetFollowerCountResponse)
	count, err := dao.GetFollowerNumsByUserId(req.UserId)
	if err != nil {
		klog.Error("get follower num error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower num error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// GetFollowingCount implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetFollowingCount(_ context.Context, req *sociality.DouyinGetFollowingCountRequest) (resp *sociality.DouyinGetFollowingCountResponse, err error) {
	resp = new(sociality.DouyinGetFollowingCountResponse)
	count, err := dao.GetFollowingNumsByUserId(req.UserId)
	if err != nil {
		klog.Error("get following num error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get following num error"))
		return resp, nil
	}
	resp.Count = count
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}
