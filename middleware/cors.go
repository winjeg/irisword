package middleware

import "github.com/kataras/iris/v12"

type CorsConfig struct {
	Domain string `json:"domain" yaml:"domain"`
}

func NewCORS(cfg *CorsConfig) iris.Handler {
	return nil
}
