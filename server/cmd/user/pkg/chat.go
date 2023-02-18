package pkg

import (
	"context"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
)

type ChatManager struct {
	client chatservice.Client
}

func NewChatManager(client chatservice.Client) *ChatManager {
	return &ChatManager{client: client}
}

func (m *ChatManager) BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*base.LatestMsg, error) {
	resp, err := m.client.BatchGetLatestMessage(ctx, &chat.DouyinMessageBatchGetLatestRequest{
		UserId:       userId,
		ToUserIdList: toUserIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.SocialityServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.LatestMsgList, nil
}
