package pkg

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
)

var (
	ErrNoSuchUser = errors.New("no such user")
	ErrUserExist  = errors.New("user already exist")
)

type RedisManager struct {
	redisClient *redis.Client
}

func NewRedisManager(client *redis.Client) *RedisManager {
	return &RedisManager{redisClient: client}
}

// GetUserById get user by userid.
func (r *RedisManager) GetUserById(ctx context.Context, uid int64) (*model.User, error) {
	uidStr := fmt.Sprintf("%d", uid)
	userJson, err := r.redisClient.Get(ctx, uidStr).Result()
	if err != nil {
		return nil, err
	}
	var user model.User
	err = sonic.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// BatchGetUserById get users by userid.
func (r *RedisManager) BatchGetUserById(ctx context.Context, uidList []int64) ([]*model.User, error) {
	if uidList == nil {
		return nil, nil
	}
	length := len(uidList)
	userList := make([]*model.User, length)
	var wg sync.WaitGroup
	wg.Add(length)
	for i := 0; i < length; i++ {
		go func(idx int) {
			defer wg.Done()
			userList[idx], _ = r.GetUserById(ctx, uidList[idx])
		}(i)
	}
	wg.Wait()
	return userList, nil
}

// CreateUser creates a user.
func (r *RedisManager) CreateUser(ctx context.Context, user *model.User) error {
	uidStr := fmt.Sprintf("%d", user.ID)
	err := r.redisClient.Get(ctx, uidStr).Err()
	if err == nil {
		return ErrUserExist
	}
	if err != redis.Nil {
		return err
	}
	userJson, err := sonic.Marshal(user)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, uidStr, userJson, 0).Err()
}

// DeleteUser delete a user by userId.
func (r *RedisManager) DeleteUser(ctx context.Context, userId int64) error {
	uidStr := fmt.Sprintf("%d", userId)
	err := r.redisClient.Get(ctx, uidStr).Err()
	if err == redis.Nil {
		return ErrNoSuchUser
	}
	if err != nil {
		return err
	}
	return r.redisClient.Del(ctx, uidStr).Err()
}
