package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/dao"
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
	//fr, err := dao.FindRecord(req.ToUserId, req.UserId)
	//if err == nil && fr == nil {
	//	err = dao.CreateFollow(&model.Follow{
	//		UserId:     req.ToUserId,
	//		FollowerId: req.UserId,
	//		ActionType: req.ActionType,
	//	})
	//	if err != nil {
	//		klog.Error("follow action error", err)
	//		resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
	//		return resp, nil
	//	}
	//	resp.BaseResp = tools.BuildBaseResp(nil)
	//	return resp, nil
	//}
	//if err != nil {
	//	klog.Error("follow error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("follow action error"))
	//	return resp, nil
	//}
	//err = dao.UpdateFollow(req.ToUserId, req.UserId, req.ActionType)
	//if err != nil {
	//	klog.Error("follow error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.InteractionServerErr.WithMessage("follow action error"))
	//	return resp, nil
	//}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

//
//// GetFollowingList implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) GetFollowingList(ctx context.Context, req *sociality.DouyinGetRelationFollowListRequest) (resp *sociality.DouyinGetRelationFollowListResponse, err error) {
//	resp = new(sociality.DouyinGetRelationFollowListResponse)
//	list, err := s.RedisManager.List(ctx, req.OwnerId, consts.FollowingList)
//	if err != nil {
//		klog.Error("get following list by redis error", err)
//		list, err = dao.GetFollowingIdList(req.OwnerId)
//		if err != nil {
//			klog.Error("get following list by mysql error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get following list error"))
//			return resp, nil
//		}
//	}
//	users, err := s.UserManager.GetUsers(ctx, list, req.ViewerId)
//	if err != nil {
//		klog.Error("get users by user manager error", err)
//		resp.BaseResp = tools.BuildBaseResp(errno.RPCUserErr.WithMessage("get following list error"))
//		return resp, nil
//	}
//	resp.UserList = users
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}
//
//// GetFollowerList implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) GetFollowerList(ctx context.Context, req *sociality.DouyinGetRelationFollowerListRequest) (resp *sociality.DouyinGetRelationFollowerListResponse, err error) {
//	resp = new(sociality.DouyinGetRelationFollowerListResponse)
//	list, err := s.RedisManager.List(ctx, req.OwnerId, consts.FollowerList)
//	if err != nil {
//		klog.Error("get follower list by redis error", err)
//		list, err = dao.GetFollowerIdList(req.OwnerId)
//		if err != nil {
//			klog.Error("get follower list by mysql error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower list error"))
//			return resp, nil
//		}
//	}
//	users, err := s.UserManager.GetUsers(ctx, list, req.ViewerId)
//	if err != nil {
//		klog.Error("get users by user manager error", err)
//		resp.BaseResp = tools.BuildBaseResp(errno.RPCUserErr.WithMessage("get follower list error"))
//		return resp, nil
//	}
//	resp.UserList = users
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}
//
//// GetFriendList implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) GetFriendList(ctx context.Context, req *sociality.DouyinGetRelationFriendListRequest) (resp *sociality.DouyinGetRelationFriendListResponse, err error) {
//	resp = new(sociality.DouyinGetRelationFriendListResponse)
//	list, err := s.RedisManager.List(ctx, req.OwnerId, consts.FriendsList)
//	if err != nil {
//		klog.Error("get friends list by redis error", err)
//		list, err = dao.GetFriendsList(req.OwnerId)
//		if err != nil {
//			klog.Error("get friends list by mysql error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get friends list error"))
//			return resp, nil
//		}
//	}
//	users, err := s.UserManager.GetFriendUsers(ctx, list, req.ViewerId)
//	if err != nil {
//		klog.Error("get users by user manager error", err)
//		resp.BaseResp = tools.BuildBaseResp(errno.RPCUserErr.WithMessage("get friends list error"))
//		return resp, nil
//	}
//	resp.UserList = users
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}
//
//// CheckFollow implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) CheckFollow(ctx context.Context, req *sociality.DouyinCheckFollowRequest) (resp *sociality.DouyinCheckFollowResponse, err error) {
//	resp = new(sociality.DouyinCheckFollowResponse)
//	flag, err := s.RedisManager.Check(ctx, req.UserId, req.ToUserId)
//	if err != nil {
//		klog.Error("check follow by redis error", err)
//		info, err := dao.FindRecord(req.UserId, req.ToUserId)
//		if err != nil {
//			klog.Error("check follow by mysql error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("check follow error"))
//			return resp, nil
//		}
//		if info == nil {
//			resp.Check = false
//		} else {
//			if info.ActionType == consts.IsFollow {
//				resp.Check = true
//			} else {
//				resp.Check = false
//			}
//		}
//		resp.BaseResp = tools.BuildBaseResp(nil)
//		return resp, nil
//	}
//	resp.Check = flag
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}
//
//// GetFollowerCount implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) GetFollowerCount(ctx context.Context, req *sociality.DouyinGetFollowerCountRequest) (resp *sociality.DouyinGetFollowerCountResponse, err error) {
//	resp = new(sociality.DouyinGetFollowerCountResponse)
//	count, err := s.RedisManager.Count(ctx, req.UserId, consts.FollowerCount)
//	if err != nil {
//		klog.Error("get follower num by redis error", err)
//		count, err = dao.GetFollowerNumsByUserId(req.UserId)
//		if err != nil {
//			klog.Error("get follower num by mysql error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get follower num error"))
//			return resp, nil
//		}
//		resp.Count = count
//		resp.BaseResp = tools.BuildBaseResp(nil)
//		return resp, nil
//	}
//	resp.Count = count
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}
//
//// GetFollowingCount implements the SocialityServiceImpl interface.
//func (s *SocialityServiceImpl) GetFollowingCount(ctx context.Context, req *sociality.DouyinGetFollowingCountRequest) (resp *sociality.DouyinGetFollowingCountResponse, err error) {
//	resp = new(sociality.DouyinGetFollowingCountResponse)
//	count, err := s.RedisManager.Count(ctx, req.UserId, consts.FollowingCount)
//	if err != nil {
//		klog.Error("get following num by redis error", err)
//		count, err = dao.GetFollowingNumsByUserId(req.UserId)
//		if err != nil {
//			klog.Error("get following num error", err)
//			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get following num error"))
//			return resp, nil
//		}
//		resp.Count = count
//		resp.BaseResp = tools.BuildBaseResp(nil)
//		return resp, nil
//	}
//	resp.Count = count
//	resp.BaseResp = tools.BuildBaseResp(nil)
//	return resp, nil
//}

// GetRelationIdList implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetRelationIdList(ctx context.Context, req *sociality.DouyinGetRelationIdListRequest) (resp *sociality.DouyinGetRelationIdListResponse, err error) {
	resp = new(sociality.DouyinGetRelationIdListResponse)
	resp.UserIdList, err = s.RedisManager.List(ctx, req.OwnerId, req.Option)
	if err != nil {
		klog.Error("get id list by redis error", err)
		resp.UserIdList, err = dao.GetFollowingIdList(req.OwnerId)
		if err != nil {
			klog.Error("get id list by mysql error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.SocialityServerErr.WithMessage("get id list error"))
			return resp, nil
		}
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// GetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) GetSocialInfo(ctx context.Context, req *sociality.DouyinGetSocialInfoRequest) (resp *sociality.DouyinGetSocialInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// BatchGetSocialInfo implements the SocialityServiceImpl interface.
func (s *SocialityServiceImpl) BatchGetSocialInfo(ctx context.Context, req *sociality.DouyinBatchGetSocialInfoRequest) (resp *sociality.DouyinBatchGetSocialInfoResponse, err error) {
	// TODO: Your code here...
	return
}
