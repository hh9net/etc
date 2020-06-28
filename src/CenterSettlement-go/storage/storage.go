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
//通过交易状态为0，卡网络号为地区， 卡类型为储值卡、记账卡，查询出交易结算数据

//查询其他省的结算数据
func QueryQitaJiessj(KaLx int, Diqu string) *[]types.BJsJiessj {

	xorm := database.XormClient
	//查询多条数据
	tests := make([]types.BJsJiessj, 0)
	// 每次查100条
	qerr := xorm.Where("F_NB_DABZT=?", 0).And("F_VC_KAWLH=?", Diqu).And("F_NB_KALX=?", KaLx).Limit(100, 0).Find(&tests)
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
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	for _, id := range Jiaoyjlid {
		Jiessj := new(types.BJsJiessj)
		Jiessj.FNbDabzt = 1
		log.Printf("交易记录id:%s 打包状态更新为：1", id)
		_, err := session.Table("b_js_jiessj").Where("F_VC_JIAOYJLID=?", id).Update(Jiessj)
		if err != nil {
			log.Println("更新打包状态 失败", err)
			return err
		}
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("更新打包状态为：打包中 时，事务错误", serr)
		return serr
	}
	log.Println("更新打包状态为：打包中 成功")
	return nil
}

//通过交易记录id 更新打包状态为 初始状态0
func UpdatePackagingInit(Jiaoyjlid []string) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	for _, id := range Jiaoyjlid {
		Jiessj := new(types.BJsJiessj)
		Jiessj.FNbDabzt = 0
		log.Printf("交易记录id:%s 打包状态更新为：0", id)
		_, err := session.Table("b_js_jiessj").Where("F_VC_JIAOYJLID=?", id).Update(Jiessj)
		if err != nil {
			log.Println("更新打包状态为：0 失败", err)
			return err
		}
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("更新打包状态为：初始未打包 时，事务错误", serr)
		return serr
	}
	log.Println("更新打包状态为：初始未打包 成功")

	return nil
}

//打包成功
//   新增打包记录【插入 原始交易消息表 b_js_yuansjyxx】
func PackagingRecordInsert(data types.BJsYuansjyxx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)
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

	_, err := session.Table("b_js_yuansjyxx").Insert(yuansjyxx)
	if err != nil {
		log.Println("新增原始交易数据消息包打包记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增原始交易数据消息包打包记录 时，事务错误", serr)
		return serr
	}
	log.Println("新增原始交易数据消息包打包记录 成功")
	return nil
}

//   新增打包明细记录
func PackagingMXRecordInsert(mx []types.BJsYuansjymx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	Yuansjymx := new(types.BJsYuansjymx)
	//赋值
	for _, v := range mx {
		Yuansjymx.FVcXiaoxxh = v.FVcXiaoxxh     //消息序号
		Yuansjymx.FNbBaonxh = v.FNbBaonxh       //包内序号
		Yuansjymx.FDtJiaoysj = v.FDtJiaoysj     //交易时间
		Yuansjymx.FNbJine = v.FNbJine           //金额
		Yuansjymx.FVcDingzjyxx = v.FVcDingzjyxx //定制交易信息 CustomizedData
		Yuansjymx.FVcJiaoybh = v.FVcJiaoybh     //交易编号 停车场编号+交易发生时间+流水号
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

		_, err := session.Table("b_js_yuansjymx").Insert(Yuansjymx)
		if err != nil {
			log.Println("新增打包明细记录 error", err)
			return err
		}
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增打包明细记录 时，事务错误", serr)
		return serr
	}
	log.Printf("原始交易消息包%d中 明细数据有：%d 条 数据 ", Yuansjymx.FVcXiaoxxh, len(mx))
	log.Println("新增打包明细记录 成功")
	return nil
}

//更新数据    根据 包号 更新原始交易消息包的【发送状态   发送中】
func UpdateYuansjyxx(Mid int64) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	yuansjyxx := new(types.BJsYuansjyxx)
	yuansjyxx.FNbFaszt = 1
	_, err := session.Table("b_js_yuansjyxx").Where("F_NB_XIAOXXH=?", Mid).Update(yuansjyxx)
	if err != nil {
		log.Println(" 根据 包号 更新原始交易消息包的发送状态 为 ： 发送中 error", err)
		return err
	}
	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("根据 包号 更新原始交易消息包的发送状态 为 ： 发送中 时，事务错误", serr)
		return serr
	}
	log.Println(" 根据 包号 更新原始交易消息包的发送状态 为 ： 发送中  成功")
	return nil
}

