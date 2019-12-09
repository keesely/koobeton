// main.go kee > 2019/12/09

package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Run(
		iris.Addr(":8000"),
		iris.WithCharset("UTF-8"),
	)
}
