// struct.go kee > 2019/12/11

package models

import (
	"koobeton/app"
	//"log"
	"time"
)

type User struct {
	Id        int       `json:"id" xorm:"int(11) not null autoincr pk"`
	Name      string    `json:"name" xorm:"index"`
	Passwd    string    `json:"passwd,omitempty" xorm:"-> 'passwd' varchar(64) not null"`
	Status    int       `json:"status" xorm:"tinyint(4) not null default 1 comment('0: 停用;1: 启用')"`
	Sex       int       `json:"sex" xorm:"tinyint(4) 'sex' default 0 comment('0:保密;1:男;2:女')"`
	Nicename  string    `json:"nicename"`
	Topic     string    `json:"topic" xorm:"comment('头像链接地址')"`
	Email     string    `json:"email" xorm:"not null unique"`
	Phone     string    `json:"phone" xorm:"varchar(11)"`
	LastedIP  string    `json:"lasted_ip" xorm:"varchar(255) 'lasted_ip'"`
	LastedAt  time.Time `json:"lasted_at" xorm:"datetime"`
	CreatedAt time.Time `json:"created_at" xorm:"created default now()"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted"`
}

type Category struct {
	Id       int    `json:"id" xorm:"int(11) not null autoincr pk"`
	Name     string `json:"name"`
	State    int    `json:"state"`
	ParentId int    `json:"parent_id"`
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
		new(Category),
		new(Author),
		new(Document),
	})
}
