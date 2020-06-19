package centerServer

import (
	"github.com/go-xorm/xorm"
	"log"
)

const pwd = "shx19930509321"

type DB struct {
	orm *xorm.Engine
}

func NewDatabase() *DB {
	var db DB
	xo, err := xorm.NewEngine("mysql", "root:"+pwd+"@tcp(127.0.0.1:3306)/center?charset=utf8")
	if err != nil {
		log.Println("连接联网中心数据库失败")
		return nil
	}
	db.orm = xo
	if db.orm != nil {
		log.Println("连接联网中心数据库成功")
	}
	return &db
}

//创建数据库
func (db *DB) NewTable() {
	db = NewDatabase()
	is, err := db.orm.IsTableExist(
		//new(JieSuanMessage),
		new(JieSuanMessage),
		//new(JieSuanMessageMx),
	)
	if err != nil {
		log.Println("创建数据库判断数据库表是否存在时 error  ", err)
	}
	if is == false {
		err := db.orm.Sync2(
			new(JieSuanMessage),
			//new(JieSuanMessageMx),
			//new(Jiessjchuli),
		)
		//err = db.orm.CreateTables(new(SJsJiessj))
		if err != nil {
			log.Println("创建数据库 映射表 error ", err)
		}
	} else {
		log.Println("创建数据库 映射表 表已存在")
	}
}

//新增结算数据消息包
func JieSuanMessageInset(msg Message) error {
	db := NewDatabase()
	jiessjmsg := new(JieSuanMessage)
	//赋值

	jiessjmsg.Version = msg.Header.Version
	jiessjmsg.MessageClass = msg.Header.MessageClass
	jiessjmsg.MessageType = msg.Header.MessageType
	jiessjmsg.SenderId = msg.Header.SenderId
	jiessjmsg.ReceiverId = msg.Header.ReceiverId
	jiessjmsg.MessageId = msg.Header.MessageId
	jiessjmsg.ClearTargetDate = msg.Body.ClearTargetDate
	jiessjmsg.ServiceProviderId = msg.Body.ServiceProviderId
	jiessjmsg.IssuerId = msg.Body.IssuerId
	jiessjmsg.Count = msg.Body.Count
	jiessjmsg.Amount = msg.Body.Amount

	_, err := db.orm.Insert(jiessjmsg)
	if err != nil {
		log.Println("联网中心 新增结算数据 时 错误", err)
		return err
	}
	log.Println("联网中心 新增结算数据 时 成功")
	return nil
}

