package pkg

import (
	"context"
	"strconv"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
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
	//TODO implement me
	panic("implement me")
}

func (r *RedisManager) Count(ctx context.Context, uid int64, option int8) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisManager) List(ctx context.Context, uid int64, option int8) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}
