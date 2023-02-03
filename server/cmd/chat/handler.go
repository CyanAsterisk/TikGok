package main

import (
	"context"
	"time"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/consts"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct{}

// ChatHistory implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) ChatHistory(_ context.Context, req *chat.DouyinMessageChatRequest) (resp *chat.DouyinMessageChatResponse, err error) {
	resp = new(chat.DouyinMessageChatResponse)
	msgs, err := dao.GetMessages(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("get chat history error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get chat history error"))
		return resp, nil
	}
	resp.MessageList = pkg.Messages(msgs)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// SentMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SentMessage(_ context.Context, req *chat.DouyinMessageActionRequest) (resp *chat.DouyinMessageActionResponse, err error) {
	resp = new(chat.DouyinMessageActionResponse)
	err = dao.ChatAction(&model.Message{
		ToUserId:   req.ToUserId,
		FromUserId: req.UserId,
		Content:    req.Content,
		CreateDate: time.Now(),
	})
	if err != nil {
		klog.Error("sent message error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("sent message error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// LatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) LatestMessage(_ context.Context, req *chat.DouyinMessageLatestRequest) (resp *chat.DouyinMessageLatestResponse, err error) {
	resp = new(chat.DouyinMessageLatestResponse)
	msg, err := dao.GetLatestMessage(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("get latest message error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get latest message error"))
		return resp, nil
	}
	if msg.FromUserId == req.UserId {
		resp.MsgType = consts.SentMessage
	} else {
		resp.MsgType = consts.ReceiveMessage
	}
	resp.Message = msg.Content
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
