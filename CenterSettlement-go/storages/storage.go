package storage

import (
	"CenterSettlement-go/sysinit"
	"CenterSettlement-go/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

//数据层负责 查询数据、插入数据、准备数据、返回数据、错误处理
//在收到联网中心的数据后，解析数据，插入数据
//在发送给联网中心前查询数据
//注意事务处理

//通过交易状态为0，卡网络号为江苏本省， 卡类型为储值卡、查询交易结算数据
func QueryJiessjcz() *[]types.BJsJiessj {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	session := TransactionBegin(xorm)
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	qerr := session.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", types.JS_NETWORK).And("F_NB_KALX=?", types.PRECARD).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}
	terr := TransactionCommit(session)
	if terr != nil {
		log.Fatal("查询本省交易记录事务提交错误")
	}
	return &tests
}

//通过交易状态为0，卡网络号为江苏本省， 卡类型为记账卡、查询交易结算数据
func QueryJiessjjz() *[]types.BJsJiessj {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	//开启事务
	session := TransactionBegin(xorm)
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	qerr := session.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", types.JS_NETWORK).And("F_NB_KALX=?", types.CREDITCARD).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}

	//关闭事务
	terr := TransactionCommit(session)
	if terr != nil {
		log.Fatal("查询本省交易记录事务提交错误")
	}
	return &tests
}

type QiTaJiessj struct {
	Network string
	KalX    int
	Jiessj  []types.BJsJiessj
}

var (
	jiessjcz []types.BJsJiessj //储值卡
	jiessjjz []types.BJsJiessj //记账卡
	jiessjs  []QiTaJiessj      //返回的数据
)

//通过交易状态0  卡网络号为（除江苏以外） 卡类型为储值卡 或者记账卡   查询交易结算数据
//map[string]map[string] []types.BJsJiessj 第一级是数据为记账卡 第二级数据是卡网络号
func QueryQiTaJiessj() *QiTaJiessj {
	xorm, err := sysinit.NewEngine()
	if err != nil {
		log.Fatal("连接数据库 error :", err)
	}
	Qitajiessj := new(QiTaJiessj)

	//开启事务
	session := TransactionBegin(xorm)
	//查询多条数据
	//networkJiessjs := make(map[string][]types.BJsJiessj, 0)
	//netsjcz:=	make(map[int]map[string] []types.BJsJiessj)

	//把交易状态为0  卡网络号为其他地区的数据记录查询出来  卡类型22/23
	for _, networkcode := range types.Gl_network {
		//查询 为储值卡的数据
		jiessjcz = getJiessj(session, networkcode, 22)
		if jiessjcz == nil {
			Qitajiessj.Network = networkcode
			Qitajiessj.KalX = 22
		}

		for _, jioayisj := range jiessjcz {

			Qitajiessj.Network = networkcode
			Qitajiessj.KalX = 22
			for _, sjs := range jiessjs {
				sjs.Jiessj = append(sjs.Jiessj, jioayisj)
				sjs.KalX = Qitajiessj.KalX
				sjs.Network = Qitajiessj.Network
			}

		}

		//查询 为记账卡的数据
		jiessjjz = getJiessj(session, networkcode, 23)
		if jiessjjz == nil {

		}

		//log.Println(networkJiessj)

	}
	//log.Println(networkJiessjs)
	//return networkJiessjs
	terr := TransactionCommit(session)
	if terr != nil {
		log.Fatalln("QueryQiTaJiessj TransactionCommit error")
	}

	return Qitajiessj
}

func getJiessj(session *xorm.Session, networkcode string, Kalx int) []types.BJsJiessj {
	tests := make([]types.BJsJiessj, 0)
	qerr := session.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", networkcode).And("F_NB_KALX=?", Kalx).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("卡网络号为 %s 总共查询出 %d 条数据\n", networkcode, len(tests))
	for _, v := range tests {
		log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}
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
	qerr := xorm.Where("F_NB_DABZT=?", 0).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("消息包序号: %d\n", v.FNbXiaoxxh)
	}
	return &tests
}
