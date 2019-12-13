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
	Mvc *mvc.Application
}

type Context iris.Context

type Json map[string]interface{}

func init() {
	NewApp()
}

func Routers(routers []Router) {
	for _, c := range routers {
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

		if typeof == "func(utils.IrisMvcRoute)" {
			// iris mvc route set
			mvc.Configure(app.Party(path), func(app *mvc.Application) {
				c.Action.(func(IrisMvcRoute))(IrisMvcRoute{app})
			})
		} else if typeof == "func(*mvc.Application)" {
			// iris mvc route func
			mvc.Configure(app.Party(path), c.Action.(func(*mvc.Application)))
		} else if typeof == "func(utils.Context)" {
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
