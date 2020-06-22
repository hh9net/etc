package centerServer

import (
	"CenterSettlement-go/log"
	"encoding/xml"
	"time"
)

//一个原始交易数据包
type Message struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type Header struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int      //消息传输的机制5
	MessageType  int      //消息的应用类型7
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节
}

type Body struct {
	XMLName           xml.Name      `xml:"Body"`
	ContentType       int           `xml:",attr"` //始终为1
	ClearTargetDate   string        //日期 如：2017-06-05 清分目标日期：取当前日期
	ServiceProviderId string        //通行宝中心系统Id，表示消息包中的交易是由收费方产生的
	IssuerId          string        //发行服务机构Id， 表示产生交易记录的发行服务机构。
	MessageId         int64         //交易消息包Id。配置文件中获得
	Count             int           //本消息包含的记录数量
	Amount            string        //交易总金额(元) 数据库为分【注意转换的小数问题】
	Transaction       []Transaction //交易原始数据

}
type Transaction struct {
	XMLName             xml.Name   `xml:"Transaction"`
	TransId             int        // 包内顺序Id，从1开始递增 ，包内唯一的交易记录
	Time                string     //交易的发生时间，需要加TAC计算 2020-05-13 14:34:34
	Fee                 string     //交易的发生金额(元)
	Service             Service    //服务信息
	ICCard              ICCard     //IC卡信息
	Validation          Validation //与校验相关的信息
	OBU                 OBU        //参加交易的电子标签信息
	CustomizedData      string     //特定发行方与通行宝收费方之间 约定格式的交易信息【  】
	Id                  string     //停车场消费交易编号(停车场编号+交易发生的时间+流水号 )
	Name                string     `xml:"name"`                //停车场名称(不超过150个字符)
	ParkTime            int        `xml:"parkTime"`            //停放时长(单位：分)
	VehicleType         string     `xml:"vehicleType"`         //收费车型
	AlgorithmIdentifier int        `xml:"algorithmIdentifier"` //算法标识 1-3DEX  2-SM4

}

//服务信息
type Service struct {
	XMLName     xml.Name `xml:"Service"`
	ServiceType string   //交易的服务类型【【2 写死】】
	Description string   //对交易的文字解释【停车场名｜停车时常 ：几时几分几秒】
	Detail      string   //交易详细信息  1|04|3201|3201020001|1104|20200513 143434|03|3201|3201020001|0001|20200513 140805
}

//IC卡信息
type ICCard struct {
	XMLName     xml.Name `xml:"ICCard"`
	CardType    int      //卡类型，22为储值卡；23记账卡
	NetNo       string   //网络编码，BCD码 Hex(4) ka网络号（16进制） 数据库10进制
	CardId      string   //IC卡物理编号，BCD码  Hex(16)   卡号
	License     string   //0015文件中记录的车牌号 卡内车牌号 【含颜色】
	PreBalance  string   //交易前余额，以元为单位 Decimal
	PostBalance string   //交易后余额，以元为单位 Decimal
}

//主要用于TAC计算
type Validation struct {
	XMLName         xml.Name `xml:"Validation"`
	TAC             string   //交易时产生的TAC码，8位16进制数   Hex(8)
	TransType       string   //交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交易  Hex(2)【09】
	TerminalNo      string   //12位16进制数据，即PSAM号，PSAM中0016文件中的终端机编号  Hex(2) 	加密卡号
	TerminalTransNo string   //8位16进制数，PSAM卡脱机交易序号，在MAC1计算过程中得到  Hex(8) 加密序列号
}

type OBU struct {
	XMLName  xml.Name `xml:"OBU"`
	NetNo    string   //4501  OBU网络号
	OBUId    string   //OBU物理编号，BCD码  4501191509252866
	OBEState string   //2字节的OBU状态
	License  string   //OBU中记录的车牌号 【含颜色】
}

//结算处理表
type Jiessjchuli struct {
	//WeiyiId int64 `xorm:"pk autoincr"`
	//FVcJiaoyjlid   string    //`xorm:"pk"` //F_VC_JIAOYJLID	交易记录ID	VARCHAR(128)
	FNbYuansjybxh  int64     //F_NB_YUANSJYBXH	原始交易包序号	BIGINT
	FNbJiaoybnxh   int       //F_NB_JIAOYBNXH	交易包内序号	INT
	FNbJizjg       int       //F_NB_JIZJG	记账结果	INT "0：未记账  1：已记账    2：争议数据"
	FNbZhengylx    int       //F_NB_ZHENGYLX	争议类型	INT 0，不是争议，1，验证未通过
	FNbJizbxh      int64     //F_NB_JIZBXH	记账包序号	INT
	FNbZhengyclbxh int64     //F_NB_ZHENGYCLBXH	争议处理包序号	BIGINT  记账结果：争议放过、坏账时
	FNbQingfbxh    int64     //F_NB_QINGFBXH	清分包序号	BIGINT
	FNbKalx        int       //F_NB_KALX	卡类型	INT  储值卡22，23 记账卡
	FNbJiaoyqye    string    //F_NB_JIAOYQYE	交易前余额	分转元 INT
	FNbJiaoyhye    string    //F_NB_JIAOYHYE	交易后余额	分转元 INT
	FNbJine        string    //F_NB_JINE	金额	INT         分转元
	FDtJiaoysj     string    //F_DT_JIAOYSJ	交易时间	DATETIME   2020-05-13 14:34:34
	FDtJiaoylx     string    //F_DT_JIAOYLX	交易类型	VARCHAR(32)
	FDtYonghrksj   time.Time //F_DT_YONGHRKSJ	用户入口时间	DATETIME
	FNbYonghtcsc   int       //F_NB_YONGHTCSC	用户停车时长(分)	INT  天时分秒
	FNbQingfjg     int       //F_NB_QINGFJG  '清分结果 0：未清分、1：已清分'
	FNbZhengycljg  int       //F_NB_ZHENGYCLJG  '争议处理结果 0:未处理、1：争议放过、2：坏账'
	FNbJusbj       int       //`F_NB_JUSBJ`   DEFAULT '0' COMMENT '拒收标记 0、正常，1、拒收',
	FVcQingfmbr    string    // `F_NB_QINGFMBR` int(11) DEFAULT NULL COMMENT '清分目标日',

}

