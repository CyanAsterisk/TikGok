package pack

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/interaction"
)

func BuildBaseResp(err error) *interaction.DouyinBaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.RequestServerFail.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *interaction.DouyinBaseResponse {
	return &interaction.DouyinBaseResponse{
		StatusCode: int32(err.ErrCode),
		StatusMsg:  err.ErrMsg,
	}
}