// 原始交易消息包发送成功更新 发送状态 发送时间 发送成功后消息包的文件路径
func SendedUpdateYuansjyxx(Mid int64, fname string) (error, string) {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	yuansjyxx := new(types.BJsYuansjyxx)
	yuansjyxx.FNbFaszt = 2
	yuansjyxx.FDtFassj = time.Now()
	yuansjyxx.FVcXiaoxwjlj = "compressed_xml/" + fname
	_, err := session.Table("b_js_yuansjyxx").Where("F_NB_XIAOXXH=?", Mid).Update(yuansjyxx)
	if err != nil {
		log.Println(" 根据 包号 更新 发送状态 发送时间 发送成功后消息包的文件路径 error", err)
		return err, ""
	}
	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("根据 包号 更新 发送状态、发送时间、发送成功后的消息包的文件路径 时，事务错误", serr)
	}

	log.Println(" 根据 包号 更新 发送状态 发送时间 发送成功后消息包的文件路径  成功")

	sj := common.DateTimeFormat(yuansjyxx.FDtFassj)
	return nil, sj
}

//   更新结算数据打包结果【打包状态：已打包、原始交易包号、包内序号、清分目标日】
func UpdateDataPackagingResults(Jiaoyjlid []string, Msgid int64, jiaoyisj *types.Message) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	for i, idstr := range Jiaoyjlid {
		Jiessj := new(types.BJsJiessj)
		Jiessj.FNbDabzt = 2                                //已打包
		Jiessj.FNbYuansjybxh = Msgid                       //消息包id
		Jiessj.FNbJiaoybnxh = i + 1                        //包内序号
		Jiessj.FVcQingfmbr = jiaoyisj.Body.ClearTargetDate //清分目标日 打包当天日期

		_, err := session.Table("b_js_jiessj").Where("F_VC_JIAOYJLID=?", idstr).Update(Jiessj)
		if err != nil {
			log.Println("更新打包状态失败", err)
			return err
		}
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("更新 打包状态：已打包、原始交易包号、包内序号、清分目标日 时，事务错误", serr)
	}
	log.Println("更新 打包状态：已打包、原始交易包号、包内序号、清分目标日 成功")
	return nil
}

//   新增应答记录
func PackagingRespRecordInsert(data *types.BJsYuansjyydxx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)

	Yuansjyydxx := new(types.BJsYuansjyydxx)
	//赋值
	Yuansjyydxx.FVcBanbh = data.FVcBanbh         //F_VC_BANBH	版本号	VARCHAR(32)
	Yuansjyydxx.FNbXiaoxlb = data.FNbXiaoxlb     //F_NB_XIAOXLB	消息类别	INT
	Yuansjyydxx.FNbXiaoxlx = data.FNbXiaoxlx     //F_NB_XIAOXLX	消息类型	INT
	Yuansjyydxx.FVcFaszid = data.FVcFaszid       //F_VC_FASZID	发送者ID	VARCHAR(32)
	Yuansjyydxx.FVcJieszid = data.FVcJieszid     //F_VC_JIESZID	接收者ID	VARCHAR(32)
	Yuansjyydxx.FNbXiaoxxh = data.FNbXiaoxxh     //F_NB_XIAOXXH	消息序号	BIGINT
	Yuansjyydxx.FNbQuerdxxxh = data.FNbQuerdxxxh //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	Yuansjyydxx.FDtChulsj = data.FDtChulsj       //F_DT_CHULSJ	处理时间	DATETIME
	Yuansjyydxx.FNbZhixjg = data.FNbZhixjg       //F_NB_ZHIXJG	执行结果	INT
	Yuansjyydxx.FVcQingfmbr = data.FVcQingfmbr   //F_VC_QINGFMBR	清分目标日	VARCHAR(32)
	_, err := session.Table("b_js_yuansjyydxx").Insert(Yuansjyydxx)
	if err != nil {
		log.Println("新增应答记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增应答记录 时，事务错误", serr)
		return serr
	}
	return nil
}
