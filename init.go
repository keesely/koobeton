// init.go kee > 2019/12/15
// 全局初始化方法
// 处理应用程序初始化数据问题
// 依赖 .lock

package main

import (
	"fmt"
	"github.com/keesely/kiris"
	"github.com/keesely/kiris/hash"
	"koobeton/app"
	"koobeton/models"
	"log"
	"strconv"
	"time"
)

func init() {
	lock := kiris.RealPath("./.lock")
	if kiris.FileExists(lock) {
		return
	}
	app.NewApp()
	securitySet(lock)
	adminUserSet()
}

// 安全配置
// 生成加密密钥，用于加密各式隐私数据
func securitySet(lock string) {
	c := app.Config
	s1 := hash.Md5(strconv.Itoa(int(time.Now().UnixNano())))[:8]
	s2 := hash.Md5(c.Get("app.name", "koobeton").(string))[:8]
	s3 := hash.Md5(strconv.Itoa(kiris.Rand(123456789, 987654321)))[:8]
	//s4 := hash.Md5()
	uniqid := kiris.Base64Encode([]byte(hash.Md5(s1 + s2 + s3)))
	kiris.FilePutContents(lock, fmt.Sprintf("base64:%s", uniqid), 0)
}

// 超级用户初始化
// 根据配置文件初始化系统超级用户
func adminUserSet() {
	c := app.Config
	name := c.Get("admin.name", "admin").(string)
	nicename := "administrator"
	email := c.Get("admin.email", "example@koobeton.com").(string)
	passwd := c.Get("admin.passwd", "12345678").(string)
	var user = models.User{Name: name, Nicename: nicename, Email: email, Passwd: passwd}
	if _, err := user.Create(); err != nil {
		log.Panic(err)
	}
}
