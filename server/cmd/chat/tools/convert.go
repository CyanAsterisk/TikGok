package tools

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

// Messages model to idl
func Messages(m []*model.Message) []*base.Message {
	if m == nil {
		return nil
	}
	ml := make([]*base.Message, 0)
	for _, ms := range ml {
		ml = append(ml, ms)
	}
	return ml
}
