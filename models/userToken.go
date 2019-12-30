// userToken.go kee > 2019/12/15

package models

import (
	"fmt"
	"github.com/keesely/kiris/hash"
	"strconv"
	"time"
)

// 用户token存储表
type UserToken struct {
	Id        int       `json:"id" xorm:"int(11) not null autoincr pk"`
	UserId    int       `json:"user_id"`
	AccessKey string    `json:"access_key" xorm:"varchar(32) not null unique"`
	SecretKey string    `json:"secret_key" xorm:"varchar(32) not null"`
	Expire    int       `json:"expire"` // 到期时间
	Status    int       `json:"status" xorm:"tinyint default 1 comment('0:无效;1:在用;2:临时')"`
	Device    string    `json:"device" xorm:"default 'unknown'"` // 设备标识
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	DeletedAt time.Time `json:"deleted_at" xorm:"deleted ->"`
}

// 获取凭证内容
func (m *UserToken) GetToken(accessKey string) *UserToken {
	if ok, _ := orm.Where("access_key=?", accessKey).Get(m); ok {
		if m.Status == 1 {
			return m
		} else if m.Status == 2 && (time.Now().Unix() < int64(m.Expire)) {
			return m
		}
	}
	return nil
}

// 校验访问凭证
func (m *UserToken) CheckToken(token string, pk ...string) (*User, error) {
	if ok := m.GetToken(m.AccessKey); ok != nil {
		s := fmt.Sprintf("%s:%s", m.AccessKey, m.SecretKey)
		for _, k := range pk {
			s += k
		}
		sToken := hash.Sha1(s)
		if sToken == token {
			return m.GetUser()
		}
		return nil, fmt.Errorf("bad token")
	}
	return nil, fmt.Errorf("bad token (access_key has been disabled)")
	//return nil, fmt.Errorf("bad token (access key not found)")
}

// 获取AccessKey 对应用户
func (m *UserToken) GetUser() (*User, error) {
	var user = &User{Id: m.UserId}
	ok, err := user.Get(true)
	if ok {
		return user, nil
	} else if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("bad token (user not found)")
}

// 新增访问凭证
// param device: 传入设备信息
// param ttl: time to live 凭证有效时长/hour / 0值则永久有效
func (m *UserToken) AddToken(device string, ttl int) (*UserToken, error) {
	// 校验用户
	if _, err := m.GetUser(); err != nil {
		return nil, err
	}

	var (
		accessKey = hash.Md5(strconv.Itoa(m.Id) + strconv.Itoa(int(time.Now().Unix())))
		secretKey = hash.Md5(strconv.Itoa(int(time.Now().Unix())) + device)
	)
	token := UserToken{
		UserId:    m.Id,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Device:    device,
		Status:    1,
	}
	if ttl > 0 {
		token.Status = 2
		token.Expire = int(time.Now().Add(time.Hour * time.Duration(ttl)).Unix())
	}
	if affected, err := orm.Insert(&token); affected > 0 {
		if err != nil {
			return nil, err
		}
		return &token, nil
	}
	return nil, fmt.Errorf("add Token fail")
}
