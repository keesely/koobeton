// user.go kee > 2019/12/11

package models

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/keesely/kiris"
	"github.com/keesely/kiris/hash"
	"koobeton/app"
	"time"
)

// 用户信息表
type User struct {
	Id        int       `json:"id" xorm:"int(11) not null autoincr pk"`
	Name      string    `json:"name" xorm:"index"`
	Passwd    string    `json:"passwd,omitempty" model:"hide" xorm:"'passwd' -> varchar(64) not null"`
	Status    int       `json:"status" xorm:"tinyint(4) not null default 1 comment('0: 停用;1: 启用')"`
	Sex       int       `json:"sex" xorm:"tinyint(4) 'sex' default 0 comment('0:保密;1:男;2:女')"`
	Nicename  string    `json:"nicename"`
	Topic     string    `json:"topic" xorm:"comment('头像链接地址')"`
	Email     string    `json:"email" xorm:"not null unique"`
	Phone     string    `json:"phone" xorm:"varchar(11)"`
	LastedIP  string    `json:"lasted_ip" xorm:"varchar(255) 'lasted_ip'"`
	LastedAt  time.Time `json:"lasted_at" xorm:"datetime"`
	CreatedAt time.Time `json:"created_at" xorm:"created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"updated"`
	DeletedAt time.Time `json:"-" xorm:"deleted ->"`
	Token     string    `json:"token,omitempty" xorm:"-"`
}

// 校验登陆密码
func (m *User) CheckPasswd(passwd, salt string) (ok bool) {
	if exists, _ := m.Get(true); exists {
		if cPasswd, err := m.GetPasswd(); err == nil {
			return passwd == hash.Sha1(cPasswd+salt)
		}
	}
	return
}

// 获取密码
func (m *User) GetPasswd() (string, error) {
	mapper := make(map[string]string)
	if ok, _ := orm.Table("user").Cols("passwd").ID(m.Id).Get(&mapper); ok {
		return mapper["passwd"], nil
	}
	return "", fmt.Errorf("user no exists")
}

// 生成新的Session信息 - JWT session
func (m *User) NewToken(device string) (string, error) {
	// 创建临时凭证
	accToken, err := (&UserToken{Id: m.Id}).AddToken(device, 2)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        m.Id,
		"accessKey": accToken.AccessKey,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(accToken.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *User) VerifyToken(token string) (app.Json, error) {
	tParse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// 验证
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		var (
			// conver to json data
			data = new(app.Json).Conver(token.Claims)
			ak   = data.Get("accessKey").(string)
			id   = int(data.Get("id").(float64))
		)

		// 提取SecretKey
		acc := new(UserToken).GetToken(ak)
		if acc == nil || acc.UserId != id {
			return nil, fmt.Errorf("bad token")
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(acc.SecretKey), nil
	})
	dJson := new(app.Json).Conver(tParse)
	return dJson, err
}

// 执行登陆后数据回写 - 创建JWT SESSION
func (m *User) Logined(ctx app.Context) error {
	lastIP := app.GetClientIP(ctx)
	ua := ctx.Request().Header.Get("User-Agent")
	token, err := m.NewToken(lastIP + ":" + ua)
	if err != nil {
		return err
	}
	m.LastedIP = lastIP
	m.LastedAt = time.Now()
	_, err = orm.Cols("lasted_ip", "lasted_at").Update(m)
	m.Token = token
	return err
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
		m.Name = kiris.Ternary(m.Name == "", m.Email, m.Name).(string)
		m.Nicename = kiris.Ternary(m.Nicename == "", m.Name, m.Nicename).(string)
		m.Status = 1

		if affected, err := orm.Insert(m); affected > 0 {
			return int(affected), nil
		} else {
			return 0, err
		}
	} else {
		return -1, fmt.Errorf("The user already exists: %s", m.Email)
	}
	return 0, fmt.Errorf("Undefined fail")
}

// 获取用户信息
// active = true 判定用户是否有效
func (m *User) Get(active ...bool) (ok bool, err error) {
	if active[0] {
		if ok, err = orm.Get(m); ok {
			return m.Active()
		}
		return false, fmt.Errorf("user no exists")
	}
	return orm.Get(m)
}

// 检测状态是否存活
func (m *User) Active() (bool, error) {
	if m.Status != 0 {
		return true, nil
	}
	return false, fmt.Errorf("the user has been disabled")
}
