package pkg

import (
	"context"
	"errors"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisFollowingClient *redis.Client
	RedisFollowerClient  *redis.Client
}

func (r *RedisManager) Action(ctx context.Context, req *sociality.DouyinRelationActionRequest) error {
	toUserIdStr := strconv.Itoa(int(req.ToUserId))
	userIdStr := strconv.Itoa(int(req.UserId))
	flPipe := r.RedisFollowingClient.TxPipeline()
	fePipe := r.RedisFollowerClient.TxPipeline()
	flag, err := fePipe.SIsMember(ctx, toUserIdStr, req.UserId).Result()
	if err != nil {
		return err
	}
	if flag {
		// Already a fan, unfollow
		if err = fePipe.SRem(ctx, toUserIdStr, req.UserId).Err(); err != nil {
			return err
		}
		if err = flPipe.SRem(ctx, userIdStr, req.ToUserId).Err(); err != nil {
			return err
		}
		return nil
	}
	if err = fePipe.SAdd(ctx, toUserIdStr, req.UserId).Err(); err != nil {
		return err
	}
	if err = flPipe.SAdd(ctx, userIdStr, req.ToUserId).Err(); err != nil {
		return err
	}
	if _, err = flPipe.Exec(ctx); err != nil {
		return err
	}
	if _, err = fePipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RedisManager) Check(ctx context.Context, uid, toUid int64) (bool, error) {
	toUserIdStr := strconv.Itoa(int(toUid))
	userIdStr := strconv.Itoa(int(uid))
	if flag1, err := r.RedisFollowingClient.SIsMember(ctx, userIdStr, toUid).Result(); err != nil || !flag1 {
		return false, err
	}
	if flag2, err := r.RedisFollowerClient.SIsMember(ctx, toUserIdStr, uid).Result(); err != nil || !flag2 {
		return false, err
	}
	return true, nil
}

func (r *RedisManager) Count(ctx context.Context, uid int64, option int8) (int64, error) {
	userIdStr := strconv.Itoa(int(uid))
	switch option {
	case consts.FollowCount:
		cnt, err := r.RedisFollowingClient.SCard(ctx, userIdStr).Result()
		if err != nil {
			return 0, err
		}
		return cnt, nil
	case consts.FollowerCount:
		cnt, err := r.RedisFollowerClient.SCard(ctx, userIdStr).Result()
		if err != nil {
			return 0, err
		}
		return cnt, nil
	default:
		return 0, errors.New("invalid option")
	}
}

func (r *RedisManager) List(ctx context.Context, uid int64, option int8) ([]int64, error) {
	userIdStr := strconv.Itoa(int(uid))
	switch option {
	case consts.FollowerList:
		args := r.RedisFollowerClient.SMembers(ctx, userIdStr).Args()
		list, err := redis.NewIntSliceCmd(ctx, args).Result()
		if err != nil {
			return nil, err
		}
		return list, nil
	case consts.FollowList:
		args := r.RedisFollowingClient.SMembers(ctx, userIdStr).Args()
		list, err := redis.NewIntSliceCmd(ctx, args).Result()
		if err != nil {
			return nil, err
		}
		return list, nil
	case consts.FriendsList:
		args1 := r.RedisFollowerClient.SMembers(ctx, userIdStr).Args()
		list1, err := redis.NewIntSliceCmd(ctx, args1).Result()
		if err != nil {
			return nil, err
		}
		args2 := r.RedisFollowingClient.SMembers(ctx, userIdStr).Args()
		list2, err := redis.NewIntSliceCmd(ctx, args2).Result()
		if err != nil {
			return nil, err
		}
		list := tools.SimpleGeneric(list1, list2)
		return list, nil
	default:
		return nil, errors.New("invalid option")
	}
}
