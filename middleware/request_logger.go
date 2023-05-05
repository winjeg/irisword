//go:build go1.16

package middleware

import "github.com/kataras/iris/v12"

type RequestLogConfig struct {
}

func NewRequestLogger(cfg *RequestLogConfig) iris.Handler {
	return nil
}
