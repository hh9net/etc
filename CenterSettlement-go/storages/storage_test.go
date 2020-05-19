package storage

import (
	"CenterSettlement-go/sysinit"
	"log"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	//log.Println(xorm)
	if xorm == nil {
		log.Println("xorm ==nil :", err)
	}
	log.Println("xorm 获取 成功")
}

//测试查询本省的结算数据 储值卡 打包数据状态为0  卡网络号为江苏
func TestQueryJiessjcz(t *testing.T) {
	QueryJiessjcz()
}

//测试查询本省的结算数据 记账卡 打包数据状态为0  卡网络号为江苏
func TestQueryJiessjjz(t *testing.T) {
	QueryJiessjjz()
}

//测试查询其他地区的结算数据
func TestQueryQiTaJiessj(t *testing.T) {

}
