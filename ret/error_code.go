package ret

type ErrorCode struct {
	Code int
	Msg  string
}

var (
	InternalError = ErrorCode{500, "服务器内部错误"}
	IllegalParam  = ErrorCode{400, "参数非法"}
	BizErr        = ErrorCode{1000, "业务错误"}
	Success       = ErrorCode{0, "success"}
)
