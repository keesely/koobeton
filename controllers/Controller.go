// Controller.go kee > 2019/12/11

package controllers

import (
	"koobeton/app"
)

type Controller struct {
	Result Result
}

type Result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Extra interface{} `json:"extra,omitempty"`
}

var (
	M = app.NewOrm()
)

func ResError(code int, msg string, extra ...interface{}) Result {
	return Result{code, msg, nil, getExtra(extra)}
}

func ResData(data interface{}, extra ...interface{}) Result {
	return Result{200, "successful", data, getExtra(extra)}
}

func getExtra(extra []interface{}) interface{} {
	if len(extra) > 0 {
		return extra
	}
	return nil
}
