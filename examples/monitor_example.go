package main

import (
	"github.com/kataras/iris/v12"
	"github.com/winjeg/irisword/middleware"
)

func regMonitor(app *iris.Application) {
	handler := middleware.NewIrisMonitor(&middleware.MonitorConfig{
		Port: 10000,
		Path: "/metrics",
		Tags: nil,
	})

	app.UseRouter(handler)
}
