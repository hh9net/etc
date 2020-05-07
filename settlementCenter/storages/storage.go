package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

func NewEngine() *xorm.Engine {
	//x, err = xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/xorm?charset=utf8")
	x, err := xorm.NewEngine("mysql", "root:mysql@/tcp(192.168.200.160:3306)/payflow?charset=utf8")
	if err != nil {
		log.Fatal("连接数据库error")
	}
	log.Fatal("连接数据库成功")
	return x
}

//查
//func Getinfo(id int64) *User {
//	user := &User{User_id: id}
//	is, _ := x.Get(user)
//	if !is {
//		log.Fatal("搜索结果不存在!")
//	}
//	return user
//}
