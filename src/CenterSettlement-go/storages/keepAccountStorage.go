package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"log"
)

//记账处理数据层

//1、新增记账处理的消息记录
func InsertMessageData(data types.BJsJizclxx) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	jizclxx := new(types.BJsJizclxx)

	//赋值
	jizclxx.FDtChulsj = data.FDtChulsj

	//插入
	_, err := xorm.Insert(jizclxx)
	if err != nil {
		log.Fatal("新增记账处理的消息记录 error")
		return err
	}
	return nil

}

//2、新增记账处理消息明细记录
func InsertMessageMXData(data types.BJsJizclmx) error {

	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	jizclmx := new(types.BJsJizclmx)

	//赋值
	jizclmx.FNbbaonxh = data.FNbbaonxh

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
