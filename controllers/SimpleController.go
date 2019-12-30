// SimpleController.go kee > 2019/12/10

package controllers

import (
	"koobeton/app"
	//"koobeton/models"
)

type SimpleController struct{}

func (c *SimpleController) Get(ctx app.Context) (int, string) {
	//user := &models.User{Email: "chinboy2012@gmail.com"}
	//passwd, _ := user.GetPasswd()
	//app.Logger().Println("> Passwd: ", passwd)
	//user.Get(true)
	//var user = app.Json{}
	//app.NewOrm().Table("user").Where("Email=?", "chinboy2012@gmail.com").Get(&user)
	//passwd := user.Get("passwd")
	//app.Logger().Println(passwd)
	//cUser := ctx.Params().Get("cUser")
	//wspace := &models.Workspace{Name: "hello world", URL: "https://location"}
	sss := "hello world"
	return 200, sss
	//return 200, cUser
}
