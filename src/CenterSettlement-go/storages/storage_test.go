package storage

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/service"
	"CenterSettlement-go/types"
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
	s := []string{"320700001111022020060209372900000432", "320700001111022020060209372900000431"}
	UpdatePackaging(s)
}
func TestGetTingcc(t *testing.T) {
	GetTingcc("4242420021")
}

func TestPackagingResRecordInsert(t *testing.T) {
	var yuansjyxx types.BJsYuansjyxx
	yuansjyxx.FVcBanbh = "00010000"                               //版本号
	yuansjyxx.FNbXiaoxlb = 5                                      //消息类别
	yuansjyxx.FNbXiaoxlx = 7                                      //消息类型
	yuansjyxx.FVcFaszid = "00000000000000FD"                      //发送者ID
	yuansjyxx.FVcJieszid = "0000000000000020"                     //接受者ID
	yuansjyxx.FNbXiaoxxh = 12343454                               //消息序号【消息包号】
	yuansjyxx.FDtDabsj = time.Now().Format("2020-01-02 15:04:05") // 打包时间
	yuansjyxx.FVcQingfmbr = "jiaoyisj.Body.ClearTargetDate "      //清分目标日
	yuansjyxx.FVcTingccqffid = "00000000000000FD"                 //停车场清分方ID
	yuansjyxx.FVcFaxfwjgid = "0000000000000020"                   //发行服务机构ID 0000000000000020
	yuansjyxx.FNbJilsl = 2                                        //记录数量
	yuansjyxx.FNbZongje = "212"                                   //总金额
	yuansjyxx.FVcXiaoxwjlj = "generatexml/+ fname  "              //消息文件路径
	err := PackagingRecordInsert(yuansjyxx)
	if err != nil {
		log.Println("PackagingRecordInsert error")
	}
}

func TestPackagingMXRecordInsert(t *testing.T) {

	mx := service.YuanshiMsgMXAssignment(jiaoyisj)
	PackagingMXRecordInsert(mx)
}
