package server

import (
	"github.com/go-xorm/xorm"
	"log"
)

const pwd = "shx19930509321"

type DB struct {
	orm *xorm.Engine
}

func NewDatabase() *DB {
	var db DB
	xo, err := xorm.NewEngine("mysql", "root:"+pwd+"@tcp(127.0.0.1:3306)/center?charset=utf8")
	if err != nil {
		log.Println("lianjie 失败")
		return nil
	}
	db.orm = xo
	log.Println("lianjie成功", db.orm)
	return &db
}

//创建数据库
func (db *DB) NewTable() {
	db = NewDatabase()
	is, err := db.orm.IsTableEmpty(
		new(SJsJiessj),
	)
	if err != nil {

	}

	if is == false {
		//err := engine.Sync2(new(User), new(Group))
		err = db.orm.CreateTables(new(SJsJiessj))
		if err != nil {

		}
	}
}
