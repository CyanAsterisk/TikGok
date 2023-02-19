package pkg

import (
	"context"
	"fmt"
	"time"

	"testing"

	"github.com/CyanAsterisk/TikGok/server/cmd/interaction/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestCommentLifecycle(t *testing.T) {
	c := context.Background()
	cleanFunc, client, err := test.RunWithRedisInDocker(consts.RedisCommentClientDB, t)
	defer cleanFunc()
	if err != nil {
		t.Fatal(err)
	}
	manager := NewCommentRedisManager(client)

	commentList := make([]*model.Comment, 0)
	timeStamp := int64(1676323214)
	for i := int64(0); i < 8; i++ {

		commentId := i + 1*100000
		uid := i%4 + 2*100000
		videoId := i%2 + 300000
		c := &model.Comment{
			ID:          commentId,
			UserId:      uid,
			VideoId:     videoId,
			ActionType:  consts.ValidComment,
			CommentText: fmt.Sprintf("user%d comment on video%d on %d", uid, videoId, timeStamp+i),
			CreateDate:  timeStamp + i,
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
				time.Sleep(time.Second)
				for _, comment := range commentList {
					if err = manager.CreateComment(c, comment); err != nil {
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
				err = manager.CreateComment(c, commentList[0])
				return "", err
			},
			wantErr: true,
		},
		{
			name: "get comment list by video id",
			op: func() (string, error) {
				list, err := manager.GetCommentListByVideoId(c, commentList[0].VideoId)
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
			wantResult: `[{"ID":100006,"UserId":200002,"VideoId":300000,"ActionType":1,"CommentText":"user200002 comment on video300000 on 1676323220","CreateDate":1676323220},{"ID":100004,"UserId":200000,"VideoId":300000,"ActionType":1,"CommentText":"user200000 comment on video300000 on 1676323218","CreateDate":1676323218},{"ID":100002,"UserId":200002,"VideoId":300000,"ActionType":1,"CommentText":"user200002 comment on video300000 on 1676323216","CreateDate":1676323216},{"ID":100000,"UserId":200000,"VideoId":300000,"ActionType":1,"CommentText":"user200000 comment on video300000 on 1676323214","CreateDate":1676323214}]`,
		},
		{
			name: "delete comment by id",
			op: func() (string, error) {
				err = manager.DeleteComment(c, commentList[0].ID)
				return "", err
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate delete comment by id",
			op: func() (string, error) {
				err = manager.DeleteComment(c, commentList[0].ID)
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
