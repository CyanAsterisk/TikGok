package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisClient *redis.Client
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{RedisClient: client}
}

func (r *RedisManager) Action(ctx context.Context, req *sociality.DouyinRelationActionRequest) error {
	toUserIdStr := fmt.Sprintf("%d", req.ToUserId)
	userIdStr := fmt.Sprintf("%d", req.UserId)
	pl := r.RedisClient.TxPipeline()
	if req.ActionType == consts.IsFollow {
		if err := pl.SAdd(ctx, userIdStr+consts.RedisFollowSuffix, req.ToUserId).Err(); err != nil {
			return err
		}
		if err := pl.SAdd(ctx, toUserIdStr+consts.RedisFollowerSuffix, req.UserId).Err(); err != nil {
			return err
		}
	} else if req.ActionType == consts.IsNotFollow {
		if err := pl.SRem(ctx, userIdStr+consts.RedisFollowSuffix, req.ToUserId).Err(); err != nil {
			return err
		}
		if err := pl.SRem(ctx, toUserIdStr+consts.RedisFollowerSuffix, req.UserId).Err(); err != nil {
			return err
		}
	} else {
		return errno.SocialityServerErr.WithMessage("invalid action type")
	}
	_, err := pl.Exec(ctx)
	return err
}

func (r *RedisManager) Check(ctx context.Context, uid, toUid int64) (bool, error) {
	toUserIdStr := fmt.Sprintf("%d", toUid)
	userIdStr := fmt.Sprintf("%d", uid)
	flag1, err := r.RedisClient.SIsMember(ctx, userIdStr+consts.RedisFollowSuffix, toUid).Result()
	if err != nil {
		return false, err
	}
	flag2, err := r.RedisClient.SIsMember(ctx, toUserIdStr+consts.RedisFollowerSuffix, uid).Result()
	if err != nil {
		return false, err
	}
	if flag1 != flag2 {
		return false, errno.SocialityServerErr.WithMessage("dirty data in redis")
	}
	return flag1, nil
}

func (r *RedisManager) Count(ctx context.Context, uid int64, option int8) (int64, error) {
	userIdStr := fmt.Sprintf("%d", uid)
	switch option {
	case consts.FollowCount:
		cnt, err := r.RedisClient.SCard(ctx, userIdStr+consts.RedisFollowSuffix).Result()
		if err != nil {
			return 0, err
		}
		return cnt, nil
	case consts.FollowerCount:
		cnt, err := r.RedisClient.SCard(ctx, userIdStr+consts.RedisFollowerSuffix).Result()
		if err != nil {
			return 0, err
		}
		return cnt, nil
	default:
		return 0, errors.New("invalid option")
	}
}

func (r *RedisManager) List(ctx context.Context, uid int64, option int8) (list []int64, err error) {
	userIdStr := fmt.Sprintf("%d", uid)
	list = make([]int64, 0)
	switch option {
	case consts.FollowerList:
		err = r.RedisClient.SMembers(ctx, userIdStr+consts.RedisFollowerSuffix).ScanSlice(&list)
	case consts.FollowList:
		err = r.RedisClient.SMembers(ctx, userIdStr+consts.RedisFollowSuffix).ScanSlice(&list)
	case consts.FriendsList:
		err = r.RedisClient.SInter(ctx, userIdStr+consts.RedisFollowSuffix, userIdStr+consts.RedisFollowerSuffix).ScanSlice(&list)
	default:
		return nil, errors.New("invalid option")
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}
