// LoginController.go kee > 2019/12/16

package controllers

import (
	"koobeton/app"
	//"koobeton/app/request"
	"koobeton/models"
)

type LoginController struct{}

func (*LoginController) Post(ctx app.Context) Result {
	type LoginForm struct {
		Email  string
		Passwd string
		Salt   string
	}
	var form LoginForm
	ctx.ReadForm(&form)

	user := models.User{Email: form.Email}

	/*
		token := ctx.Request().Header.Get("token")
		if data, err := user.VerifyToken(token); err != nil {
			return ResError(400, err.Error())
		} else {
			return ResData(data)
		}
	*/

	if ok := user.CheckPasswd(form.Passwd, form.Salt); ok {
		if err := user.Logined(ctx); err != nil {
			return ResError(401, err.Error())
		}
		return ResData(&user)
	}
	return ResError(400, "bad request")
}
