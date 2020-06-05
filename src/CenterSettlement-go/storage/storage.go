package storage

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

//数据层负责 查询数据、插入数据、准备数据、返回数据、错误处理
//在收到联网中心的数据后，解析数据，插入数据
//在发送给联网中心前查询数据
//注意事务处理
//通过交易状态为0，卡网络号为江苏本省， 卡类型为储值卡、查询交易结算数据

type DB struct {
	//Db database.XormClient
}

func QueryQitaJiessj(KaLx int, Diqu string) *[]types.BJsJiessj {

	//database.DBInit()
	xorm := database.XormClient
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	//qerr := xorm.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", Diqu).And("F_NB_KALX=?", KaLx).Limit(100, 0).Find(&tests)
	//测试 每次查两条
	qerr := xorm.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", Diqu).And("F_NB_KALX=?", KaLx).Limit(10, 0).Find(&tests)
	if qerr != nil {
		log.Fatalln("查询结算数据出错", qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}
	return &tests
}

//江苏本省
func QueryJiessj(KaLx int) *[]types.BJsJiessj {
	//database.DBInit()
	xorm := database.XormClient
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	qerr := xorm.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", types.JS_NETWORK).And("F_NB_KALX=?", KaLx).Limit(100, 0).Find(&tests)
	if qerr != nil {
		log.Fatalln("查询结算数据出错", qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("打包状态: %d, 交易记录id: %s, 卡网络号: %s\n", v.FNbDabzt, v.FVcJiaoyjlid, v.FVcKawlh)
	}
	return &tests
}

//通过交易记录id 更新打包状态为打包中
func UpdatePackaging(Jiaoyjlid []string) error {
	//database.DBInit()
	xorm := database.XormClient

	for _, id := range Jiaoyjlid {
		Jiessj := new(types.BJsJiessj)
		Jiessj.FNbDabzt = 1
		log.Printf("交易记录id:%s 打包状态更新为：1", id)
		_, err := xorm.Table("b_js_jiessj").Where("F_VC_JIAOYJLID=?", id).Update(Jiessj)
		if err != nil {
			log.Println("更新打包状态失败", err)
			return err
		}
	}
	log.Println("更新打包状态为：打包中 成功")

	return nil
}

//打包成功
//   新增打包记录【插入表b_js_yuansjyxx】
func PackagingRecordInsert(data types.BJsYuansjyxx) error {
	//database.DBInit()
	xorm := database.XormClient
	yuansjyxx := new(types.BJsYuansjyxx)
	log.Println("PackagingRecordInsert data : ", data)
	yuansjyxx.FVcBanbh = data.FVcBanbh             //版本号
	yuansjyxx.FNbXiaoxlb = data.FNbXiaoxlb         //消息类别
	yuansjyxx.FNbXiaoxlx = data.FNbXiaoxlx         //消息类型
	yuansjyxx.FVcFaszid = data.FVcFaszid           //发送者ID
	yuansjyxx.FVcJieszid = data.FVcJieszid         //接受者ID
	yuansjyxx.FNbXiaoxxh = data.FNbXiaoxxh         //消息序号【消息包号】
	yuansjyxx.FDtDabsj = data.FDtDabsj             // 打包时间
	yuansjyxx.FVcQingfmbr = data.FVcQingfmbr       //清分目标日
	yuansjyxx.FVcTingccqffid = data.FVcTingccqffid //停车场清分方ID
	yuansjyxx.FVcFaxfwjgid = data.FVcFaxfwjgid     //发行服务机构ID 0000000000000020
	yuansjyxx.FNbJilsl = data.FNbJilsl             //记录数量
	yuansjyxx.FNbZongje = data.FNbZongje           //总金额
	yuansjyxx.FVcXiaoxwjlj = data.FVcXiaoxwjlj     //消息文件路径

	_, err := xorm.Table("b_js_yuansjyxx").Insert(yuansjyxx)
	if err != nil {
		log.Println("新增原始交易数据消息包打包记录 error", err)
		return err
	}
	log.Println("新增原始交易数据消息包打包记录 成功")
	return nil
}

//   新增打包明细记录
func PackagingMXRecordInsert(mx []types.BJsYuansjymx) error {
	//database.DBInit()
	xorm := database.XormClient

	Yuansjymx := new(types.BJsYuansjymx)
	//赋值
	for _, v := range mx {
		Yuansjymx.FVcXiaoxxh = v.FVcXiaoxxh
		Yuansjymx.FNbBaonxh = v.FNbBaonxh
		Yuansjymx.FDtJiaoysj = v.FDtJiaoysj
		Yuansjymx.FNbJine = v.FNbJine
		Yuansjymx.FVcDingzjyxx = v.FVcDingzjyxx
		Yuansjymx.FVcJiaoybh = v.FVcJiaoybh
		Yuansjymx.FVcTingccmc = v.FVcTingccmc
		Yuansjymx.FNbTingfsc = v.FNbTingfsc

		Yuansjymx.FNbShoufcx = v.FNbShoufcx
		Yuansjymx.FNbSuanfbs = v.FNbSuanfbs
		Yuansjymx.FNbFuwlx = v.FNbFuwlx
		Yuansjymx.FVcZhangdsm = v.FVcZhangdsm
		Yuansjymx.FVcJiaoyxxxx = v.FVcJiaoyxxxx
		Yuansjymx.FNbKalx = v.FNbKalx
		Yuansjymx.FVcWanglbm = v.FVcWanglbm
		Yuansjymx.FVcKawlbh = v.FVcKawlbh
		Yuansjymx.FVcKancph = v.FVcKancph
		//yuansjymx.FVcKajyxh=v.ICCard.   //卡交易序号

		Yuansjymx.FNbJiaoyqye = v.FNbJiaoyqye
		Yuansjymx.FNbJiaoyhye = v.FNbJiaoyhye
		Yuansjymx.FVcTacm = v.FVcTacm
		Yuansjymx.FVcJiaoybs = v.FVcJiaoybs
		Yuansjymx.FVcZongdjh = v.FVcZongdjh
		Yuansjymx.FVcZongdjyxh = v.FVcZongdjyxh
		Yuansjymx.FVcObuwlbh = v.FVcObuwlbh
		Yuansjymx.FVcObuzt = v.FVcObuzt
		Yuansjymx.FVcObuncph = v.FVcObuncph

		log.Println("原始交易消息明细数据：", Yuansjymx)
		_, err := xorm.Table("b_js_yuansjymx").Insert(Yuansjymx)
		if err != nil {
			log.Println("新增打包明细记录 error", err)
			return err
		}
	}
	log.Printf("原始交易消息包%d中 明细数据有：%d 条 数据 ", Yuansjymx.FVcXiaoxxh, len(mx))
	log.Println("新增打包明细记录 成功")
	return nil
}

//更新数据    根据 包号 更新原始交易消息包的【发送状态   发送中】

func UpdateYuansjyxx(Mid int64) error {
	//database.DBInit()
	xorm := database.XormClient
	yuansjyxx := new(types.BJsYuansjyxx)
	yuansjyxx.FNbFaszt = 1
	_, err := xorm.Table("b_js_yuansjyxx").Where("F_NB_XIAOXXH=?", Mid).Update(yuansjyxx)
	if err != nil {
		log.Println(" 根据 包号 更新原始交易消息包的发送状态 为 ： 发送中 error", err)
		return err
	}
	log.Println(" 根据 包号 更新原始交易消息包的发送状态 为 ： 发送中  成功")
	return nil
}

// 原始交易消息包发送成功更新 发送状态 发送时间 发送成功后消息包的文件路径
func SendedUpdateYuansjyxx(Mid int64, fname string) (error, string) {
	//database.DBInit()
	xorm := database.XormClient
	log.Println(fname)
	yuansjyxx := new(types.BJsYuansjyxx)
	yuansjyxx.FNbFaszt = 2
	yuansjyxx.FDtFassj = time.Now()
	yuansjyxx.FVcXiaoxwjlj = "compressed_xml/" + fname
	_, err := xorm.Table("b_js_yuansjyxx").Where("F_NB_XIAOXXH=?", Mid).Update(yuansjyxx)
	if err != nil {
		log.Println(" 根据 包号 更新 发送状态 发送时间 发送成功后消息包的文件路径 error", err)
		return err, ""
	}
	log.Println(" 根据 包号 更新 发送状态 发送时间 发送成功后消息包的文件路径  成功")

	sj := common.DateTimeFormat(yuansjyxx.FDtFassj)
	return nil, sj
}

//   更新结算数据打包结果【打包状态：已打包、原始交易包号、包内序号、清分目标日】
func UpdateDataPackagingResults(Jiaoyjlid []string, Msgid int64, jiaoyisj *types.Message) error {
	//database.DBInit()
	xorm := database.XormClient
	for i, idstr := range Jiaoyjlid {
		Jiessj := new(types.BJsJiessj)
		Jiessj.FNbDabzt = 2
		Jiessj.FNbYuansjybxh = Msgid
		//Xuh, _ := strconv.Atoi(idstr)
		Jiessj.FNbJiaoybnxh = i + 1
		Jiessj.FVcQingfmbr = jiaoyisj.Body.ClearTargetDate //打包当天日期

		_, err := xorm.Table("b_js_jiessj").Where("F_VC_JIAOYJLID=?", idstr).Update(Jiessj)
		if err != nil {
			log.Println("更新打包状态失败", err)
			return err
		}
	}
	log.Println("更新打包状态为：已达包、原始交易包号、包内序号 成功")

	return nil
}

//   新增打包应答记录
func PackagingResRecordInsert(data types.BJsYuansjyydxx) error {
	//database.DBInit()
	xorm := database.XormClient

	Yuansjyydxx := new(types.BJsYuansjyydxx)
	//赋值
	_, err := xorm.Table("b_js_jiessj").Insert(Yuansjyydxx)
	if err != nil {
		log.Println("新增打包明细记录 error")
		return err
	}
	return nil
}
