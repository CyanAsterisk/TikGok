package middleware

import (
	"context"
	"net/http"

	"github.com/CyanAsterisk/TikGok/server/shared/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Recovery() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
			hlog.SystemLogger().CtxErrorf(ctx, "[Recovery] err=%v\nstack=%s", err, stack)
			c.JSON(http.StatusOK, utils.H{
				"status_code": int32(errno.ServiceErr.ErrCode),
				"status_msg":  errno.ServiceErr.ErrMsg,
			})
		},
	))
}
