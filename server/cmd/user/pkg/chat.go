package pkg

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/chat/chatservice"
)

type ChatManager struct {
	ChatService chatservice.Client
}

func (m *ChatManager) BatchGetLatestMessage(ctx context.Context, userId int64, toUserIdList []int64) ([]*chat.LatestMsg, error) {
	resp, err := m.ChatService.BatchGetLatestMessage(ctx, &chat.DouyinMessageBatchGetLatestRequest{
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
