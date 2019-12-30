// UserController.go kee > 2019/12/11
// koobeton 用户管理

package controllers

import (
	"koobeton/app"
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

// 获取用户列表
func (*UserController) Get(ctx app.Context) Result {
	var users = []models.User{}
	M.Find(&users)
	return ResData(&users)
}

// 获取用户信息
func (*UserController) GetBy(id int, ctx app.Context) Result {
	user := models.User{Id: id}
	if ok, _ := user.Get(); ok {
		return ResData(&user)
	}
	return ResError(404, "user no exists")
}
