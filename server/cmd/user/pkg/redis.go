package pkg

import (
	"context"
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/go-redis/redis/v8"
)

const (
	UsernameFiled = "username"
	CryptPwdFiled = "password"
)

type RedisManager struct {
	redisClient *redis.Client
}

func NewRedisManger(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}

// GetUserById get user by userid.
func (r *RedisManager) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	uidStr := fmt.Sprintf("%d", uid)
	values, err := r.redisClient.HMGet(ctx, uidStr, UsernameFiled, CryptPwdFiled).Result()
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       uid,
		Username: values[0].(string),
		Password: values[1].(string),
	}, err
}

// BatchGetUserById get users by userid.
func (r *RedisManager) BatchGetUserById(ctx context.Context, uidList []int64) ([]*model.User, error) {
	if uidList == nil {
		return nil, nil
	}
	var userList []*model.User
	for _, uid := range uidList {
		user, err := r.GetUserById(ctx, uid)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}

// CreateUser creates a user.
func (r *RedisManager) CreateUser(ctx context.Context, user *model.User) error {
	uidStr := fmt.Sprintf("%d", user.ID)
	batchData := make(map[string]string)
	batchData[UsernameFiled] = user.Username
	batchData[CryptPwdFiled] = user.Password
	return r.redisClient.HMSet(ctx, uidStr, batchData).Err()
}
