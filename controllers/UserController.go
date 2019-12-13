// UserController.go kee > 2019/12/11
// koobeton 用户管理

package controllers

import (
	"github.com/keesely/kiris/hash"
	"koobeton/app"
	"koobeton/app/request"
	"koobeton/models"
)

type UserController struct{}

// 创建用户
func (*UserController) Post(ctx app.Context) Result {
	var user models.User
	ctx.ReadJSON(&user)
	if affected, err := user.Create(); 0 >= affected {
		return ResError(500, err.Error())
	}
	return ResData(app.Json{
		"data": user,
	})
}

func (*UserController) Get(ctx app.Context) Result {
	var users = []models.User{}
	M.Find(&users)
	return ResData(&users)
}

func (*UserController) GetBy(id int, ctx app.Context) Result {
	user := models.User{Id: id}
	if ok, _ := user.Get(); ok {
		return ResData(&user)
	}
	return ResError(404, "user no exists")
}

func (*UserController) GetCreate() app.Json {
	return app.Json{
		"action": "create user",
	}
}

// Login
func (*UserController) PostLogin(ctx app.Context) Result {
	req := request.New(ctx)

	user := models.User{Email: req.Get("email").(string)}

	passwd := req.Get("passwd").(string)
	salt := "1234b"

	passwd = hash.Sha1(hash.Sha1(passwd) + salt)
	if ok := user.CheckPasswd(passwd, salt); ok {
		remotAddr := app.GetClientIP(ctx)
		if err := user.Logined(remotAddr); err != nil {
			return ResError(401, err.Error())
		}
		return ResData(&user)
	}
	return ResError(400, "bad request")
}
