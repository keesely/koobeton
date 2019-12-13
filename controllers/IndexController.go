// IndexController.go kee > 2019/12/09

package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/keesely/kiris"
	"koobeton/app"
	"koobeton/models"
)

type IndexController struct{}

func (C *IndexController) Get(c iris.Context) {
	name := app.Config.Get("app.name", "Hell")
	author := &models.Author{Name: "kee_" + fmt.Sprintf("%d", kiris.Rand(100, 110))}
	app.NewOrm().Insert(author)

	author = &models.Author{Name: "kee_" + fmt.Sprintf("%d", kiris.Rand(100, 110))}
	if ok, _ := app.NewOrm().Get(author); ok {
		c.JSON(author)
	} else {
		c.Writef("Welcome to koobeton -> %s", name)
	}
}
