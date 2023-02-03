package pkg

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
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
		CreateTime: m.CreateDate.Format("mm-dd"),
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
