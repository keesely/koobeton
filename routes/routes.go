// routes.go kee > 2019/12/09

package routes

import (
	//"github.com/kataras/iris"
	//"github.com/kataras/iris/mvc"
	core "koobeton/app"
	c "koobeton/controllers"
)

type router struct {
	Path   string
	Action interface{}
}

var (
	routers = []core.Router{
		core.Router{"/", new(c.IndexController)},
		core.Router{"/user", new(c.UserController)},
		core.Router{"/simple", new(c.SimpleController)},
		core.Router{"/test", func(app core.IrisMvcRoute) {
			app.Mvc.Handle(new(c.SimpleController))
		}},
		core.Router{"get:/test2/{id:int}", func(app core.Context) {
			app.Writef("test %s", app.Params().Get("id"))
		}},
	}

	errors = map[int]core.OnError{
		404: core.OnError(func(ctx core.Context) {
			ctx.JSON(core.Json{
				"code": 404,
				"msg":  "not found",
			})
		}),
		500: core.OnError(func(ctx core.Context) {
			ctx.JSON(c.Result{500, "server error", nil, nil})
		}),
	}
)

func init() {
	core.Routers(routers)
	core.OnErrorCodes(errors)
}
