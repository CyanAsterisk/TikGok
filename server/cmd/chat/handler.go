package main

import (
	"context"
	chat "github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct{}

// ChatHistory implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) ChatHistory(ctx context.Context, req *chat.DouyinMessageChatRequest) (resp *chat.DouyinMessageChatResponse, err error) {
	// TODO: Your code here...
	return
}

// SentMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SentMessage(ctx context.Context, req *chat.DouyinMessageActionRequest) (resp *chat.DouyinMessageActionResponse, err error) {
	// TODO: Your code here...
	return
}

// LatestMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) LatestMessage(ctx context.Context, req *chat.DouyinMessageLatestRequest) (resp *chat.DouyinMessageLatestResponse, err error) {
	// TODO: Your code here...
	return
}
