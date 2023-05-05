//go:build go1.16

package ret

// ErrorCode for most defined errors
type ErrorCode struct {
	Code string
	Msg  string
}

// Ret return value generic data structure
type Ret struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

// BizError should be greater than 1000
// Since we try to use the http code as much as we can
var (
	Success       = ErrorCode{"0", "success"}
	NoAuth        = ErrorCode{"401", "unauthorized!"}
	NoFound       = ErrorCode{"404", "route not found!"}
	InternalError = ErrorCode{"500", "internal error"}
	IllegalParam  = ErrorCode{"400", "illegal param"}
	UnknownErr    = ErrorCode{"999", "unknown error"}
	BizErr        = ErrorCode{"1000", "biz error"}
)
