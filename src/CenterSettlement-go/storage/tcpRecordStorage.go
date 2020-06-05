package storage

import (
	"CenterSettlement-go/database"
	log "github.com/sirupsen/logrus"
)

//tcp发送与接收包记录

//表字段
//请求记录表
type BJsTcpqqjl struct {
	FVcXiaoxxh  string //`F_VC_XIAOXXH` '消息序号',
	FNbFasz     int    //`F_NB_FASZ`   DEFAULT '1' COMMENT '发送者   1、ETC结算开放平台，2、联网中心',
	FNbChongfcs int    //`F_NB_CHONGFCS`   '重复次数',
	FDtZuixsj   string //`F_DT_ZUIXSJ`   '最新时间',
	FNbXiaoxcd  int    //`F_NB_XIAOXCD`  '消息长度',
	FNbMd5      string //`F_NB_MD5`   MD5',
}
type BJsTcpydjl struct {
	FVcXiaoxxh  string //  `F_VC_XIAOXXH`   '消息序号',
	FNbFasz     int    //  `F_NB_FASZ` int(11) DEFAULT '1' COMMENT '发送者 1、ETC结算开放平台，2、联网中心',
	FNbChongfcs int    //  `F_NB_CHONGFCS` int(11) DEFAULT '0' COMMENT '重复次数',
	FDtZuixsj   string //  `F_DT_ZUIXSJ` '最新时间',
}

//表操作
//1、 发送tcp包，记录tcp数据包
func TcpSendRecordInsert(record BJsTcpqqjl) error {
	//database.DBInit()
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	tcprecord := new(BJsTcpqqjl)
	//赋值
	tcprecord.FDtZuixsj = record.FDtZuixsj
	tcprecord.FVcXiaoxxh = record.FVcXiaoxxh
	tcprecord.FNbFasz = record.FNbFasz
	tcprecord.FNbMd5 = record.FNbMd5
	tcprecord.FNbXiaoxcd = record.FNbXiaoxcd
	tcprecord.FNbChongfcs = record.FNbChongfcs
	//插入
	_, err := xorm.Insert(tcprecord)
	if err != nil {
		log.Fatal("发送tcp包，记录tcp数据包 error", err)
		return err
	}
	log.Println("tcpsend记录 成功")
	return nil
}

func TcpSendRecordUpdate(record BJsTcpqqjl) error {
	//database.DBInit()
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	tcprecord := new(BJsTcpqqjl)
	//赋值

	//tcprecord.FDtZuixsj = record.FDtZuixsj
	//tcprecord.FVcXiaoxxh = record.FVcXiaoxxh
	//tcprecord.FNbFasz = record.FNbFasz
	//	tcprecord.FNbMd5 = record.FNbMd5
	//tcprecord.FNbXiaoxcd = record.FNbXiaoxcd
	tcprecord.FNbChongfcs = record.FNbChongfcs
	//插入
	log.Println(record)
	_, err := xorm.Where("F_VC_XIAOXXH=?", record.FVcXiaoxxh).Update(tcprecord)
	if err != nil {
		log.Fatal("发送tcp包，Update tcp数据包 error", err)
		return err
	}
	log.Println("tcpsend记录 Update 成功")
	return nil
}

func GetTcpSendRecord(msgid string) (bool, error, int) {
	//database.DBInit()
	xorm := database.XormClient
	tcprecord := &BJsTcpqqjl{FVcXiaoxxh: msgid}
	has, err := xorm.Get(tcprecord)
	log.Println("获取数据：", has, err, tcprecord.FNbChongfcs)
	return has, err, tcprecord.FNbChongfcs
}

//2、即使应答包记录存储
func TcpResponseRecordInsert(resRecord BJsTcpydjl) error {
	//database.DBInit()
	xorm := database.XormClient
	//session := TransactionBegin(xorm)
	tcpResRecord := new(BJsTcpydjl)

	//赋值
	tcpResRecord.FVcXiaoxxh = resRecord.FVcXiaoxxh
	tcpResRecord.FNbFasz = resRecord.FNbFasz
	tcpResRecord.FNbChongfcs = resRecord.FNbChongfcs
	tcpResRecord.FDtZuixsj = resRecord.FDtZuixsj

	//插入
	_, err := xorm.Insert(tcpResRecord)
	if err != nil {
		log.Fatal("即使应答包记录存储error", err)
		return err
	}
	log.Println("tcpsend记录 即使应答包记录存储 成功")

	return nil
}
