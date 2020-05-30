package storage

import (
	"CenterSettlement-go/database"
	log "github.com/sirupsen/logrus"
)

//tcp发送与接收包记录

//表字段
//请求记录表
type BJsTcpqqjl struct {
	F_VC_XIAOXXH  string //`F_VC_XIAOXXH` '消息序号',
	F_NB_FASZ     int    //`F_NB_FASZ`   DEFAULT '1' COMMENT '发送者   1、ETC结算开放平台，2、联网中心',
	F_NB_CHONGFCS int    //`F_NB_CHONGFCS`   '重复次数',
	F_DT_ZUIXSJ   string //`F_DT_ZUIXSJ`   '最新时间',
	F_NB_XIAOXCD  int    //`F_NB_XIAOXCD`  '消息长度',
	F_NB_MD5      string //`F_NB_MD5`   MD5',
}
type BJsTcpydjl struct {
	F_VC_XIAOXXH  string //  `F_VC_XIAOXXH`   '消息序号',
	F_NB_FASZ     int    //  `F_NB_FASZ` int(11) DEFAULT '1' COMMENT '发送者 1、ETC结算开放平台，2、联网中心',
	F_NB_CHONGFCS int    //  `F_NB_CHONGFCS` int(11) DEFAULT '0' COMMENT '重复次数',
	F_DT_ZUIXSJ   string //  `F_DT_ZUIXSJ` '最新时间',
}

//表操作
//1、 发送tcp包，记录tcp数据包

func TcpSendRecordInsert(record BJsTcpqqjl) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	tcprecord := new(BJsTcpqqjl)
	//赋值
	tcprecord.F_DT_ZUIXSJ = record.F_DT_ZUIXSJ
	tcprecord.F_NB_CHONGFCS = record.F_NB_CHONGFCS
	tcprecord.F_NB_FASZ = record.F_NB_FASZ
	tcprecord.F_NB_MD5 = record.F_NB_MD5
	tcprecord.F_NB_XIAOXCD = record.F_NB_XIAOXCD
	tcprecord.F_VC_XIAOXXH = record.F_VC_XIAOXXH
	//
	log.Println("tcpsend记录")
	//插入
	_, err := xorm.Insert(tcprecord)
	if err != nil {
		log.Fatal("发送tcp包，记录tcp数据包 error")
		return err
	}
	return nil
}

//2、即使应答包记录存储
func TcpResponseRecordInsert(resRecord BJsTcpydjl) error {
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	tcpResRecord := new(BJsTcpydjl)

	//赋值
	tcpResRecord.F_VC_XIAOXXH = resRecord.F_VC_XIAOXXH
	tcpResRecord.F_NB_FASZ = resRecord.F_NB_FASZ
	tcpResRecord.F_NB_CHONGFCS = resRecord.F_NB_CHONGFCS
	tcpResRecord.F_DT_ZUIXSJ = resRecord.F_DT_ZUIXSJ

	//插入
	_, err := xorm.Insert(tcpResRecord)
	if err != nil {
		log.Fatal("即使应答包记录存储error")
		return err
	}
	return nil
}
