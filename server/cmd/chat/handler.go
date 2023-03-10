package main

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct {
	Publisher
	Subscriber
	Dao *dao.Message
}

// Publisher defines the publisher interface.
type Publisher interface {
	Publish(context.Context, *chat.DouyinMessageActionRequest) error
}

// Subscriber defines a car update subscriber.
type Subscriber interface {
	Subscribe(context.Context) (ch chan *chat.DouyinMessageActionRequest, cleanUp func(), err error)
}

// GetChatHistory implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) GetChatHistory(_ context.Context, req *chat.DouyinMessageGetChatHistoryRequest) (resp *chat.DouyinMessageGetChatHistoryResponse, err error) {
	resp = new(chat.DouyinMessageGetChatHistoryResponse)
	msgs, err := s.Dao.GetMessages(req.UserId, req.ToUserId, req.PreMsgTime)
	if err != nil {
		klog.Error("get chat history by mysql error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get chat history error"))
		return resp, nil
	}
	resp.MessageList = pkg.Messages(msgs)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// SentMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SentMessage(ctx context.Context, req *chat.DouyinMessageActionRequest) (resp *chat.DouyinMessageActionResponse, err error) {
	resp = new(chat.DouyinMessageActionResponse)
	err = s.Publisher.Publish(ctx, req)
	if err != nil {
		klog.Error("publish message error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("sent message error"))
		return resp, nil
	}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetLatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) BatchGetLatestMessage(_ context.Context, req *chat.DouyinMessageBatchGetLatestRequest) (resp *chat.DouyinMessageBatchGetLatestResponse, err error) {
	resp = new(chat.DouyinMessageBatchGetLatestResponse)
	msgList, err := s.Dao.BatchGetLatestMessage(req.UserId, req.ToUserIdList)
	if err != nil {
		klog.Error("batch get latest message by mysql error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get latest message error"))
		return resp, nil
	}
	resp.LatestMsgList = pkg.LatestMsgs(msgList, req.UserId)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// GetLatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) GetLatestMessage(_ context.Context, req *chat.DouyinMessageGetLatestRequest) (resp *chat.DouyinMessageGetLatestResponse, err error) {
	resp = new(chat.DouyinMessageGetLatestResponse)
	msg, err := s.Dao.GetLatestMessage(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("get latest message by mysql error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get latest message error"))
		return resp, nil
	}
	resp.LatestMsg = pkg.LatestMsg(msg, req.UserId)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
