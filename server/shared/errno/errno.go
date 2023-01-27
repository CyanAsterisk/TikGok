package errno

import (
	"fmt"

	"github.com/CyanAsterisk/TikGok/server/shared/kitex_gen/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

// NewErrNo return ErrNo
func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success  = NewErrNo(int64(errno.Err_Success), "Success")
	ParamsEr = NewErrNo(int64(errno.Err_ParamsErr), "Params err")

	RPCInteractionErr    = NewErrNo(int64(errno.Err_InteractionServerErr), "rpc call interaction server error")
	InteractionServerErr = NewErrNo(int64(errno.Err_ParamsErr), "interaction server error")

	RPCSocialityErr    = NewErrNo(int64(errno.Err_ParamsErr), "rpc call sociality server error")
	SocialityServerErr = NewErrNo(int64(errno.Err_ParamsErr), "sociality server error")

	RPCUserErr          = NewErrNo(int64(errno.Err_ParamsErr), "rpc call user server error")
	UserServerErr       = NewErrNo(int64(errno.Err_ParamsErr), "user server error")
	UserAlreadyExistErr = NewErrNo(int64(errno.Err_ParamsErr), "user already exist")
	UserNotFoundErr     = NewErrNo(int64(errno.Err_ParamsErr), "user not found")
	AuthorizeFailErr    = NewErrNo(int64(errno.Err_ParamsErr), "authorize failed")

	RPCVideoErr    = NewErrNo(int64(errno.Err_ParamsErr), "rpc call video server error")
	VideoServerErr = NewErrNo(int64(errno.Err_ParamsErr), "video server error")
)

// SendResponse pack response
func SendResponse(c *app.RequestContext, err ErrNo, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    err.ErrCode,
		Message: err.ErrMsg,
		Data:    data,
	})
}
