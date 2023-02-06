package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/go-redis/redis/v8"
)

type RedisManager struct {
	RedisSentClient    *redis.Client
	RedisReceiveClient *redis.Client
}

func (r *RedisManager) Action(ctx context.Context, request *chat.DouyinMessageActionRequest) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisManager) GetMessages(uid int64, toUid int64) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisManager) GetLatestMessage(uid int64, toUid int64) (*model.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisManager) BatchGetLatestMessage(uid int64, toUid []int64) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}
