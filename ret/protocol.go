package ret

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"github.com/winjeg/go-commons/log"
)

var logger = log.GetLogger(nil)

type Ret struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

// Ok 正常返回
func Ok(ctx iris.Context, data ...interface{}) {
	var d interface{} = nil
	if len(data) > 0 {
		d = data[0]
	}
	_, err := ctx.JSON(Ret{
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
	_, err := ctx.JSON(Ret{
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
	_, err := ctx.JSON(Ret{
		Code: InternalError.Code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// BizError 业务错误一般调用此方法
func BizError(ctx iris.Context, code int, msg ...string) {
	m := strings.Join(msg, "\t")
	if len(strings.TrimSpace(m)) == 0 {
		m = BizErr.Msg
	}
	_, err := ctx.JSON(Ret{
		Code: BizErr.Code,
		Msg:  m,
	})
	logErr(ctx, err)
}

// NotFound 找不到资源
func NotFound(ctx iris.Context) {
	_, err := ctx.JSON(Ret{
		Code: 404,
		Msg:  fmt.Sprintf("route not found: %s", ctx.Path()),
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
