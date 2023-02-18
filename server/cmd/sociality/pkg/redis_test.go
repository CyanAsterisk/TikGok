package pkg

import (
	"context"
	"github.com/bytedance/sonic"
	"testing"
	"time"

	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/sociality"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
)

func TestFollowLifeCycle(t *testing.T) {
	ctx := context.Background()
	cleanUpFunc, client, err := test.RunWithRedisInDocker(consts.RedisSocialClientDB, t)
	defer cleanUpFunc()
	if err != nil {
		t.Fatal(err)
	}

	manager := NewRedisManager(client)

	aid1 := int64(100001)
	aid2 := int64(100002)
	aid3 := int64(100003)

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "action",
			op: func() (string, error) {
				time.Sleep(1 * time.Second)
				err = manager.Action(ctx, &sociality.DouyinRelationActionRequest{
					UserId:     aid3,
					ToUserId:   aid1,
					ActionType: consts.IsFollow,
				})
				err = manager.Action(ctx, &sociality.DouyinRelationActionRequest{
					UserId:     aid2,
					ToUserId:   aid1,
					ActionType: consts.IsFollow,
				})
				err = manager.Action(ctx, &sociality.DouyinRelationActionRequest{
					UserId:     aid3,
					ToUserId:   aid2,
					ActionType: consts.IsFollow,
				})
				err = manager.Action(ctx, &sociality.DouyinRelationActionRequest{
					UserId:     aid2,
					ToUserId:   aid3,
					ActionType: consts.IsFollow,
				})
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get follow id list",
			op: func() (string, error) {
				list, err := manager.List(ctx, aid2, consts.FollowList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100001,100003]",
		},
		{
			name: "get follower id list",
			op: func() (string, error) {
				list, err := manager.List(ctx, aid2, consts.FollowerList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100003]",
		},
		{
			name: "get friend id list",
			op: func() (string, error) {
				list, err := manager.List(ctx, aid2, consts.FriendsList)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "[100003]",
		},
		{
			name: "get follow num",
			op: func() (string, error) {
				count, err := manager.Count(ctx, aid2, consts.FollowCount)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(count)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "2",
		},
		{
			name: "get follower num",
			op: func() (string, error) {
				count, err := manager.Count(ctx, aid2, consts.FollowerCount)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(count)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: "1",
		},
		{
			name: "check before unfollow",
			op: func() (string, error) {
				check, err := manager.Check(ctx, aid2, aid1)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(check)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `true`,
		},
		{
			name: "unfollow",
			op: func() (string, error) {
				err = manager.Action(ctx, &sociality.DouyinRelationActionRequest{
					UserId:     aid2,
					ToUserId:   aid1,
					ActionType: consts.IsNotFollow,
				})
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "check after unfollow",
			op: func() (string, error) {
				check, err := manager.Check(ctx, aid2, aid1)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(check)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `false`,
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
