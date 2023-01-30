package main

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/model"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	sTools "github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

// SocialityServiceImpl implements the last service interface defined in the IDL.
type SocialityServiceImpl struct {
	UserManager
}

// UserManager defines the Anti Corruption Layer
// for get user logic.
type UserManager interface {
	GetUsers([]int64) ([]*base.User, error)
}

// Action implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) Action(_ context.Context, req *sociality.DouyinRelationActionRequest) (resp *sociality.DouyinRelationActionResponse, err error) {
	resp = new(sociality.DouyinRelationActionResponse)
	_, err = dao.FindRecord(req.ToUserId, req.UserId)
	if err == gorm.ErrRecordNotFound {
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
func (s *SocialityServiceImpl) FollowingList(_ context.Context, req *sociality.DouyinRelationFollowListRequest) (resp *sociality.DouyinRelationFollowListResponse, err error) {
	resp = new(sociality.DouyinRelationFollowListResponse)
	list, err := dao.GetFollowingIdList(req.UserId)
	if err != nil {
		klog.Error("get following list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get following list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetUsers(list)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get following list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FollowerList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) FollowerList(_ context.Context, req *sociality.DouyinRelationFollowerListRequest) (resp *sociality.DouyinRelationFollowerListResponse, err error) {
	resp = new(sociality.DouyinRelationFollowerListResponse)
	list, err := dao.GetFollowerIdList(req.UserId)
	if err != nil {
		klog.Error("get follower list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetUsers(list)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get follower list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}

// FriendList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) FriendList(_ context.Context, req *sociality.DouyinRelationFriendListRequest) (resp *sociality.DouyinRelationFriendListResponse, err error) {
	resp = new(sociality.DouyinRelationFriendListResponse)
	list, err := dao.GetFriendsList(req.UserId)
	if err != nil {
		klog.Error("get friends list error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get friends list error"))
		return resp, nil
	}
	users, err := s.UserManager.GetUsers(list)
	if err != nil {
		klog.Error("get users by user manager error", err)
		resp.BaseResp = sTools.BuildBaseResp(errno.InteractionServerErr.WithMessage("get friends list error"))
		return resp, nil
	}
	resp.UserList = users
	resp.BaseResp = sTools.BuildBaseResp(nil)
	return resp, nil
}
