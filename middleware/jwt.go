package middleware

import "github.com/kataras/iris/v12"

// JWTConfig Json Web Token config
type JWTConfig struct {
	Key string `json:"key" yaml:"key"`
}

func NewJWT(config *JWTConfig) iris.Handler {
	return nil
}
