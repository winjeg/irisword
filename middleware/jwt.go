//go:build go1.16

// Package middleware call NewJWT first
// basicly one app need all these three functions exposed

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/jwt"
	"github.com/winjeg/go-commons/log"
	"github.com/winjeg/irisword/ret"
)

// JWTConfig Json Web Token config
type JWTConfig struct {
	Name   string `json:"name" yaml:"name"`     // Name of the JWT
	Expire int    `json:"expire" yaml:"expire"` // JWT expire time
	Domain string `json:"domain" yaml:"domain"` // JWT domain
	Secret string `json:"key" yaml:"key"`       // Secret of the JWT

	Claims       func(ctx iris.Context) (interface{}, error) // set if the return value from NewJWT is used as middleware
	Deserializer func(data []byte) (interface{}, error)      // must be set
}

var (
	localJWTCfg *JWTConfig = nil
	lock                   = sync.Mutex{}
)

// NewJWT return the JWT middleware for iris web framework
// This returned middleware can be used in the project
// only if you need to set Token only here, also you can use in any project if you want,
// just make sure the configuration is the same.
func NewJWT(cfg *JWTConfig) iris.Handler {
	if localJWTCfg != nil {
		return func(ctx iris.Context) {
			ctx.Next()
		}
	}
	lock.Lock()
	localJWTCfg = cfg
	lock.Unlock()

	return func(ctx iris.Context) {
		session := GetFromJWT(ctx)
		if session != nil {
			ret.Ok(ctx, session)
			return
		}
		sigKey := []byte(cfg.Secret)
		claims, err := cfg.Claims(ctx)
		if err != nil {
			ret.Unauthorized(ctx, err.Error())
			return
		}
		token, err := jwt.Sign(jwt.HS256, sigKey, claims, jwt.MaxAge(time.Duration(cfg.Expire)*time.Second))
		if err != nil {
			ret.ServerError(ctx, "error generating token")
			return
		}
		loc, _ := time.LoadLocation("Asia/Shanghai")
		ctx.SetCookie(&http.Cookie{
			Name:       cfg.Name,
			Domain:     cfg.Domain,
			Value:      string(token),
			Path:       "/",
			Expires:    time.Now().Add(time.Second * time.Duration(cfg.Expire)).In(loc).Local(),
			RawExpires: "",
			Secure:     false,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
		// when all set, set the session for later use.
		ret.Ok(ctx, claims)
	}
}

func JWTSession(ctx iris.Context) {
	if localJWTCfg == nil {
		ret.Unauthorized(ctx, "user not login")
		return
	}
	if GetFromJWT(ctx) != nil {
		ctx.Next()
		return
	}
	sigKey := []byte(localJWTCfg.Secret)
	tk := ctx.GetCookie(localJWTCfg.Name)
	verifiedToken, err := jwt.Verify(jwt.HS256, sigKey, []byte(tk))
	if err != nil || verifiedToken.StandardClaims.IssuedAt > time.Now().Unix() {
		ret.Unauthorized(ctx, "unauthorized!")
		return
	}
	session, err := localJWTCfg.Deserializer(verifiedToken.Payload)
	if session != nil && err == nil {
		ctx.Values().Set(localJWTCfg.Name, session)
	}
	ctx.Next()
}

// GetFromJWT get raw claims info  set from config
// When called this method will need JWTConfig initialed first, call NewJWT if necessary
// You can get token from the
func GetFromJWT(ctx iris.Context) interface{} {
	if localJWTCfg == nil {
		log.GetLogger(nil).Warnf("GetFromJWT middleware not initialized!")
		return nil
	}
	if session := ctx.Values().Get(localJWTCfg.Name); session != nil {
		return session
	}
	tk := ctx.GetCookie(localJWTCfg.Name)
	sigKey := []byte(localJWTCfg.Secret)
	verifiedToken, err := jwt.Verify(jwt.HS256, sigKey, []byte(tk))
	if err != nil {
		return nil
	}
	session, err := localJWTCfg.Deserializer(verifiedToken.Payload)
	if session != nil && err == nil {
		ctx.Values().Set(localJWTCfg.Name, session)
		return session
	}
	return nil
}
