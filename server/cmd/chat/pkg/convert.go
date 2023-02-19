package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// Message model to idl
func Message(m *model.Message) *base.Message {
	if m == nil {
		return nil
	}
	return &base.Message{
		Id:         m.ID,
		ToUserId:   m.ToUserId,
		FromUserId: m.FromUserId,
		Content:    m.Content,
		CreateTime: m.CreateTime,
	}
}

// Messages model to idl
func Messages(m []*model.Message) []*base.Message {
	if m == nil {
		return nil
	}
	ml := make([]*base.Message, 0)
	for _, ms := range m {
		ml = append(ml, Message(ms))
	}
	return ml
}

// LatestMsg coverts Message model to LatestMsg.
func LatestMsg(m *model.Message, uid int64) *base.LatestMsg {
	if m == nil {
		return nil
	}
	msg := &base.LatestMsg{
		Message: m.Content,
	}
	if m.FromUserId == uid {
		msg.MsgType = consts.SentMessage
	} else {
		msg.MsgType = consts.ReceiveMessage
	}
	return msg
}

// LatestMsgs batch coverts Message list model to LatestMsg list.
func LatestMsgs(ml []*model.Message, uid int64) []*base.LatestMsg {
	if ml == nil {
		return nil
	}
	msgl := make([]*base.LatestMsg, 0)
	for _, m := range ml {
		msgl = append(msgl, LatestMsg(m, uid))
	}
	return msgl
}
