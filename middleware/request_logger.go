//go:build go1.16

package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/winjeg/go-commons/log"
	"time"
)

func DefaultLogConfig() *RequestLogConfig {
	return &RequestLogConfig{
		WithLabels: true,
		LogLong:    true,
		Threshold:  1000000,
	}
}

type RequestLogConfig struct {
	WithLabels bool  `json:"withLabels" yaml:"with-labels"`
	LogLong    bool  `json:"logLong" yaml:"log-long"`
	Threshold  int64 `json:"threshold" yaml:"threshold"`
}

func NewRequestLogger(cfg *RequestLogConfig) iris.Handler {
	if cfg == nil {
		cfg = DefaultLogConfig()
	}
	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		cost := time.Now().Sub(start).Microseconds()
		errMsg := ""
		if ctx.Err() != nil {
			errMsg = ctx.Err().Error()
		}
		if cfg.WithLabels {
			log.GetLogger(nil).Infof("code: %d, cost: %dµs, method: %s, path: %s, host: %s %s\n",
				ctx.GetStatusCode(), cost, ctx.Method(), ctx.Path(), ctx.RemoteAddr(), errMsg)
			if cfg.LogLong && cost >= cfg.Threshold {
				log.GetLogger(nil).Warnf("high cost: code: %d, cost: %dµs, method: %s, path: %s, host: %s %s\n",
					ctx.GetStatusCode(), cost, ctx.Method(), ctx.Path(), ctx.RemoteAddr(), errMsg)
			}
		} else {
			log.GetLogger(nil).Infof("%d, %dµs, %s, %s, %s %s\n",
				ctx.GetStatusCode(), cost, ctx.Method(), ctx.Path(), ctx.RemoteAddr(), errMsg)
			if cfg.LogLong && cost >= cfg.Threshold {
				log.GetLogger(nil).Warnf("high cost: %d, %dµs, %s, %s, %s %s\n",
					ctx.GetStatusCode(), cost, ctx.Method(), ctx.Path(), ctx.RemoteAddr(), errMsg)
			}
		}
	}
}
