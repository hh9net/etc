package storage

import (
	"CenterSettlement-go12/sysinit"
	"log"
	"testing"
)

//测试查询本省的结算数据
func TestQueryJiessj(t *testing.T) {
	QueryJiessjcz()
}
func TestNewDatabase(t *testing.T) {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	log.Println(xorm)
}

//测试查询其他地区的结算数据
func TestQueryQiTaJiessj(t *testing.T) {
	//QueryQiTaJiessj()
}