//新增结算数据明细
func JieSuanMessageMxInset(msg Message) error {
	db := NewDatabase()
	//jiessjmsgs := make([]JieSuanMessageMx, 0)

	jiessjmsg := new(JieSuanMessageMx)
	//赋值
	jiessjmsg.Version = msg.Header.Version
	jiessjmsg.MessageClass = msg.Header.MessageClass
	jiessjmsg.MessageType = msg.Header.MessageType
	jiessjmsg.SenderId = msg.Header.SenderId
	jiessjmsg.ReceiverId = msg.Header.ReceiverId
	jiessjmsg.MessageId = msg.Header.MessageId
	jiessjmsg.ClearTargetDate = msg.Body.ClearTargetDate
	jiessjmsg.ServiceProviderId = msg.Body.ServiceProviderId
	jiessjmsg.IssuerId = msg.Body.IssuerId

	for _, T := range msg.Body.Transaction {
		jiessjmsg.TransId = T.TransId
		jiessjmsg.Time = T.Time
		jiessjmsg.Fee = T.Fee

		jiessjmsg.CustomizedData = T.CustomizedData           //特定发行方与通行宝收费方之间 约定格式的交易信息【  】
		jiessjmsg.BhId = T.Id                                 //停车场消费交易编号(停车场编号+交易发生的时间+流水号 )
		jiessjmsg.Name = T.Name                               //停车场名称(不超过150个字符)
		jiessjmsg.ParkTime = T.ParkTime                       //停放时长(单位：分)
		jiessjmsg.VehicleType = T.VehicleType                 //收费车型
		jiessjmsg.AlgorithmIdentifier = T.AlgorithmIdentifier //算法标识 1-3DEX  2-SM4

		jiessjmsg.ServiceType = T.Service.ServiceType //交易的服务类型【【2 写死】】
		jiessjmsg.Description = T.Service.Description //对交易的文字解释【停车场名｜停车时常 ：几时几分几秒】
		jiessjmsg.Detail = T.Service.Detail           //交易详细信息  1|04|3201|3201020001|1104|20200513 143434|03|3201|3201020001|0001|20200513 140805

		jiessjmsg.CardType = T.ICCard.CardType       //卡类型，22为储值卡；23记账卡
		jiessjmsg.NetNo = T.ICCard.NetNo             //网络编码，BCD码 Hex(4) ka网络号（16进制） 数据库10进制
		jiessjmsg.CardId = T.ICCard.CardId           //IC卡物理编号，BCD码  Hex(16)   卡号
		jiessjmsg.License = T.ICCard.License         //0015文件中记录的车牌号 卡内车牌号 【含颜色】
		jiessjmsg.PreBalance = T.ICCard.PreBalance   //交易前余额，以元为单位 Decimal
		jiessjmsg.PostBalance = T.ICCard.PostBalance //交易后余额，以元为单位 Decimal

		jiessjmsg.TAC = T.Validation.TAC                         //交易时产生的TAC码，8位16进制数   Hex(8)
		jiessjmsg.TransType = T.Validation.TransType             //交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交易  Hex(2)【09】
		jiessjmsg.TerminalNo = T.Validation.TerminalNo           //12位16进制数据，即PSAM号，PSAM中0016文件中的终端机编号  Hex(2) 	加密卡号
		jiessjmsg.TerminalTransNo = T.Validation.TerminalTransNo //8位16进制数，PSAM卡脱机交易序号，在MAC1计算过程中得到  Hex(8) 加密序列号

		jiessjmsg.OBUId = T.OBU.OBUId       //OBU物理编号，BCD码  4501191509252866
		jiessjmsg.OBEState = T.OBU.OBEState //2字节的OBU状态

		//jiessjmsgs = append(jiessjmsgs, *jiessjmsg)

		_, err := db.orm.Insert(jiessjmsg)
		if err != nil {
			log.Println("联网中心 新增结算数据 时 错误", err)
			return err
		}
	}

	log.Println("联网中心 新增结算数据明细 时 成功")
	return nil
}

//新增结算处理
func JieSuanMessageChuliInset(msg Message) error {
	db := NewDatabase()
	jiessjmsg := new(Jiessjchuli)
	//赋值

	//jiessjmsg.FVcJiaoyjlid   string   //F_VC_JIAOYJLID	交易记录ID	VARCHAR(128)

	jiessjmsg.FNbYuansjybxh = msg.Header.MessageId //F_NB_YUANSJYBXH	原始交易包序号	BIGINT
	for _, T := range msg.Body.Transaction {

		jiessjmsg.FNbJiaoybnxh = T.TransId           //F_NB_JIAOYBNXH	交易包内序号	INT
		jiessjmsg.FNbKalx = T.ICCard.CardType        //F_NB_KALX	卡类型	INT  储值卡22，23 记账卡
		jiessjmsg.FNbJiaoyqye = T.ICCard.PreBalance  //F_NB_JIAOYQYE	交易前余额	分转元 INT
		jiessjmsg.FNbJiaoyhye = T.ICCard.PostBalance //F_NB_JIAOYHYE	交易后余额	分转元 INT
		jiessjmsg.FNbJine = T.Fee                    //F_NB_JINE	金额	INT         分转元
		jiessjmsg.FDtJiaoysj = T.Time                //F_DT_JIAOYSJ	交易时间	DATETIME   2020-05-13 14:34:34
		jiessjmsg.FDtJiaoylx = T.Service.ServiceType //F_DT_JIAOYLX	交易类型	VARCHAR(32)
		//jiessjmsg.FDtYonghrksj =T. //F_DT_YONGHRKSJ	用户入口时间	DATETIME
		jiessjmsg.FNbYonghtcsc = T.ParkTime //F_NB_YONGHTCSC	用户停车时长(分)	INT  天时分秒
		_, err := db.orm.Insert(jiessjmsg)
		if err != nil {
			log.Println("联网中心 新增结算数据 时 错误", err)
			return err
		}
	}
	log.Println("联网中心 新增结算处理数据 时 成功")
	return nil
}

