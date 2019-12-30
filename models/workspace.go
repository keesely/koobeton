// workspace.go kee > 2019/12/27

package models

import (
	"time"
)

type Workspace struct {
	Id          int       `json:"id" xorm:"int(11) not null autoincr pk"`
	Name        string    `json:"name" xorm:"notnull unique"`  // 工作区名称
	Token       string    `json:"token" xorm:"notnull unique"` // 工作区地址标识
	UserId      int       `json:"user_id" xorm:"notnull"`      // 关联创建者ID, 所有权凭证
	Email       string    `json:"email"`                       // 工作区公开邮箱、虽然暂时不知道有什么用，就先添加了吧
	Description string    `json:"description" xorm:"text"`     // 工作区描述
	Location    string    `json:"location"`                    // 所在位置，看情况，喜欢就添加呗
	Cover       string    `json:"cover"`                       // 封面地址
	Publisher   string    `json:"publisher"`                   // 公司标识，出版机构标识 - 用以导出时的页脚设置
	CreatedAt   time.Time `json:"created_at" xorm:"created"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at" xorm:"updated"`   // 更新时间
	DeletedAt   time.Time `json:"deleted_at" xorm:"deleted"`   // 移除时间 / 标识数据有效性唯一标识 / NULL = 有效
}

func (m Workspace) Create() (int, error) {
	if 0 >= m.UserId {
		return -1, errorf("The owner user id is empty")
	}
	u := User{Id: m.UserId}
	if ok, _ := orm.Get(&u); !ok {
		return -2, errorf("The owner user is fail")
	}
	var ws = new(Workspace)
	if ok, _ := orm.Where("name=?", m.Name).Or("token=?", m.Token).Get(ws); ok {
		if ws.Token == m.Token {
			return -3, errorf("the workspace token already exists.")
		}
		return -3, errorf("the workspace has already exists.")
	}

	if ins, err := orm.Insert(&m); ins > 0 {
		return int(ins), nil
	} else {
		return 0, err
	}
	return 0, errorf("insert fail")
}

func (m Workspace) CheckToken(token string) bool {
	if ok, _ := orm.Get(&Workspace{Token: token}); ok {
		return true
	}
	return false
}

// 设置/转移所有人
func (m *Workspace) SetOwer(user *User) error {
	m.UserId = user.Id
	orm.Update(&m)
	return errorf("set the owner fail")
}

// 获取创建者ID
func (m Workspace) GetOwer() (*User, error) {
	if 0 >= m.UserId {
		return nil, errorf("owner user id as empty")
	}
	var user = User{Id: m.UserId}
	if ok, _ := orm.Get(&user); ok {
		return &user, nil
	}
	return nil, errorf("owner user no exist")
}
