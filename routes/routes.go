// routes.go kee > 2019/12/09

package routes

import (
	"koobeton/app"
	c "koobeton/controllers"
)

func init() {
	app.Routers([]app.Router{
		//app.Router{"/users", new(c.UserController)},
		// 用户登陆
		app.Router{"/users/login", new(c.LoginController)},

		// 标准API
		app.Router{"/", func(r app.IrisMvcRoute) {
			r.Router.Use(app.RouteHandle(authorization))
			r.Party("/workspace").Handle(new(c.Workspace))
		}},

		app.Router{"get:/test/{id:int}", func(ctx app.Context) {
			ctx.Writef("test 1 %s", ctx.Params().Get("id"))
		}},
	})

	// 错误处理
	app.OnErrorCodes(map[int]app.OnError{
		401: app.OnError(func(ctx app.Context) {
			ctx.JSON(app.Json{
				"code": 401,
				"msg":  "autorization fail",
			})
		}),
		404: app.OnError(func(ctx app.Context) {
			ctx.JSON(app.Json{
				"code": 404,
				"msg":  "not found",
			})
		}),
	})
}
