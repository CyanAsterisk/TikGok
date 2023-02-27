package dao

import (
	"testing"

	"github.com/CyanAsterisk/TikGok/server/cmd/sociality/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestFollowLifeCycle(t *testing.T) {
	cleanUpFunc, db, err := test.RunWithMySQLInDocker(t)
	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	dao := NewFollow(db)

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
			name: "create follow",
			op: func() (string, error) {
				err = dao.CreateFollow(&model.Follow{
					UserId:     aid1,
					FollowerId: aid2,
					ActionType: 1,
				})
				err = dao.CreateFollow(&model.Follow{
					UserId:     aid1,
					FollowerId: aid3,
					ActionType: 1,
				})
				err = dao.CreateFollow(&model.Follow{
					UserId:     aid2,
					FollowerId: aid3,
					ActionType: 1,
				})
				err = dao.CreateFollow(&model.Follow{
					UserId:     aid3,
					FollowerId: aid2,
					ActionType: 1,
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
				list, err := dao.GetFollowIdList(aid2)
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
				list, err := dao.GetFollowerIdList(aid2)
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
				list, err := dao.GetFriendsList(aid2)
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
				count, err := dao.GetFollowNumsByUserId(aid2)
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
				count, err := dao.GetFollowerNumsByUserId(aid2)
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
			name: "get record before unfollow",
			op: func() (string, error) {
				record, err := dao.FindRecord(aid1, aid2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(record)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"UserId":100001,"FollowerId":100002,"ActionType":1}`,
		},
		{
			name: "unfollow",
			op: func() (string, error) {
				err := dao.UpdateFollow(aid1, aid2, consts.IsNotFollow)
				if err != nil {
					return "", err
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get record after unfollow",
			op: func() (string, error) {
				record, err := dao.FindRecord(aid1, aid2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(record)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"UserId":100001,"FollowerId":100002,"ActionType":2}`,
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
