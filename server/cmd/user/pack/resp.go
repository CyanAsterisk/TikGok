package pack

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/user"
)

func BuildBaseResp(err error) *user.DouyinBaseResponse {
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

func baseResp(err errno.ErrNo) *user.DouyinBaseResponse {
	return &user.DouyinBaseResponse{
		StatusCode: int32(err.ErrCode),
		StatusMsg:  err.ErrMsg,
	}
}
