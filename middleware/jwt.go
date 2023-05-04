package middleware

import (
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/jwt"
	"github.com/winjeg/go-commons/log"
	"github.com/winjeg/irisword/ret"
)

// JWTConfig Json Web Token config
type JWTConfig struct {
	Name         string `json:"name" yaml:"name"`
	Secret       string `json:"key" yaml:"key"`
	Expire       int    `json:"expire" yaml:"expire"`
	Claims       func() interface{}
	Deserializer func(data []byte) (interface{}, error)
}

var localJWTCfg *JWTConfig = nil

func NewJWT(cfg *JWTConfig) iris.Handler {
	localJWTCfg = cfg
	return func(ctx iris.Context) {
		session := GetFromJWT(ctx)
		if session != nil {
			ctx.Next()
			return
		}
		sigKey := []byte(cfg.Secret)
		claims := cfg.Claims()
		token, err := jwt.Sign(jwt.HS256, sigKey, claims, jwt.MaxAge(15*time.Minute))
		if err != nil {
			ret.ServerError(ctx, "error generating token")
			return
		}
		loc, _ := time.LoadLocation("Asia/Shanghai")
		ctx.SetCookie(&http.Cookie{
			Name:       cfg.Name,
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
		setJWTSession(ctx)
	}
}

func setJWTSession(ctx iris.Context) {
	if localJWTCfg == nil {
		log.GetLogger(nil).Warnf("JWTSession middleware not initialized!")
		ret.BadRequest(ctx, "JWT not initialized!")
		return
	}
	sigKey := []byte(localJWTCfg.Secret)
	tk := ctx.GetCookie(localJWTCfg.Name)
	verifiedToken, err := jwt.Verify(jwt.HS256, sigKey, []byte(tk))
	if err != nil || verifiedToken.StandardClaims.IssuedAt > time.Now().Unix() {
		ret.BadRequest(ctx, "unauthorized!")
		return
	}
	session, err := localJWTCfg.Deserializer(verifiedToken.Payload)
	if session != nil && err == nil {
		ctx.Values().Set(localJWTCfg.Name, session)
	}
	ctx.Next()
}

// GetFromJWT get raw claims info  set from config
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
