package storage

import (
	"CenterSettlement-go/sysinit"
	"CenterSettlement-go/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

//数据层负责查询数据、准备数据、

//通过交易状态为0，卡网络号为江苏本省， 卡类型为储值卡、   查询交易结算数据
func QueryJiessj() *[]types.BJsJiessj {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	qerr := xorm.Where("F_NB_JIAOYZT=?", 0).And("F_VC_KAWLH=?", types.JS_NETWORK).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("交易状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbJiaoyzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}
	return &tests
}

//通过交易状态  卡网络号为（除江苏以外） 查询交易结算数据
func QueryQiTaJiessj() []types.QiTaJiessj {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	//查询多条数据
	networkJiessjs := make([]types.QiTaJiessj, 0)

	//var networkcode string
	var networkJiessj types.QiTaJiessj
	//把交易状态为0 卡网络号为其他地区的数据记录查询出来
	for _, networkcode := range types.Gl_network {
		networkJiessj.Networkcode = networkcode
		//查询 code de数据
		jiessj := getJiessj(xorm, networkcode)

		networkJiessj.QitaJiessj = jiessj
		//log.Println(networkJiessj)
		networkJiessjs = append(networkJiessjs, networkJiessj)
	}
	log.Println(networkJiessjs)
	return networkJiessjs
}

func getJiessj(xorm *xorm.Engine, networkcode string) []types.BJsJiessj {
	tests := make([]types.BJsJiessj, 0)
	qerr := xorm.Where("F_NB_JIAOYZT=?", 0).And("F_VC_KAWLH=?", networkcode).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	//log.Printf("卡网络号为 %s 总共查询出 %d 条数据\n", networkcode,len(tests))
	//for _, v := range tests {
	// log.Printf("交易状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbJiaoyzt,v.FVcJiaoyjlid,v.FVcKawlh)
	//}
	return tests
}

//   B_JS_YUANSJYXX【原始交易消息包表】
func YUANSJYXX() *[]types.BJsYuansjyxx {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	//查询多条数据
	tests := make([]types.BJsYuansjyxx, 0)
	qerr := xorm.Where("F_NB_JIAOYZT=?", 0).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("消息包序号: %d\n", v.FNbXiaoxxh)
	}
	return &tests
}
