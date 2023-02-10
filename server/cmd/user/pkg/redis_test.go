package pkg

import (
	"context"
	"os"
	"testing"

	"github.com/CyanAsterisk/TikGok/server/cmd/user/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestUserLifecycle(t *testing.T) {
	c := context.Background()
	manager := newManager(c, t)

	aid1 := int64(1024)
	aid2 := int64(2048)

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create account1",
			op: func() (string, error) {
				err := manager.CreateUser(c, &model.User{
					ID:       aid1,
					Username: "account1",
					Password: "12345",
				})
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate create account1",
			op: func() (string, error) {
				err := manager.CreateUser(c, &model.User{
					ID:       aid1,
					Username: "account1",
					Password: "12345",
				})
				return "", err
			},
			wantErr: true,
		},
		{
			name: "create account2",
			op: func() (string, error) {
				err := manager.CreateUser(c, &model.User{
					ID:       aid2,
					Username: "account2",
					Password: "666666",
				})
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get user by id",
			op: func() (string, error) {
				user, err := manager.GetUserById(c, aid1)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(user)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"ID":1024,"Username":"account1","Password":"12345"}`,
		},
		{
			name: "batch get user by id",
			op: func() (string, error) {
				userList, err := manager.BatchGetUserById(c, []int64{aid1, aid2})
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(userList)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":1024,"Username":"account1","Password":"12345"},{"ID":2048,"Username":"account2","Password":"666666"}]`,
		},
		{
			name: "delete user by id",
			op: func() (string, error) {
				err := manager.DeleteUser(c, aid1)
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get deleted user",
			op: func() (string, error) {
				user, err := manager.GetUserById(c, aid1)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(user)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr: true,
		},
	}

	for _, cc := range cases {
		result, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:want error;got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s:operation failed: %v", cc.name, err)
		}
		if result != cc.wantResult {
			t.Errorf("%s:result err: want %s,got %s", cc.name, cc.wantResult, result)
		}
	}
}

func newManager(c context.Context, t *testing.T) *RedisManager {
	rc, err := test.NewRedisClient(c, consts.RedisUserClientDB)
	if err != nil {
		t.Fatalf("cannot create redis client: %v", err)
	}
	return NewRedisManger(rc)
}

func TestMain(m *testing.M) {
	os.Exit(test.RunWithRedisInDocker(m))
}
