// router.go kee > 2019/12/10

package app

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"strings"
)

type Router struct {
	Path   string
	Action interface{}
}

type OnError func(Context)

type IrisMvcRoute struct {
	*mvc.Application
}

type Context iris.Context

type Party struct {
	iris.Party
}

func init() {
	NewApp()
}

func Routers(routers []Router) {
	//for _, c := range routers {
	for i := 0; i < len(routers); i++ {
		c := routers[i]
		var (
			methods = []string{
				"GET", "POST", "PUT", "DELETE", "OPTIONS", "TRACE", "HEAD", "CONNECT", "PATCH", "ANY",
			}
			method = "GET"
			path   = c.Path
			typeof = fmt.Sprintf("%T", c.Action)
		)

		if idx := strings.Index(c.Path, ":"); idx != -1 {
			m := strings.ToUpper(c.Path[:idx])
			for _, v := range methods {
				if m == v {
					method = m
					path = c.Path[idx+1:]
				}
			}
		}

		if typeof == "func(app.IrisMvcRoute)" {
			// iris mvc route set
			mvc.Configure(app.Party(path), func(app *mvc.Application) {
				c.Action.(func(IrisMvcRoute))(IrisMvcRoute{app})
			})
		} else if typeof == "func(*mvc.Application)" {
			// iris mvc route func
			mvc.Configure(app.Party(path), c.Action.(func(*mvc.Application)))
		} else if typeof == "func(app.Context)" {
			// iris application route
			app.Handle(method, path, func(ctx iris.Context) {
				c.Action.(func(Context))(ctx)
			})
		} else {
			// new iris mvc route
			mvc.New(app.Party(c.Path)).Handle(c.Action)
		}
	}
}

func PartyFunc(route func(Party)) func(iris.Party) {
	return func(r iris.Party) {
		party := Party{r}
		route(party)
	}
}

func RouteHandle(handle func(Context)) func(iris.Context) {
	return func(ctx iris.Context) {
		handle(ctx.(Context))
	}
}

func OnErrorCodes(errors map[int]OnError) {
	for errCode, callback := range errors {
		if 0 <= errCode && 600 >= errCode {
			act := func(ctx iris.Context) {
				callback(ctx)
			}
			switch errCode {
			case 0:
				app.OnAnyErrorCode(act)
				break
			default:
				app.OnErrorCode(errCode, act)
				break
			}
		}
	}
}

func GetClientIP(ctx iris.Context) string {
	ip := ctx.Request().Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request().Header.Get("X-Real-Ip")
	}

	if ip == "" {
		if ip = ctx.Request().Header.Get("X-Appengine-Remote-Addr"); ip == "" {
			return ctx.RemoteAddr()
		}
	}
	return ip
}
