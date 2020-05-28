package storage

import (
	"CenterSettlement-go/database"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

//测试查询本省的结算数据 记账卡
func TestQueryJiessjjz(t *testing.T) {
	QueryJiessjjz()
}

//测试查询本省的结算数据 储值卡
func TestQueryJiessjcz(t *testing.T) {
	QueryJiessjcz()
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
	database.DBInit()
	var sRecord BJsTcpqqjl
	var err error
	sRecord.F_VC_XIAOXXH = "111111"
	sRecord.F_NB_FASZ = 1 //结算平台
	sRecord.F_DT_ZUIXSJ = time.Now().Format("2006-01-02 15:04:05")
	sRecord.F_NB_MD5 = "jgjkg42314j5jh231v5321551"
	sRecord.F_NB_XIAOXCD = 650 //消息长度
	sRecord.F_NB_CHONGFCS = 0
	log.Println(sRecord)
	err = TcpSendRecordInsert(sRecord)
	if err != nil {
		log.Println("测试tcp发送记录 error")
	}
}

//测试tcp接收即使应答记录
func TestTcpResponseRecordInsert(t *testing.T) {
	var resRecord BJsTcpydjl
	//赋值
	resRecord.F_NB_CHONGFCS = 0
	resRecord.F_DT_ZUIXSJ = time.Now().Format("2006-01-02 15:04:05")
	resRecord.F_VC_XIAOXXH = "999999"
	resRecord.F_NB_FASZ = 2
	err := TcpResponseRecordInsert(resRecord)
	if err != nil {
		log.Println("测试tcp发送记录 error")
	}
	log.Println(resRecord)
}

//
func TestUpdatePackaging(t *testing.T) {
	s := []string{"051700000111022020042615313200000173", "320102000111032020042414232900000585"}
	UpdatePackaging(s)
}
