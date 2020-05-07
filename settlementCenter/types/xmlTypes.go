package types

import (
	"encoding/xml"
	"time"
)

//一个交易包
type Message struct {
	XMLName xml.Name `xml:"Message"`
	Header  `xml:"Message>Header"`
	Body    `xml:"Message>Body"`
}

type Header struct {
	Version      string `xml:"Message>Header>Version"`      //统一 00010000 Hex(8) Header
	MessageClass int32  `xml:"Message>Header>MessageClass"` //消息传输的机制
	MessageType  int32  `xml:"Message>Header>MessageType"`  //消息的应用类型
	SenderId     string `xml:"Message>Header>SenderId"`     // Hex(16位，不足补零)发送方Id
	ReceiverId   string `xml:"Message>Header>ReceiverId"`   //Hex(16位，不足补零)接收方Id
	MessageId    int64  `xml:"Message>Header>MessageId"`    //消息序号，从1开始，逐1递增 ，8字节
}

type Body struct {
	ContentType       int           `xml:"ContentType,attr"`               //始终为1
	ClearTargetDate   time.Time     `xml:"Message>Body>ClearTargetDate"`   //日期 如：2017-06-05
	ServiceProviderId string        `xml:"Message>Body>ServiceProviderId"` //通行宝中心系统Id，表示消息包中的交易是由收费方产生的
	IssuerId          string        `xml:"Message>Body>IssuerId"`          //发行服务机构Id， 表示产生交易记录的发行服务机构。
	MessageId         int64         `xml:"Message>Body>MessageId"`         //交易消息包Id。
	Count             int           `xml:"Message>Body>Count"`             //本消息包含的记录数量
	Amount            int           `xml:"Message>Body>Amount"`            //交易总金额(分)
	Transaction       []Transaction //交易原始数据   `xml:"Message>Body>Transaction"`
	Description       string        `xml:",innerxml"`
}
type Transaction struct {
	TransId        int       // 包内顺序Id，从1开始递增 ，包内唯一的交易记录
	Time           time.Time //交易的发生时间，需要加TAC计算
	Fee            int       //交易的发生金额(分)
	Service                  //服务信息
	ICCard                   //IC卡信息
	Validation               //与校验相关的信息
	OBU                      //参加交易的电子标签信息
	CustomizedData string    //特定发行方与通行宝收费方之间 约定格式的交易信息
}

type Service struct {
	ServiceType int    //交易的服务类型
	Description string //对交易的文字解释
	Detail      string //交易详细信息
}

type ICCard struct {
	CardType    int    //卡类型，22为储值卡；23记账卡
	NetNo       string //网络编码，BCD码 Hex(4)
	CardId      string //IC卡物理编号，BCD码  Hex(16)
	License     string //0015文件中记录的车牌号
	PreBalance  string //交易前余额，以元为单位 Decimal
	PostBalance string //交易后余额，以元为单位 Decimal
}

//主要用于TAC计算
type Validation struct {
	TAC             string //交易时产生的TAC码，8位16进制数   Hex(8)
	TransType       string //交易标识，2位16进制数，PBOC定义，如06为传统交易，09为复合交易  Hex(2)
	TerminalNo      string //12位16进制数据，即PSAM号，PSAM中0016文件中的终端机编号  Hex(2)
	TerminalTransNo string //8位16进制数，PSAM卡脱机交易序号，在MAC1计算过程中得到  Hex(8)
}

type OBU struct {
	NetNo    string //4501
	OBUId    string //OBU物理编号，BCD码  4501191509252866
	OBEState string //2字节的OBU状态
	License  string //OBU中记录的车牌号
}
