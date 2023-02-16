package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
	"testing"
	"time"
)

func TestFavorite(t *testing.T) {
	cleanUpFunc, db, err := test.RunWithMySQLInDocker(t)
	defer cleanUpFunc()
	if err != nil {
		t.Fatal(err)
	}
	dao := NewFavorite(db)
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
				for _, f := range favList {
					if err = dao.CreateFavorite(f); err != nil {
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
			name: "duplicate create favorite",
			op: func() (string, error) {
				err = dao.CreateFavorite(favList[0])
				return "", err
			},
			wantErr: true,
		},
		{
			name: "get favorite info",
			op: func() (string, error) {
				fav, err := dao.GetFavoriteInfo(favList[0].UserId, favList[1].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(fav)
				return string(result), err
			},
			wantErr:    false,
			wantResult: `{"UserId":100000,"VideoId":200001,"ActionType":1,"CreateDate":"2023-02-14T05:20:20+08:00"}`,
		},
		{
			name: "get favorite video id list by userid",
			op: func() (string, error) {
				list, err := dao.GetFavoriteVideoIdListByUserId(favList[0].UserId)
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
			name: "get favorite user id list by video id",
			op: func() (string, error) {
				list, err := dao.GetFavoriteUserList(favList[0].VideoId)
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
			wantResult: `[100002,100000]`,
		},
		{
			name: "get favorite count by video id",
			op: func() (string, error) {
				count, err := dao.GetFavoriteCountByVideoId(favList[0].VideoId)
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
				return "", dao.UpdateFavorite(favList[0].UserId, favList[0].VideoId, consts.UnLike)
			},
			wantErr: false,
		},
		{
			name: "get favorite user id list by video id after cancel favorite",
			op: func() (string, error) {
				list, err := dao.GetFavoriteUserList(favList[0].VideoId)
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
			wantResult: `[100002]`,
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
