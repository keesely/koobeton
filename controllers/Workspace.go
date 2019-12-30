// Workspace.go kee > 2019/12/28

package controllers

import (
	"koobeton/app"
	m "koobeton/models"
)

type Workspace struct{}

func (w Workspace) Post(ctx app.Context) (result Result) {
	var ws = new(m.Workspace)
	ctx.ReadJSON(ws)

	if user := w.getUser(ctx); user != nil {
		ws.UserId = user.Id
	}
	if _, werr := ws.Create(); werr != nil {
		return ResError(500, "create workspace fail: "+werr.Error())
	}
	return ResData(&ws)
}

func (w Workspace) Get(ctx app.Context) Result {
	var rows = make([]m.Workspace, 0)
	orm := app.Engine.NewSession()
	if user := w.getUser(ctx); user != nil {
		orm.Where("user_id=?", user.Id)
	}

	if err := orm.Find(&rows); err != nil {
		return ResError(500, err.Error())
	}
	return ResData(rows)
}

func (w Workspace) GetBy(name string) Result {
	ws := m.Workspace{Name: name}
	if ok, _ := app.Engine.Get(&ws); ok {
		return ResData(&ws)
	}
	return ResError(500, "workspace get fail")
}

func (w Workspace) getUser(ctx app.Context) *m.User {
	var user, cUser = m.User{}, ctx.Params().Get("cUser")
	new(app.Json).JsonUnmarshal(cUser, &user)
	return &user
}
