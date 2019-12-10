package model

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
)

type BaseModel struct {
	Id        int        `xorm:"id pk autoincr index comment(默认主键)" json:"id"`
	CreatedAt time.Time  `xorm:"created timestamp comment(创建时间) 'created_at'" json:"created"`
	UpdatedAt time.Time  `xorm:"updated timestamp comment(更新时间) 'updated_at'" json:"updated,omitempty" `
	DeletedAt *time.Time `xorm:"deleted timestamp comment(删除时间) 'deleted_at'" json:"deleted,omitempty"`
	Version   int        `xorm:"version comment(版本)" json:"version"`
}

var db *xorm.Engine
var conn string
var BufferSize = 100
var Limit = 10

func InitMysql() {
	var e error

	database := viper.Get("mysql.database")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	param := "?"
	loc := url.QueryEscape("Asia/Shanghai")
	conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8mb4&parseTime=true", user, password, host, port, database, param, loc)

	db, e = NewDatabase("mysql", conn)

	if viper.GetBool("system.debug") {
		db.ShowSQL()
	}

	if e != nil {
		panic(e)
	}
}

func NewDatabase(driver, source string) (db *xorm.Engine, err error) {
	db, err = xorm.NewEngine(driver, source)
	return db, err
}

func DB() *xorm.Engine {
	var e error
	if db == nil || db.Ping() != nil {
		db, e = NewDatabase("mysql", conn)
		if e != nil {
			panic(e)
		}
	}
	return db
}

func WhereGet(v interface{}, query interface{}, args ...interface{}) (bool, error) {
	return DB().Where(query, args...).Get(v)
}
