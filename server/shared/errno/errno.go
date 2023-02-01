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
	Success    = NewErrNo(int64(errno.Err_Success), "success")
	ParamsEr   = NewErrNo(int64(errno.Err_ParamsErr), "params err")
	ServiceErr = NewErrNo(int64(errno.Err_ServiceErr), "service err")

	RPCInteractionErr    = NewErrNo(int64(errno.Err_RPCInteractionErr), "rpc call interaction server error")
	InteractionServerErr = NewErrNo(int64(errno.Err_InteractionServerErr), "interaction server error")

	RPCSocialityErr    = NewErrNo(int64(errno.Err_RPCSocialityErr), "rpc call sociality server error")
	SocialityServerErr = NewErrNo(int64(errno.Err_SocialityServerErr), "sociality server error")

	RPCUserErr          = NewErrNo(int64(errno.Err_RPCUserErr), "rpc call user server error")
	UserServerErr       = NewErrNo(int64(errno.Err_UserServerErr), "user server error")
	UserAlreadyExistErr = NewErrNo(int64(errno.Err_UserAlreadyExistErr), "user already exist")
	UserNotFoundErr     = NewErrNo(int64(errno.Err_UserNotFoundErr), "user not found")
	AuthorizeFailErr    = NewErrNo(int64(errno.Err_AuthorizeFailErr), "authorize failed")

	RPCVideoErr    = NewErrNo(int64(errno.Err_RPCVideoErr), "rpc call video server error")
	VideoServerErr = NewErrNo(int64(errno.Err_VideoServerErr), "video server error")

	RPCChatErr    = NewErrNo(int64(errno.Err_RPCChatErr), "rpc call chat server error")
	ChatServerErr = NewErrNo(int64(errno.Err_ChatServerErr), "chat server error")
)

// SendResponse pack response
func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}
