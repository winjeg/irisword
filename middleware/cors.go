//go:build go1.16

// Package middleware recommended use:  app.UseRouter
// cors configs
package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"net/http"
	"strings"
)

const (
	originRequestHeader    = "Origin"
	allowOriginHeader      = "Access-Control-Allow-Origin"
	allowCredentialsHeader = "Access-Control-Allow-Credentials"
	referrerPolicyHeader   = "Referrer-Policy"
	exposeHeadersHeader    = "Access-Control-Expose-Headers"
	requestMethodHeader    = "Access-Control-Request-Method"
	requestHeadersHeader   = "Access-Control-Request-Headers"
	allowMethodsHeader     = "Access-Control-Allow-Methods"
	allowAllMethodsValue   = "*"
	allowHeadersHeader     = "Access-Control-Allow-Headers"
	maxAgeHeader           = "Access-Control-Max-Age"
	varyHeader             = "Vary"
)

type CorsConfig struct {
	AllowOrigin []string `json:"allowOrigin" yaml:"allow-origin"`
	AllowHeader string   `json:"allowHeader" yaml:"allow-header"`
	MaxAge      string   `json:"maxAge" yaml:"max-age"`
}

// NewCORS POST、PUT、PATCH和DELETE 标准上要求浏览器在这些请求上都要加上 Origin Header
func NewCORS(cfg *CorsConfig) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Header(varyHeader, originRequestHeader)
		if ctx.Method() == http.MethodOptions {
			ctx.Header(varyHeader, requestMethodHeader)
			ctx.Header(varyHeader, requestHeadersHeader)
		}
		requestOrigin := ctx.GetHeader(originRequestHeader)
		// 一般都是 *.xx.com 或者 zz.xx.com
		for _, origin := range cfg.AllowOrigin {
			if strings.Index(origin, "*") == 0 {
				// 默认情况下 origin 获取不到
				if strings.Index(requestOrigin, origin[1:]) == -1 {
					ctx.StopWithStatus(http.StatusForbidden)
					return
				}
			} else {
				if !strings.EqualFold(requestOrigin, origin) {
					ctx.StopWithStatus(http.StatusForbidden)
					return
				}
			}
		}

		if len(cfg.AllowOrigin) == 0 { // if we allow empty origins, set it to wildcard.
			cfg.AllowOrigin = []string{"*"}
		}

		if len(cfg.MaxAge) == 0 {
			cfg.MaxAge = "86400"
		}
		if len(cfg.AllowHeader) == 0 {
			cfg.AllowHeader = "*"
		}

		ctx.Header(allowOriginHeader, strings.Join(cfg.AllowOrigin, ","))
		ctx.Header(allowCredentialsHeader, "true")
		ctx.Header(referrerPolicyHeader, cors.NoReferrerWhenDowngrade.String())
		ctx.Header(exposeHeadersHeader, "*, Authorization, X-Authorization")
		if ctx.Method() == http.MethodOptions {
			ctx.Header(allowMethodsHeader, allowAllMethodsValue)
			ctx.Header(allowHeadersHeader, cfg.AllowHeader)
			ctx.Header(maxAgeHeader, cfg.MaxAge)
			ctx.StatusCode(http.StatusNoContent)
			return
		}
	}
}
