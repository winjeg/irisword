package main

import "github.com/kataras/iris/v12"

// this package is mainly used as test case purpose.
func main() {
	app := iris.New()

	app.Use()

	app.Get("/ping", func(ctx iris.Context) {
		ctx.Text("pong")
	})

	app.Listen(":8000")
}
