package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"log"
)

//争议处理
//新增争议处理结果消息包记录
func DisputeProcessXXInsert(data *types.BJsZhengyjyclxx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	zhengyclxx := new(types.BJsZhengyjyclxx)

	//赋值
	zhengyclxx.FVcBanbh = data.FVcBanbh               //F_VC_BANBH	版本号	VARCHAR(32)
	zhengyclxx.FNbXiaoxlb = data.FNbXiaoxlb           //F_NB_XIAOXLB	消息类别	INT
	zhengyclxx.FNbXiaoxlx = data.FNbXiaoxlx           //F_NB_XIAOXLX	消息类型	INT
	zhengyclxx.FVcFaszid = data.FVcFaszid             //F_VC_FASZID	发送者ID	VARCHAR(32)
	zhengyclxx.FVcJieszid = data.FVcJieszid           //F_VC_JIESZID	接收者ID	VARCHAR(32)
	zhengyclxx.FNbXiaoxxh = data.FNbXiaoxxh           //F_NB_XIAOXXH	消息序号	BIGINT
	zhengyclxx.FDtJiessj = data.FDtJiessj             //F_DT_JIESSJ	接收时间	DATETIME
	zhengyclxx.FVcQingffid = data.FVcQingffid         //F_VC_QINGFFID	清分方ID	VARCHAR(32)
	zhengyclxx.FVcLianwzxid = data.FVcLianwzxid       //F_VC_LIANWZXID	联网中心ID	VARCHAR(32)
	zhengyclxx.FVcFaxfid = data.FVcFaxfid             //F_VC_FAXFID	发行方ID	VARCHAR(32)
	zhengyclxx.FVcZhengyjgwjid = data.FVcZhengyjgwjid //F_VC_ZHENGYJGWJID	争议结果文件ID	INT
	zhengyclxx.FDtZhengyclsj = data.FDtZhengyclsj     //F_DT_ZHENGYCLSJ	争议处理时间	DATETIME
	zhengyclxx.FNbZhengysl = data.FNbZhengysl         //F_NB_ZHENGYSL	争议数量	INT
	zhengyclxx.FNbQuerxyjzdzje = data.FNbQuerxyjzdzje //F_NB_QUERXYJZDZJE	确认需要记账的总金额	INT
	zhengyclxx.FVcXiaoxwjlj = data.FVcXiaoxwjlj       //消息路径

	//插入
	_, err := session.Insert(zhengyclxx)
	if err != nil {
		log.Fatal("新增争议处理包的消息记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增争议处理结果消息包记录 时，事务错误", serr)

	}

	return nil
}

//新增争议处理结果消息包明细记录
func DisputeProcessMxInsert(data *types.BJsZhengyjyclmx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	zhengyjyclmx := new(types.BJsZhengyjyclmx)

	//赋值
	zhengyjyclmx.FNbZhengyjyclxxxh = data.FNbZhengyjyclxxxh //F_NB_ZHENGYJYCLXXXH	争议交易处理消息序号	BIGINT
	zhengyjyclmx.FNbFenzxh = data.FNbFenzxh                 //F_NB_FENZXH	分组序号	INT
	zhengyjyclmx.FNbYuansjyxxxh = data.FNbYuansjyxxxh       //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	zhengyjyclmx.FNbZunjlsl = data.FNbZunjlsl               //F_NB_ZUNJLSL	组内记录数量	INT
	zhengyjyclmx.FNbZunjezh = data.FNbZunjezh               //F_NB_ZUNJEZH	组内金额总和	INT
	zhengyjyclmx.FNbYuansbnxh = data.FNbYuansbnxh           //F_NB_YUANSBNXH	原始包内序号	INT
	zhengyjyclmx.FNbChuljg = data.FNbChuljg                 //F_NB_CHULJG	处理结果	INT

	//插入
	_, err := session.Insert(zhengyjyclmx)
	if err != nil {
		log.Fatal("新增争议处理结果消息包明细记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增争议处理结果消息包明细记录 时，事务错误", serr)

	}
	return nil
}

//新增争议处理结果应答消息包记录
func DisputeProcessResInsert(data types.BJsZhengyclydxx) error {
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	zhengyclydxx := new(types.BJsZhengyclydxx)

	//赋值
	zhengyclydxx.FVcBanbh = data.FVcBanbh         //F_VC_BANBH	版本号	VARCHAR(32)
	zhengyclydxx.FNbXiaoxlb = data.FNbXiaoxlb     //F_NB_XIAOXLB	消息类别	INT
	zhengyclydxx.FNbXiaoxlx = data.FNbXiaoxlx     //F_NB_XIAOXLX	消息类型	INT
	zhengyclydxx.FVcFaszid = data.FVcFaszid       //F_VC_FASZID	发送者ID	VARCHAR(32)
	zhengyclydxx.FVcJieszid = data.FVcJieszid     //F_VC_JIESZID	接收者ID	VARCHAR(32)
	zhengyclydxx.FNbXiaoxxh = data.FNbXiaoxxh     //F_NB_XIAOXXH	消息序号	BIGINT
	zhengyclydxx.FNbQuerdxxxh = data.FNbQuerdxxxh //F_NB_QUERDXXXH	确认的消息序号	BIGINT
	zhengyclydxx.FVcChulsj = data.FVcChulsj       //F_DT_CHULSJ	处理时间	DATETIME
	zhengyclydxx.FNbZhixjg = data.FNbZhixjg       //F_NB_ZHIXJG	执行结果	INT
	zhengyclydxx.FVcXiaoxwjlj = data.FVcXiaoxwjlj //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)

	//插入
	_, err := session.Insert(zhengyclydxx)
	if err != nil {
		log.Fatal("新增争议处理结果应答消息包记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增争议处理结果应答消息包记录时，事务错误", serr)

	}
	return nil
}

//更新结算数据【争议处理结果】争议处理包序号】

//更新争议处理结果消息记录【执行结果、处理时间】
