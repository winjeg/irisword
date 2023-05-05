package main

import (
	"github.com/kataras/iris/v12"
	"github.com/winjeg/irisword/middleware"
)

func regCORS(app *iris.Application) {
	app.UseRouter(middleware.NewCORS(&middleware.CorsConfig{AllowOrigin: []string{"*:8000"}}))
}
