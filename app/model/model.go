// model.go kee > 2019/12/16

package model

import (
	"log"
	"reflect"
)

type Model struct {
}

func New(m interface{}) interface{} {
	//valueOf := reflect.ValueOf(model)
	//typeof := reflect.TypeOf(model)
	elem := reflect.TypeOf(m).Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		mTag := field.Tag.Get("model")
		oTag := field.Tag.Get("xorm")
		log.Printf("%s: %s / %s \n", field.Name, mTag, oTag)
	}
	return m
}
