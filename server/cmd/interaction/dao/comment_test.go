package dao

import (
	"fmt"
	"github.com/bytedance/sonic"
	"testing"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
)

func TestCommentLifecycle(t *testing.T) {
	cleanUpFunc, db, err := test.RunWithMySQLInDocker(t)
	defer cleanUpFunc()
	if err != nil {
		t.Fatal(err)
	}
	dao := NewComment(db)

	commentList := make([]*model.Comment, 0)
	timeStamp := int64(1676323214)
	for i := int64(0); i < 8; i++ {
		date := time.Unix(timeStamp+i, 0)
		commentId := i + 1*100000
		uid := i%4 + 2*100000
		videoId := i%2 + 300000
		c := &model.Comment{
			ID:          commentId,
			UserId:      uid,
			VideoId:     videoId,
			ActionType:  consts.ValidComment,
			CommentText: fmt.Sprintf("user%d comment on video%d on %s", uid, videoId, date.Format("2006-01-02 15:04:05")),
			CreateDate:  date,
		}
		commentList = append(commentList, c)
	}

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "create comment",
			op: func() (string, error) {
				for _, c := range commentList {
					if err = dao.CreateComment(c); err != nil {
						return "", err
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate create comment",
			op: func() (string, error) {
				err = dao.CreateComment(commentList[0])
				return "", err
			},
			wantErr: true,
		},
		{
			name: "get comment list by video id",
			op: func() (string, error) {
				list, err := dao.GetCommentListByVideoId(commentList[0].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", nil
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":100006,"UserId":200002,"VideoId":300000,"ActionType":1,"CommentText":"user200002 comment on video300000 on 2023-02-14 05:20:20","CreateDate":"2023-02-14T05:20:20+08:00"},{"ID":100004,"UserId":200000,"VideoId":300000,"ActionType":1,"CommentText":"user200000 comment on video300000 on 2023-02-14 05:20:18","CreateDate":"2023-02-14T05:20:18+08:00"},{"ID":100002,"UserId":200002,"VideoId":300000,"ActionType":1,"CommentText":"user200002 comment on video300000 on 2023-02-14 05:20:16","CreateDate":"2023-02-14T05:20:16+08:00"},{"ID":100000,"UserId":200000,"VideoId":300000,"ActionType":1,"CommentText":"user200000 comment on video300000 on 2023-02-14 05:20:14","CreateDate":"2023-02-14T05:20:14+08:00"}]`,
		},
		{
			name: "get comment id list by video id",
			op: func() (string, error) {
				list, err := dao.GetCommentIdListByVideoId(commentList[0].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", nil
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `["100000","100002","100004","100006"]`,
		},
		{
			name: "delete comment by id",
			op: func() (string, error) {
				err = dao.DeleteComment(commentList[0].ID)
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get comment id list by video id after delete a comment",
			op: func() (string, error) {
				list, err := dao.GetCommentIdListByVideoId(commentList[0].VideoId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(list)
				if err != nil {
					return "", nil
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `["100002","100004","100006"]`,
		},
		{
			name: "duplicate delete comment by id",
			op: func() (string, error) {
				err = dao.DeleteComment(commentList[0].ID)
				return "", err
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
