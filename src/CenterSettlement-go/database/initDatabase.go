package database

import (
	"CenterSettlement-go/conf"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var XormClient *xorm.Engine

//连接数据库
func DBInit() {
	config := conf.ConfigInit()
	params := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", config.MUserName, config.MPass, config.MHostname, config.MPort, config.Mdatabasename)
	x, err := xorm.NewEngine("mysql", params)
	if x == nil {
		log.Println("获取xorm为空")
		x = new(xorm.Engine)
	}
	if err != nil {
		log.Fatal("连接数据库error")
	}
	log.Println("连接数据库成功")
	if XormClient == nil {
		XormClient = new(xorm.Engine)
	}
	XormClient = x
}
