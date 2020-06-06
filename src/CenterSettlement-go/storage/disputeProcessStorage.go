package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"log"
)

//争议处理
//新增争议处理结果消息包记录
func DisputeProcessXXInsert(data *types.BJsZhengyclxx) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	zhengyclxx := new(types.BJsZhengyclxx)

	//赋值
	zhengyclxx.FDtJiessj = data.FDtJiessj

	//插入
	_, err := xorm.Insert(zhengyclxx)
	if err != nil {
		log.Fatal("新增争议处理包的消息记录 error")
		return err
	}
	return nil
}

//新增争议处理结果消息包明细记录
func DisputeProcessMxInsert(data *types.BJsZhengyjyclmx) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	zhengyjyclmx := new(types.BJsZhengyjyclmx)

	//赋值
	zhengyjyclmx.FNbFenzxh = data.FNbFenzxh

	//插入
	_, err := xorm.Insert(zhengyjyclmx)
	if err != nil {
		log.Fatal("新增争议处理结果消息包明细记录 error")
		return err
	}
	return nil
}

//新增争议处理结果应答消息包记录
func DisputeProcessResInsert(data types.BJsZhengyclydxx) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	zhengyclydxx := new(types.BJsZhengyclydxx)

	//赋值
	zhengyclydxx.FNbQuerdxxxh = data.FNbQuerdxxxh

	//插入
	_, err := xorm.Insert(zhengyclydxx)
	if err != nil {
		log.Fatal("新增争议处理结果应答消息包记录 error")
		return err
	}
	return nil
}

//更新结算数据【争议处理结果】争议处理包序号】

//更新争议处理结果消息记录【执行结果、处理时间】