//查询记账包数据  	//查询记账状态为1的数据，表示为已记账
func QueryKeepAccountdata() (error, *map[int64]*[]Jiessjchuli, *[]JieSuanMessage) {
	db := NewDatabase()
	//查询多条数据
	tests := make([]Jiessjchuli, 0)

	msgerr, msgids := QueryKeepAccountMsgdata()
	if msgerr != nil {
		log.Println("查询可以记账的原始交易消息包出错", msgerr)
	}

	jzshuju := make(map[int64]*[]Jiessjchuli)

	for _, mid := range *msgids {

		//同一个数据包 可以记账的数据
		qerr := db.orm.Where("FNbYuansjybxh=?", mid.MessageId).And("FNbJizjg=?", 1).Find(&tests)
		if qerr != nil {
			log.Fatalln("查询结算数据出错", qerr)
			return qerr, nil
		}
		log.Printf("总共查询出 %d 条数据\n", len(tests))

		if len(tests) == mid.Count {
			log.Println("此包没有争议数据")
		}
		for _, v := range tests {
			log.Printf("记账状态: %d, 原始交易包序号: %d\n", v.FNbJizjg, v.FNbYuansjybxh)
		}
		//消息包中可以记账的数据
		jzshuju[mid.MessageId] = (&tests)
	}

	return nil, &jzshuju, msgids
}

//查询出需要记账的消息包
func QueryKeepAccountMsgdata() (error, *[]JieSuanMessage) {
	db := NewDatabase()
	//查询多条数据
	tests := make([]JieSuanMessage, 0)
	//查询可以记账的 消息记录
	qerr := db.orm.Where("ji_zhang_zt=?", 1).Find(&tests)
	if qerr != nil {
		log.Fatalln("查询原始交易包数据出错", qerr)
		return qerr, nil
	}
	log.Printf("总共查询出 %d 条数据\n", len(tests))
	for _, v := range tests {
		log.Printf("记账状态: %d, 原始交易包序号: %d\n", v.JiZhangZt, v.MessageId)
	}

	return nil, &tests
}

//查询争议包数据
func QueryDisputedata() {

}

//查询争议包数据
func QueryClearingdata() {

}

//工具函数 把明细表置为1
func Updatedata() error {
	db := NewDatabase()
	//查询多条数据
	tests := new(Jiessjchuli)
	for i := 100347; i <= 100392; i++ {
		//测试 更新记账状态
		tests.FNbJizjg = 1

		count, err := db.orm.Table("jiessjchuli").Where("f_nb_yuansjybxh=?", i).Update(tests)
		if err != nil {
			log.Println("更新打包状态 失败", err, count)
			return err
		}
		log.Printf("更新包号 %d 打包状态 成功%d", i, count)

	}
	return nil
}

//工具函数 把原始交易包记账状态置为可以记账
func UpdateJZdata() error {
	db := NewDatabase()
	//查询多条数据
	tests := new(JieSuanMessage)
	for i := 100347; i <= 100392; i++ {
		//测试 更新记账状态 表示在记账中
		tests.JiZhangZt = 1

		count, err := db.orm.Table("jie_suan_message").Where("message_id=?", i).Update(tests)
		if err != nil {
			log.Println("更新打包状态 失败", err, count)
			return err
		}
		log.Printf("更新包号 %d 打包状态 成功%d", i, count)

	}
	return nil
}
