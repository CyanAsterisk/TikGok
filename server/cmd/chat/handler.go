package main

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/dao"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/pkg"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/CyanAsterisk/TikGok/server/shared/tools"
	"github.com/cloudwego/kitex/pkg/klog"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct {
	RedisManager
	Publisher
	Subscriber
}

// RedisManager defines the redis interface.
type RedisManager interface {
	Action(context.Context, *chat.DouyinMessageActionRequest) error
	GetMessages(uid int64, toUid int64) ([]*model.Message, error)
	GetLatestMessage(uid int64, toUid int64) (*model.Message, error)
	BatchGetLatestMessage(uid int64, toUid []int64) ([]*model.Message, error)
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
	msgs, err := s.RedisManager.GetMessages(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("get chat history by redis error", err)
		msgs, err = dao.GetMessages(req.UserId, req.ToUserId)
		if err != nil {
			klog.Error("get chat history by mysql error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get chat history error"))
			return resp, nil
		}
		resp.MessageList = pkg.Messages(msgs)
		resp.BaseResp = tools.BuildBaseResp(nil)
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
	err = s.RedisManager.Action(ctx, req)
	if err != nil {
		klog.Error("sent message by redis error", err)
		resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("sent message error"))
		return resp, nil
	}
	//err = dao.ChatAction(&model.Message{
	//	ToUserId:   req.ToUserId,
	//	FromUserId: req.UserId,
	//	Content:    req.Content,
	//	CreateDate: time.Now(),
	//})
	//if err != nil {
	//	klog.Error("sent message error", err)
	//	resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("sent message error"))
	//	return resp, nil
	//}
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}

// BatchGetLatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) BatchGetLatestMessage(ctx context.Context, req *chat.DouyinMessageBatchGetLatestRequest) (resp *chat.DouyinMessageBatchGetLatestResponse, err error) {
	resp = new(chat.DouyinMessageBatchGetLatestResponse)
	msgList, err := s.RedisManager.BatchGetLatestMessage(req.UserId, req.ToUserIdList)
	if err != nil {
		klog.Error("batch get latest message by redis error", err)
		msgList, err = dao.BatchGetLatestMessage(req.UserId, req.ToUserIdList)
		if err != nil {
			klog.Error("batch get latest message by mysql error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get latest message error"))
			return resp, nil
		}
	}
	resp.LatestMsgList = pkg.LatestMsgs(msgList, req.UserId)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return
}

// GetLatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) GetLatestMessage(_ context.Context, req *chat.DouyinMessageGetLatestRequest) (resp *chat.DouyinMessageGetLatestResponse, err error) {
	resp = new(chat.DouyinMessageGetLatestResponse)
	msg, err := s.RedisManager.GetLatestMessage(req.UserId, req.ToUserId)
	if err != nil {
		klog.Error("get latest message by redis error", err)
		msg, err = dao.GetLatestMessage(req.UserId, req.ToUserId)
		if err != nil {
			klog.Error("get latest message by mysql error", err)
			resp.BaseResp = tools.BuildBaseResp(errno.ChatServerErr.WithMessage("get latest message error"))
			return resp, nil
		}
	}
	resp.LatestMsg = pkg.LatestMsg(msg, req.UserId)
	resp.BaseResp = tools.BuildBaseResp(nil)
	return resp, nil
}
