package service

import (
	"CenterSettlement-go/storage"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

//测试查询本省的结算数据 记账卡
func TestQueryJiessjjz(t *testing.T) {
	//QueryJiessjjz()
}

//测试查询本省的结算数据 储值卡
func TestQueryJiessjcz(t *testing.T) {
	//QueryJiessjcz()
}

//测试查询本省的结算数据
func TestQueryJiessj(t *testing.T) {

}

//测试查询其他地区的结算数据
func TestQueryQiTaJiessj(t *testing.T) {
	//QueryQiTaJiessj()
}

//测试tcp发送记录
func TestTcpSendRecordInsert(t *testing.T) {
	//database.DBInit()
	var sRecord storage.BJsTcpqqjl
	var err error
	//FVcXiaoxxh  string //`F_VC_XIAOXXH` '消息序号',
	//	FNbFasz     int    //`F_NB_FASZ`   DEFAULT '1' COMMENT '发送者   1、ETC结算开放平台，2、联网中心',
	//	FNbChongfcs int    //`F_NB_CHONGFCS`   '重复次数',
	//	FDtZuixsj   string //`F_DT_ZUIXSJ`   '最新时间',
	//	FNbXIiaoxcd int    //`F_NB_XIAOXCD`  '消息长度',
	//	FNbMd5      string //`F_NB_MD5`   MD5',
	sRecord.FVcXiaoxxh = "111111"
	sRecord.FNbFasz = 1 //结算平台
	sRecord.FDtZuixsj = time.Now().Format("2006-01-02 15:04:05")
	sRecord.FNbMd5 = "jgjkg42314j5jh231v5321551"
	sRecord.FNbXiaoxcd = 650 //消息长度
	sRecord.FNbChongfcs = 0
	err = storage.TcpSendRecordInsert(sRecord)
	if err != nil {
		log.Println("测试tcp发送记录 error")
	}
}

//测试tcp接收即使应答记录
func TestTcpResponseRecordInsert(t *testing.T) {
	var resRecord storage.BJsTcpydjl
	//赋值
	resRecord.FNbChongfcs = 0
	resRecord.FDtZuixsj = time.Now().Format("2006-01-02 15:04:05")
	resRecord.FVcXiaoxxh = "999999"
	resRecord.FNbFasz = 2
	err := storage.TcpResponseRecordInsert(resRecord)
	if err != nil {
		log.Println("测试tcp发送记录 error")
	}
	log.Println(resRecord)
}

//
func TestUpdatePackaging(t *testing.T) {
	s := []string{"320700001111022020060209372900000432", "320700001111022020060209372900000431"}
	storage.UpdatePackaging(s)
}
func TestGetTingcc(t *testing.T) {
	storage.GetTingcc("4242420021")
}
