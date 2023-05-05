package main

import (
	"github.com/kataras/iris/v12"
)

// this package is mainly used as test case purpose.
func main() {
	app := iris.New()
	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.Text("pong")
	})
	regJWT(app)
	_ = app.Listen(":8000")
}
