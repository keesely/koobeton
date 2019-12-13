// model.go kee > 2019/12/11

package app

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	_ "github.com/mattn/go-sqlite3"
	//"reflect"
	//"log"
)

var (
	Engine *xorm.Engine
)

func NewEngine(driver, DSN string) (*xorm.Engine, error) {
	orm, err := xorm.NewEngine(driver, DSN)
	return orm, err
}

func NewOrm() *xorm.Engine {
	if nil == Engine {
		var (
			DSN string
			err error
		)
		driver := Config.Get("db.connection", "sqlite3").(string)
		switch driver {
		case "mysql":
			DSN = getMysqlDSN()
			break
		case "sqlite3":
			DSN = getSqlite3DSN()
			break
		default:
			driver = "sqlite3"
			DSN = getSqlite3DSN()
			break
		}
		Engine, err = NewEngine(driver, DSN)
		if err != nil {
			app.Logger().Fatal(err)
		}

		iris.RegisterOnInterrupt(func() {
			Engine.Close()
		})
	}
	return Engine
}

func CreateModels(models []interface{}) {
	var (
		err error
		orm = NewOrm()
	)
	for _, model := range models {
		switch model.(type) {
		default:
			{
				if ok, _ := orm.IsTableExist(model); !ok {
					err = orm.Sync2(model)
				}
			}
		}
		if err != nil {
			app.Logger().Fatal(err)
		}
	}
}

//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func getMysqlDSN() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s",
		Config.Get("db.user", "root"),
		Config.Get("db.password", "").(string),
		Config.Get("db.protocol", "tcp"),
		Config.Get("db.host", "localhost"),
		Config.Get("db.port", 3306).(int),
		Config.Get("db.database", "test"),
		Config.Get("db.charset", "utf8"),
	)
}

func getSqlite3DSN() string {
	dsn := fmt.Sprintf("file:%s", Config.Get("db.database", "./test.db"))
	if options := Config.Get("db.options", "").(string); options != "" {
		dsn += "?" + options
	}
	return dsn
}
