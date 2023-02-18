package dao

import (
	"fmt"
	"testing"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/shared/test"
	"github.com/bytedance/sonic"
)

func TestChatLifeCycle(t *testing.T) {
	cleanUpFunc, db, err := test.RunWithMySQLInDocker(t)
	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	dao := NewMessage(db)
	var msgList []*model.Message
	timeStamp := int64(1676323214)
	for i := int64(0); i < 10; i++ {
		uid1 := 1 - i%2 + 100000
		uid2 := 200001 - uid1
		msg := &model.Message{
			ID:         200000 + i,
			ToUserId:   uid2,
			FromUserId: uid1,
			Content:    fmt.Sprintf("User %d send message%d to %d", uid1, i, uid2),
			CreateTime: timeStamp + i,
		}
		msgList = append(msgList, msg)
	}

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "set message",
			op: func() (string, error) {
				for _, msg := range msgList {
					err = dao.ChatAction(msg)
					if err != nil {
						return "", err
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "get messages",
			op: func() (string, error) {
				msg, err := dao.GetMessages(msgList[0].ToUserId, msgList[0].FromUserId, timeStamp+2)
				if err != nil {
					return "", nil
				}
				res, err := sonic.Marshal(msg)
				return string(res), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":200001,"ToUserId":100001,"FromUserId":100000,"Content":"User 100000 send message1 to 100001","CreateTime":1676323215},{"ID":200000,"ToUserId":100000,"FromUserId":100001,"Content":"User 100001 send message0 to 100000","CreateTime":1676323214}]`,
		},
		{
			name: "get latest Message",
			op: func() (string, error) {
				m, err := dao.GetLatestMessage(msgList[0].FromUserId, msgList[0].ToUserId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(m)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `{"ID":200009,"ToUserId":100001,"FromUserId":100000,"Content":"User 100000 send message9 to 100001","CreateTime":1676323223}`,
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
