package tools

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/base"
)

func BuildBaseResp(err error) *base.DouyinBaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *base.DouyinBaseResponse {
	return &base.DouyinBaseResponse{
		StatusCode: int32(err.ErrCode),
		StatusMsg:  err.ErrMsg,
	}
}