type JieSuanMessage struct {
	//WeiyiId int64 `xorm:"pk autoincr"`
	Version      string //统一 00010000 Hex(8) Header
	MessageClass int    //消息传输的机制5
	MessageType  int    //消息的应用类型7
	SenderId     string // Hex(16位，不足补零) 发送方Id
	ReceiverId   string //Hex(16位，不足补零) 接收方Id
	MessageId    int64  //消息序号，从1开始，逐1递增 ，8字节

	ClearTargetDate   string //日期 如：2017-06-05 清分目标日期：取当前日期
	ServiceProviderId string //通行宝中心系统Id，表示消息包中的交易是由收费方产生的
	IssuerId          string //发行服务机构Id， 表示产生交易记录的发行服务机构。
	Count             int    //本消息包含的记录数量
	Amount            string //交易总金额(元) 数据库为分【注意转换的小数问题】

	JiZhangZt int //0：未记账处理 1：在记账中 2：已记账

	ZhengyChulizt int //0：无争议 1：有争议

}

type JieSuanMessageMx struct {
	//WeiyiId int64 `xorm:"pk autoincr"`
	Version      string //统一 00010000 Hex(8) Header
	MessageClass int    //消息传输的机制5
	MessageType  int    //消息的应用类型7
	SenderId     string // Hex(16位，不足补零) 发送方Id
	ReceiverId   string //Hex(16位，不足补零) 接收方Id
	MessageId    int64  //消息序号，从1开始，逐1递增 ，8字节

	ClearTargetDate   string //日期 如：2017-06-05 清分目标日期：取当前日期
	ServiceProviderId string //通行宝中心系统Id，表示消息包中的交易是由收费方产生的
	IssuerId          string //发行服务机构Id， 表示产生交易记录的发行服务机构。

	TransId int    // 包内顺序Id，从1开始递增 ，包内唯一的交易记录
	Time    string //交易的发生时间，需要加TAC计算 2020-05-13 14:34:34
	Fee     string //交易的发生金额(元)

	CustomizedData      string //特定发行方与通行宝收费方之间 约定格式的交易信息【  】
	BhId                string //停车场消费交易编号(停车场编号+交易发生的时间+流水号 )
	Name                string //停车场名称(不超过150个字符)
	ParkTime            int    //停放时长(单位：分)
	VehicleType         string //收费车型
	AlgorithmIdentifier int    //算法标识 1-3DEX  2-SM4

	ServiceType string //交易的服务类型【【2 写死】】
	Description string //对交易的文字解释【停车场名｜停车时常 ：几时几分几秒】
	Detail      string //交易详细信息  1|04|3201|3201020001|1104|20200513 143434|03|3201|3201020001|0001|20200513 140805

	CardType    int    //卡类型，22为储值卡；23记账卡
	NetNo       string //网络编码，BCD码 Hex(4) ka网络号（16进制） 数据库10进制
	CardId      string //IC卡物理编号，BCD码  Hex(16)   卡号
	License     string //0015文件中记录的车牌号 卡内车牌号 【含颜色】
	PreBalance  string //交易前余额，以元为单位 Decimal
	PostBalance string //交易后余额，以元为单位 Decimal

	TAC             string //交易时产生的TAC码，8位16进制数   Hex(8)
	TransType       string //交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交易  Hex(2)【09】
	TerminalNo      string //12位16进制数据，即PSAM号，PSAM中0016文件中的终端机编号  Hex(2) 	加密卡号
	TerminalTransNo string //8位16进制数，PSAM卡脱机交易序号，在MAC1计算过程中得到  Hex(8) 加密序列号

	OBUId    string //OBU物理编号，BCD码  4501191509252866
	OBEState string //2字节的OBU状态

}

type XormTest struct {
	WeiyiId int64 `xorm:"pk autoincr"`
	Name    string
}

func XorminitTest() {
	db := NewDatabase()
	is, err := db.orm.IsTableExist(
		//new(JieSuanMessage),
		new(XormTest),
		//new(JieSuanMessageMx),
	)
	if err != nil {
		log.Fatal("创建数据库判断数据库表是否存在时 error  ", err)
	}
	if is == false {
		err := db.orm.Sync2(
			new(XormTest),
		)
		//err = db.orm.CreateTables(new(SJsJiessj))
		if err != nil {
			log.Fatal()
			log.Fatal("创建数据库 映射表 error ", err)
		}
	} else {
		log.Fatal("创建数据库 映射表 表已存在")
	}
}

func XormInsert() {
	db := NewDatabase()
	x := new(XormTest)
	x.Name = "sun"

	db.orm.Insert(x)

}
