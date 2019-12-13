// user.go kee > 2019/12/11

package models

import (
	"fmt"
	"github.com/keesely/kiris/hash"
	"koobeton/app"
	"time"
)

func (m *User) Logined(lastIP string) error {
	m.LastedIP = lastIP
	m.LastedAt = time.Now()
	_, err := orm.Update(m)
	return err
}

func (m *User) CheckPasswd(passwd, salt string) (flag bool) {
	if ok, _ := orm.Where("email=?", m.Email).Get(m); ok {
		if hash.Sha1(hash.Sha1(passwd)+salt) == hash.Sha1(m.Passwd+salt) {
			return true
		}
		return false
	}
	return
}

// 创建用户
// return -1: user exists 0: error >0: insert successful
func (m *User) Create() (int, error) {
	if m.Email == "" {
		return 0, fmt.Errorf("email is empty")
	}
	if m.Passwd == "" {
		return 0, fmt.Errorf("password is empty")
	}
	if ok, _ := orm.Get(&User{Email: m.Email}); !ok {
		m.Passwd = hash.Sha1(m.Passwd)
		if m.Name == "" {
			m.Name = m.Email
		}
		if m.Nicename == "" {
			m.Nicename = m.Name
		}
		if affected, err := orm.Insert(m); affected > 0 {
			return int(affected), nil
		} else {
			app.Logger().Error(err)
			return 0, err
		}
	} else {
		return -1, fmt.Errorf("The user already exists: %s", m.Email)
	}
	return 0, fmt.Errorf("Undefined fail")
}

func (m *User) Get() (ok bool, err error) {
	return orm.Get(m)
}
