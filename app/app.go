// app.go kee > 2019/12/09

package app

import (
	"flag"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/keesely/kiris"
	"path"
	"sync"
)

var (
	app     *iris.Application
	once    sync.Once
	cnfPath = flag.String("conf-path", "./conf", "application configure path")
	Config  *appConf
)

func NewApp() *iris.Application {
	once.Do(func() {
		flag.Parse()
		// new Application
		app = iris.New()
		// import application configure
		Config = &appConf{kiris.NewYamlLoad(path.Join(*cnfPath, "app.yml"))}
	})

	return app
}

func Run() {
	var (
		irisApp = NewApp()
		listen  = Config.Get("app.listen", ":3000").(string)
		charset = Config.Get("app.charset", "UTF-8").(string)
	)

	logLevel := Config.Get("app.error", "debug").(string)
	irisApp.Logger().SetLevel(logLevel)

	// import iris configure
	c := path.Join(*cnfPath, "iris.yml")
	if kiris.FileExists(c) {
		irisApp.Configure(iris.WithConfiguration(iris.YAML(c)))
	}
	irisApp.Run(iris.Addr(listen), iris.WithCharset(charset))
}

type appConf struct {
	conf *kiris.Yaml
}

func (c *appConf) Get(name string, nfVal interface{}) interface{} {
	if nil == c.conf {
		return nil
	}
	return c.conf.Get(name, nfVal)
}

func Logger() *golog.Logger {
	return app.Logger()
}
