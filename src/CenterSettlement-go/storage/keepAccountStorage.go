package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"log"
)

//记账处理数据层

//1、新增记账处理的消息记录
func InsertMessageData(data *types.BJsJizclxx) error {
	database.DBInit()
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	jizclxx := new(types.BJsJizclxx)

	//赋值
	jizclxx.FVcBanbh = data.FVcBanbh             //F_VC_BANBH	版本号	VARCHAR(32)
	jizclxx.FNbXiaoxlb = data.FNbXiaoxlb         //F_NB_XIAOXLB	消息类别	INT
	jizclxx.FNbXiaoxlx = data.FNbXiaoxlx         //F_NB_XIAOXLX	消息类型	INT
	jizclxx.FVcFaszid = data.FVcFaszid           //F_VC_FASZID	发送者ID	VARCHAR(32)
	jizclxx.FVcJieszid = data.FVcJieszid         //F_VC_JIESZID	接收者ID	VARCHAR(32)
	jizclxx.FNbXiaoxxh = data.FNbXiaoxxh         //F_NB_XIAOXXH	消息序号	BIGINT
	jizclxx.FDtJiessj = data.FDtJiessj           //F_DT_JIESSJ	接收时间	DATETIME
	jizclxx.FNbYuansjyxxxh = data.FNbYuansjyxxxh //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	jizclxx.FNbJilsl = data.FNbJilsl             //F_NB_JILSL	记录数量	INT
	jizclxx.FNbZongje = data.FNbZongje           //F_NB_ZONGJE	总金额	INT
	jizclxx.FNbZhengysl = data.FNbZhengysl       //F_NB_ZHENGYSL	争议数量	INT
	jizclxx.FNbZhixjg = data.FNbZhixjg           //F_NB_ZHIXJG	执行结果	INT
	jizclxx.FDtChulsj = data.FDtChulsj           //F_DT_CHULSJ	处理时间	DATETIME
	jizclxx.FVcXiaoxwjlj = data.FVcXiaoxwjlj     //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
	log.Println("新增记账处理的消息记录 内容", *jizclxx)
	//插入
	_, inserterr := xorm.Insert(jizclxx)
	if inserterr != nil {
		log.Println("新增记账处理的消息记录 error", inserterr)
		return inserterr
	}
	log.Println("新增记账处理的消息记录 成功")

	return nil

}

//2、新增记账处理消息明细记录
func InsertMessageMXData(data *types.BJsJizclmx) error {

	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	jizclmx := new(types.BJsJizclmx)

	//赋值
	jizclmx.FNbYuansjyxxxh = data.FNbYuansjyxxxh //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	jizclmx.FNbChuljg = data.FNbChuljg           //处理结果
	jizclmx.FNbbaonxh = data.FNbbaonxh           //包内序号

	//插入
	_, err := xorm.Insert(jizclmx)
	if err != nil {
		log.Fatal("新增记账处理消息明细记录 error")
		return err
	}
	return nil
}

//3、新增记账处理应答消息记录
func InsertResMessageData(data types.BJsJizclydxx) error {

	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	jizclydxx := new(types.BJsJizclydxx)

	//赋值
	jizclydxx.FNbQuerdxxxh = data.FNbQuerdxxxh

	//插入
	_, err := xorm.Insert(jizclydxx)
	if err != nil {
		log.Fatal("新增记账处理应答消息记录 error")
		return err
	}
	return nil
}

//4、更新记账处理消息记录  执行结果、处理时间
func update() {

}

//5、查询原始交易明细表 查询消息序号 内 的所有数据 获取包内序号
func QueryYuanshiMx(msgid int64) *[]types.BJsYuansjymx {
	database.DBInit()
	xorm := database.XormClient
	//查询多条数据
	tests := make([]types.BJsYuansjymx, 0)
	//yuanshimx := new(types. BJsYuansjymx)
	//赋值

	qerr := xorm.Where("F_VC_XIAOXXH=?", msgid).Find(&tests)
	if qerr != nil {
		log.Fatalln("查询结算数据出错", qerr)
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("消息序号为%d的原始交易包，包内序号 %d \n", v.FVcXiaoxxh, v.FNbBaonxh)
	}
	return &tests
}

//6、 更新结算数据  记账结果：已记账(1)
func KeepAccountUpdate(Msgid int64, bnxh int, jzjg int) error {
	database.DBInit()
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	Jiessj := new(types.BJsJiessj)
	Jiessj.FNbJizjg = jzjg
	_, err := session.Table("b_js_jiessj").Where("F_NB_YUANSJYBXH=?", Msgid).Where("F_NB_JIAOYBNXH=?", bnxh).Update(Jiessj)
	if err != nil {
		log.Println("更新结算数据失败", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("更新结算数据  记账结果 时，事务错误")
	}
	log.Printf("更新结算数据  记账结果 为%d 成功", jzjg)
	return nil
}
