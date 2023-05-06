//go:build go1.16

package main

import (
	"github.com/kataras/iris/v12"
	rec "github.com/kataras/iris/v12/middleware/recover"
	"github.com/winjeg/irisword/middleware"
)

// this package is mainly used as test case purpose.
func main() {
	app := iris.New()
	app.UseRouter(middleware.NewRequestLogger(nil))
	regMonitor(app)
	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.Text("pong")
	})

	app.UseRouter(rec.New())
	regCORS(app)
	regJWT(app)

	_ = app.Listen(":8000")
}
