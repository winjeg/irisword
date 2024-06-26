//go:build go1.16

package ret

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/winjeg/go-commons/log"
)

var logger = log.GetLogger(nil)

// Ok 正常返回
func Ok(ctx iris.Context, data ...interface{}) {
	var d interface{} = nil
	if len(data) > 0 {
		d = data[0]
	}
	err := ctx.JSON(Ret{
		Code: Success.Code,
		Msg:  Success.Msg,
		Data: d,
	})
	logErr(ctx, err)
}

// BadRequest 参数错误
func BadRequest(ctx iris.Context, msg ...string) {
	m := strings.Join(msg, "\n")
	if len(strings.TrimSpace(m)) == 0 {
		m = IllegalParam.Msg
	}
	err := ctx.JSON(Ret{
		Code: IllegalParam.Code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// ServerError 不可遇见的服务器错误， 可以使用此方法
func ServerError(ctx iris.Context, msg ...string) {
	m := strings.Join(msg, "\t")
	if len(strings.TrimSpace(m)) == 0 {
		m = InternalError.Msg
	}
	err := ctx.JSON(Ret{
		Code: InternalError.Code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// BizError 业务错误一般调用此方法
func BizError(ctx iris.Context, code string, msg ...string) {
	m := strings.Join(msg, "\t")
	if len(strings.TrimSpace(m)) == 0 {
		m = BizErr.Msg
	}
	err := ctx.JSON(Ret{
		Code: code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// UnknownError error that could not be identified
func UnknownError(ctx iris.Context, msg ...string) {
	m := strings.Join(msg, "\t")
	if len(strings.TrimSpace(m)) == 0 {
		m = UnknownErr.Msg
	}
	err := ctx.JSON(Ret{
		Code: UnknownErr.Code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// NotFound 找不到资源
func NotFound(ctx iris.Context) {
	err := ctx.JSON(Ret{
		Code: NoFound.Code,
		Msg:  fmt.Sprintf("%s: %s", NoFound.Msg, ctx.Path()),
	})
	logErr(ctx, err)
}

// Unauthorized 未授权
func Unauthorized(ctx iris.Context, msg ...string) {
	m := strings.Join(msg, "\t")
	if len(strings.TrimSpace(m)) == 0 {
		m = NoAuth.Msg
	}
	err := ctx.JSON(Ret{
		Code: NoAuth.Code,
		Msg:  fmt.Sprintf("%s: %s", m, ctx.Path()),
	})
	logErr(ctx, err)
}

func logErr(ctx iris.Context, err error) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("logErr", "error writing response")
	}
}
