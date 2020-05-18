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

//通过交易状态0  卡网络号为（除江苏以外） 卡类型为储值卡 或者记账卡   查询交易结算数据
//map[string]map[string] []types.BJsJiessj 第一级是数据为记账卡 第二级数据是卡网络号
//func QueryQiTaJiessj() map[int]map[string] []types.BJsJiessj{
//	xorm, err := sysinit.NewEngine()
//	if err != nil {
//		log.Fatal("连接数据库 error :", err)
//	}
//	//查询多条数据
//	networkJiessjs := make(map[string][]types.BJsJiessj, 0)
//	netsjcz:=	make(map[int]map[string] []types.BJsJiessj)
//
//
//
//	//把交易状态为0  卡网络号为其他地区的数据记录查询出来  卡类型22/23
//	for _, networkcode := range types.Gl_network {
//
//		//查询 code de数据
//		switch networkcode {
//			//SH_NETWORK  string = "3101" // 上海
//			case types.SH_NETWORK :
//				shanghaijiessjcz := getJiessj(xorm, networkcode,22)
//				//netsj[22]=networkJiessjs[types.SH_NETWORK]
//				shanghaijiessjjz := getJiessj(xorm, networkcode,23)
//
//
//			//ZJ_NETWORK  string = "3301" // 浙江
//			//AH_NETWORK  string = "3401" // 安徽
//			//FJ_NETWORK  string = "3501" // 福建
//			//JX_NETWORK  string = "3601" // 江西
//			//SD_NETWORK  string = "3701" // 山东
//			//SD_NETWORK2 string = "3702" // 山东
//			//
//			///* 华北区路网代码定义*/
//			//BJ_NETWORK  string = "1101" // 北京
//			//TJ_NETWORK  string = "1201" // 天津
//			//HEB_NETWORK string = "1301" // 河北
//			//SX_NETWORK  string = "1401" // 山西
//			//NM_NETWORK  string = "1501" // 内蒙古
//			//
//			///* 东北区路网代码定义*/
//			//LN_NETWORK  string = "2101" // 辽宁
//			//JL_NETWORK  string = "2201" // 吉林
//			//HLJ_NETWORK string = "2301" // 黑龙江
//			//
//			///*华中、华南区路网代码定义*/
//			//HEN_NETWORK  string = "4101" // 河南
//			//HUB_NETWORK  string = "4201" // 湖北
//			//HUB_NETWORK2 string = "4202" // 湖北
//			//
//			//HUN_NETWORK  string = "4301" // 湖南
//			//GD_NETWORK   string = "4401" // 广东
//			//GX_NETWORK   string = "4501" // 广西
//			//HAIN_NETWORK string = "4601" // 海南
//			//
//			///*西南区路网代码定义*/
//			//CQ_NETWORK  string = "5001" // 重庆
//			//SC_NETWORK  string = "5101" // 四川
//			//SC_NETWORK2 string = "5102" // 四川
//			//SC_NETWORK3 string = "5103" // 四川
//			//SC_NETWORK4 string = "5104" // 四川
//			//SC_NETWORK5 string = "5105" // 四川
//			//GZ_NETWORK  string = "5201" // 贵州
//			//YN_NETWORK  string = "5301" // 云南
//			//XZ_NETWORK  string = "5401" // 西藏
//			//
//			///*西北区路网代码定义*/
//			//SHANXI_NETWORK  string = "6101" // 陕西
//			//SHANXI_NETWORK2 string = "6102" // 陕西
//			//SHANXI_NETWORK3 string = "6103" // 陕西
//			//SHANXI_NETWORK4 string = "6104" // 陕西
//			//SHANXI_NETWORK5 string = "6105" // 陕西
//			//SHANXI_NETWORK6 string = "6106" // 陕西
//			//SHANXI_NETWORK7 string = "6107" // 陕西
//			//
//			//GS_NETWORK string = "6201" // 甘肃
//			//QH_NETWORK string = "6301" // 青海
//			//NX_NETWORK string = "6401" // 宁夏
//			//XJ_NETWORK string = "6501" // 新疆
//			//
//			//ARMY_CARDNETWORK string = "501" // 军车卡的网络编号
//			//
//		}
//		jiessjcz := getJiessj(xorm, networkcode,22)
//
//
//		//log.Println(networkJiessj)
//
//	}
//	log.Println(networkJiessjs)
//	return networkJiessjs
//}

func getJiessj(xorm *xorm.Engine, networkcode string, Kalx int) []types.BJsJiessj {
	tests := make([]types.BJsJiessj, 0)
	qerr := xorm.Where("F_F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", networkcode).And("F_NB_KALX=?", Kalx).Limit(100, 0).Find(&tests)
	if qerr != nil {
		panic(qerr)
	}
	//log.Printf("卡网络号为 %s 总共查询出 %d 条数据\n", networkcode,len(tests))
	//for _, v := range tests {
	// log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt,v.FVcJiaoyjlid,v.FVcKawlh)
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
