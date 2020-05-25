package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"log"
)

//清分统计

//新增清分包记录
func ClearingInsert(data types.BJsQingftjxx) error {

	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	qingftongjixx := new(types.BJsQingftjxx)

	//赋值
	qingftongjixx.FDtJiessj = data.FDtJiessj

	//插入
	_, err := xorm.Insert(qingftongjixx)
	if err != nil {
		log.Fatal("新增记账处理的消息记录 error")
		return err
	}
	return nil

}

//新增清分统计消息明细包记录
func ClearingMXInsert(data types.BJsQingftongjimx) error {

	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	qingftongjimx := new(types.BJsQingftongjimx)

	//赋值
	qingftongjimx.FNbFenzxh = data.FNbFenzxh

	//插入
	_, err := xorm.Insert(qingftongjimx)
	if err != nil {
		log.Fatal("新增记账处理的消息记录 error")
		return err
	}
	return nil

}

//新增清分应答消息包记录

//更新结算消息记录【清分执行结果，执行时间】

//更新结算数据
