package pkg

import (
	"context"
	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction/interactionserver"
)

type InteractionManager struct {
	client interactionserver.Client
}

func NewInteractionManager(client interactionserver.Client) *InteractionManager {
	return &InteractionManager{client: client}
}

func (i *InteractionManager) GetInteractInfo(ctx context.Context, userId int64) (*base.UserInteractInfo, error) {
	resp, err := i.client.GetUserInteractInfo(ctx, &interaction.DouyinGetUserInteractInfoRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.InteractInfo, nil
}

func (i *InteractionManager) BatchGetInteractInfo(ctx context.Context, userIdList []int64) ([]*base.UserInteractInfo, error) {
	resp, err := i.client.BatchGetUserInteractInfo(ctx, &interaction.DouyinBatchGetUserInteractInfoRequest{
		UserIdList: userIdList,
	})
	if err != nil {
		return nil, err
	}
	if resp.BaseResp.StatusCode != int32(errno.Success.ErrCode) {
		return nil, errno.InteractionServerErr.WithMessage(resp.BaseResp.StatusMsg)
	}
	return resp.InteractInfoList, nil
}
