package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisFollowingClient *redis.Client
	RedisFollowerClient  *redis.Client
}

func (r *RedisManager) Action(ctx context.Context, request *sociality.DouyinRelationActionRequest) error {
	//TODO implement me
	panic("implement me")
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
