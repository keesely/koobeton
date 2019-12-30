// init.go kee > 2019/12/11

package models

import (
	"koobeton/app"
	//"log"
	"fmt"
	"time"
)

// 分类表
type Category struct {
	Id       int    `json:"id" xorm:"int(11) not null autoincr pk"`
	Name     string `json:"name"`
	Status   int    `json:"status" xorm:"tinyint"`
	ParentId int    `json:"parent_id"`
	UserId   int    `json:"user_id"`
}

type Author struct {
	Id   int    `json:"id" xorm:"int(11) not null autoincr pk"`
	Name string `json:"name"`
}

type Document struct {
	Id          int       `json:"id" xorm:"int(11) not null autoincr pk"`
	Title       string    `json:"title"`                     // 文本标题
	Intro       string    `json:"intro"`                     // 简介-max:255
	Tags        []string  `json:"tags"`                      // 标签列表 - multiple
	Authors     []Author  `json:"author"`                    // 编辑者 - multiple
	Content     string    `json:"content"`                   // 正文
	EditContent string    `json:"edit_content"`              // 编辑内容
	CreatedAt   time.Time `json:"created_at" xorm:"created"` // 创建时间
	ReleasedAt  time.Time `json:"released_at"`               // 发布时间
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated"` // 更新时间
}

var orm = app.NewOrm()

func init() {
	app.CreateModels([]interface{}{
		new(User),
		new(UserToken),
		new(Category),
		new(Author),
		new(Document),
		new(Workspace),
	})
}

func errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}
