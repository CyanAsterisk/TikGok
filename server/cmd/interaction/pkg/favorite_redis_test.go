package pkg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestFavorite(t *testing.T) {
	c := context.Background()
	cleanFunc, client, err := test.RunWithRedisInDocker(consts.RedisFavoriteClientDB, t)
	defer cleanFunc()
	if err != nil {
		t.Fatal(err)
	}
	manager := NewFavoriteRedisManager(client)
	favList := make([]*model.Favorite, 0)
	timeStamp := int64(1676323214)
	for i := int64(0); i < 10; i++ {
		f := &model.Favorite{
			UserId:     100000 + i%3,
			VideoId:    200000 + i%5,
			ActionType: consts.IsLike,
			CreateDate: time.Unix(timeStamp+i, 0),
		}
		favList = append(favList, f)

	}
	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create favorite",
			op: func() (string, error) {
				time.Sleep(1 * time.Second) // wait redis docker completely start
				for _, f := range favList {
					if err = manager.Like(c, f.UserId, f.VideoId, f.CreateDate.Unix()); err != nil {
						if err != nil {
							return "", err
						}
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "check favorite",
			op: func() (string, error) {
				flag1, err := manager.Check(c, favList[0].UserId, favList[0].VideoId)
				if err != nil {
					return "", err
				}
				flag2, err := manager.Check(c, favList[0].UserId, favList[2].VideoId)
				if err != nil {
					return "", err
				}

				return fmt.Sprintf("%v %v", flag1, flag2), err
			},
			wantErr:    false,
			wantResult: `true false`,
		},
		{
			name: "get favorite video id list by userid",
			op: func() (string, error) {
				list, err := manager.GetFavoriteVideoIdListByUserId(c, favList[0].UserId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[200004,200001,200003,200000]`,
		},
		{
			name: "get favorite video count by userid",
			op: func() (string, error) {
				count, err := manager.GetFavoriteVideoCountByUserId(c, favList[0].UserId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(count)
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `4`,
		},
		{
			name: "get favorite count by video id",
			op: func() (string, error) {
				count, err := manager.GetFavoriteCountByVideoId(c, favList[0].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(count)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `2`,
		},
		{
			name: "cancel favorite",
			op: func() (string, error) {
				return "", manager.Unlike(c, favList[0].UserId, favList[0].VideoId)
			},
			wantErr: false,
		},
		{
			name: "get favorite user id list by video id after cancel favorite",
			op: func() (string, error) {
				list, err := manager.GetFavoriteCountByVideoId(c, favList[0].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `1`,
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
