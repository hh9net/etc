package storage

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/database"
	"CenterSettlement-go/types"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//清分统计

//新增清分包记录
func ClearingInsert(data *types.BJsQingftjxx) error {
	database.DBInit()
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	qingftongjixx := new(types.BJsQingftjxx)

	//赋值
	qingftongjixx.FVcBanbh = data.FVcBanbh                 //F_VC_BANBH	版本号	VARCHAR(32)
	qingftongjixx.FNbXiaoxlb = data.FNbXiaoxlb             //F_NB_XIAOXLB	消息类别	INT
	qingftongjixx.FNbXiaoxlx = data.FNbXiaoxlx             //F_NB_XIAOXLX	消息类型	INT
	qingftongjixx.FVcFaszid = data.FVcFaszid               //F_VC_FASZID	发送者ID	VARCHAR(32)
	qingftongjixx.FVcJieszid = data.FVcJieszid             //F_VC_JIESZID	接收者ID	VARCHAR(32)
	qingftongjixx.FNbXiaoxxh = data.FNbXiaoxxh             //F_NB_XIAOXXH	消息序号	BIGINT
	qingftongjixx.FDtJiessj = data.FDtJiessj               //F_DT_JIESSJ	接收时间	DATETIME
	qingftongjixx.FVcQingfmbr = data.FVcQingfmbr           //F_VC_QINGFMBR	清分目标日	DATE
	qingftongjixx.FNbQingfzje = data.FNbQingfzje           //F_NB_QINGFZJE	清分总金额	INT
	qingftongjixx.FNbQingfsl = data.FNbQingfsl             //F_NB_QINGFSL	清分数量	INT
	qingftongjixx.FDtQingftjclsj = data.FDtQingftjclsj     //F_DT_QINGFTJCLSJ	清分统计处理时间	DATETIME
	qingftongjixx.FNbYuansjysl = data.FNbYuansjysl         //F_NB_YUANSJYSL	原始包交易数量	INT
	qingftongjixx.FNbZhengycljgbsl = data.FNbZhengycljgbsl //F_NB_ZHENGYCLJGBSL	争议处理结果包数量	INT
	qingftongjixx.FVcXiaoxwjlj = data.FVcXiaoxwjlj         //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
	qingftongjixx.FDtChulsj = data.FDtChulsj               //`F_DT_CHULSJ` datetime DEFAULT NULL COMMENT '处理时间',

	//插入
	_, err := session.Insert(qingftongjixx)
	if err != nil {
		log.Println("新增清分处理包的消息记录 error", err)
		return err
	}
	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增清分包记录 时，事务错误", serr)

	}
	log.Println("新增清分处理包的消息记录 成功", qingftongjixx)
	return nil

}

//新增清分统计消息明细包记录
func ClearingMXInsert(data *types.BJsQingftjmx) error {
	database.DBInit()
	xorm := database.XormClient
	session := TransactionBegin(xorm)
	qingftongjimx := new(types.BJsQingftjmx)

	//赋值
	qingftongjimx.FNbQingftjxxxh = data.FNbQingftjxxxh
	qingftongjimx.FVcTongxbzxxtid = data.FVcTongxbzxxtid     //通行宝中心系统ID
	qingftongjimx.FNbFenzxh = data.FNbFenzxh                 //分组序号 入库者自行生成，可取数组下标
	qingftongjimx.FNbYuansjyxxxh = data.FNbYuansjyxxxh       //原始交易消息序号
	qingftongjimx.FNbZhengycljgwjid = data.FNbZhengycljgwjid //争议处理结果文件ID

	//插入
	_, err := session.Insert(qingftongjimx)
	if err != nil {
		log.Println("新增清分统计消息明细包记录 成功")
		return err
	}
	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增清分统计消息明细包记录 时，事务错误", serr)

	}
	log.Println("新增清分统计消息明细包记录 成功", qingftongjimx)
	return nil

}

//新增清分应答消息包记录
func ClearingYDInsert(respmsg *types.ResponseCTMessage) error {

	xorm := database.XormClient
	session := TransactionBegin(xorm)
	qingftongjixxyd := new(types.BJsQingftjxxyd)

	//赋值
	qingftongjixxyd.FVcBanbh = respmsg.Header.Version        //F_VC_BANBH	版本号	VARCHAR(32)
	qingftongjixxyd.FNbXiaoxlb = respmsg.Header.MessageClass //F_NB_XIAOXLB	消息类别	INT
	qingftongjixxyd.FNbXiaoxlx = respmsg.Header.MessageType  //F_NB_XIAOXLX	消息类型	INT
	qingftongjixxyd.FVcFaszid = respmsg.Header.SenderId      //F_VC_FASZID	发送者ID	VARCHAR(32)
	qingftongjixxyd.FVcJieszid = respmsg.Header.ReceiverId   //F_VC_JIESZID	接收者ID	VARCHAR(32)
	qingftongjixxyd.FNbXiaoxxh = respmsg.Header.MessageId    //F_NB_XIAOXXH	消息序号	BIGINT
	qingftongjixxyd.FNbQuerdxxxh = respmsg.Body.MessageId    //F_NB_QUERDXXXH	确认的消息序号	BIGINT

	qingftongjixxyd.FDtChulsj = common.StrTimeTotime(common.DataTimeFormatHandle(respmsg.Body.ProcessTime)) //F_DT_CHULSJ	处理时间	DATETIME

	qingftongjixxyd.FNbZhixjg = respmsg.Body.Result //F_NB_ZHIXJG	执行结果	INT

	qingftongjixxyd.FVcXiaoxwjlj = "generatexml/" + "QF_YDB_" + fmt.Sprintf("%020d", respmsg.Header.MessageId) + ".xml" //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)

	//插入
	_, err := session.Insert(qingftongjixxyd)
	if err != nil {
		log.Fatal("新增清分应答消息包记录 error", err)
		return err
	}

	serr := TransactionCommit(session)
	if serr != nil {
		log.Println("新增清分应答消息包记录时，事务错误", serr)

	}
	return nil
}

//更新结算消息记录【清分执行结果，执行时间】

//更新结算数据
