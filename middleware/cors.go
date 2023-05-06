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

// NewDefaultCorsCfg 默认cors 配置
func NewDefaultCorsCfg() *CorsConfig {
	return &CorsConfig{
		AllowOrigin: []string{"*"},
		MaxAge:      "7200",
	}
}

// NewCORS POST、PUT、PATCH和DELETE 标准上要求浏览器在这些请求上都要加上 Origin Header
// CORS 安全策略主要应用于浏览器页面跨域访问的时候，对于非浏览器页面请求，业界通常不予以特殊拦截
func NewCORS(cfg *CorsConfig) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Header(varyHeader, originRequestHeader)
		if ctx.Method() == http.MethodOptions {
			ctx.Header(varyHeader, requestMethodHeader)
			ctx.Header(varyHeader, requestHeadersHeader)
		}

		requestOrigin := ctx.GetHeader(originRequestHeader)

		if ctx.Method() == http.MethodPost ||
			ctx.Method() == http.MethodPut ||
			ctx.Method() == http.MethodDelete ||
			ctx.Method() == http.MethodPatch {
			if len(requestOrigin) == 0 {
				ctx.StopWithStatus(http.StatusForbidden)
				return
			}
		}

		// 对于一些访问用途，如CURL, 或者浏览器扩展程序， 传递的origin头则予以放行
		if len(requestOrigin) > 0 && (strings.Index(requestOrigin, "https://") == 0 ||
			strings.Index(requestOrigin, "http://") == 0) {
			for _, origin := range cfg.AllowOrigin {
				// 一般都是 *.xx.com 或者 zz.xx.com
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
		ctx.Next()
	}
}
