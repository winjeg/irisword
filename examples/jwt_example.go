package main

import (
	"encoding/json"

	"github.com/kataras/iris/v12"
	"github.com/winjeg/irisword/middleware"
	"github.com/winjeg/irisword/ret"
)

func getInfo(ctx iris.Context) (interface{}, error) {
	id, _ := ctx.URLParamInt("id")
	name := ctx.URLParam("name")
	return map[string]interface{}{
		"userId":   id,
		"userName": name,
	}, nil
}

func dec(d []byte) (interface{}, error) {
	m := make(map[string]interface{}, 4)
	err := json.Unmarshal(d, &m)
	return m, err
}

func regJWT(app *iris.Application) {

	// 配置
	jwtCfg := &middleware.JWTConfig{
		Name:         "jwt_token",
		Expire:       1500,
		Domain:       "",
		Secret:       "llklzn1231kz1",
		Claims:       getInfo,
		Deserializer: dec,
	}

	// example, 登录接口设置JWT，
	jwt := middleware.NewJWT(jwtCfg)
	app.Get("/api/login", jwt)

	// 其他接口校验 session
	group := app.Party("/api")
	group.Use(middleware.JWTSession)
	{
		group.Get("/info", func(ctx iris.Context) {
			session := middleware.GetFromJWT(ctx)
			ret.Ok(ctx, session)
			return
		})
	}
}
