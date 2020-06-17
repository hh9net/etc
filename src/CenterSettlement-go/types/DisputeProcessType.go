package types

import (
	"encoding/xml"
	"time"
)

//争议处理消息结构
//	争议交易处理（可疑帐调整数据）
type DisputeProcessMessage struct {
	XMLName xml.Name             `xml:"Message"`
	Header  DisputeProcessHeader `xml:"Header"`
	Body    DisputeProcessBody   `xml:"Body"`
}

type DisputeProcessHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int      //消息传输的机制 5
	MessageType  int      //消息的应用类型 7
	SenderId     string   // Hex(16位，不足补零) 发送方Id 联网中心id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id 通行宝
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

//争议body
type DisputeProcessBody struct {
	XMLName            xml.Name `xml:"Body"`
	ContentType        int      `xml:",attr"` //争议消息的ContentType始终为2
	ClearingOperatorId string   //清分方ID
	ServiceProviderId  string   //发行服务id， 联网中心id
	IssuerId           string   //发行方服务机构Id，通行宝
	ProcessTime        string   //处理时间
	Count              int      //本消息包含的记录数量，包括经讨论确认付款和记录和坏帐 记录数量

	Amount      string        //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】浮点数
	FileId      int           //争议结果文件Id ，使用此栏位确认争议归属清算日期 （现使用OBU中心清算日ClearTargetDate作为唯一标识,YYYYMMDD） 如果争议结果文件每天需要交互多次，建议FileID编码规则改为前8位为争议处理结果归属清算日yyyymmdd,加1位序号
	MessageList []MessageList //包含争议结果 记录的原始交易包
}

//争议交易记录List的格式
//表示其包含的争议结果记录来自于收费方某个原始交易包
type MessageList struct {
	XMLName     xml.Name                  `xml:"MessageList"`
	MessageId   int64                     //通行宝中心系统原始交易包ID
	Count       int                       //属于当前交易包的争议结果记录数量
	Amount      string                    //属于当前交易包的争议结果记录金额
	Transaction DisputeProcessTransaction //争议处理结果记录
}

type DisputeProcessTransaction struct {
	TransId int //表示该条交易在原始数据包中的交易记录Id
	Result  int //处理结果 为0表示正常支付；为1表示此交易作坏账处理。
}

//<?xml version="1.0" encoding="UTF-8"?>
//<Message>
//<Header>
//<Version>00010000</Version>
//<MessageClass>5</MessageClass>
//<MessageType>7</MessageType>
//<SenderId>0000000000000020</SenderId>
//<ReceiverId>00000000000000FD</ReceiverId>
//<MessageId>114464</MessageId>
//</Header>
//<Body ContentType="2">
//<ClearingOperatorId>0000000000000020</ClearingOperatorId>//清分方id
//<ServiceProviderId>0000000000000020</ServiceProviderId>//发行方
//<IssuerId>00000000000000FD</IssuerId>//通行宝
//<FileId>20180702</FileId>
//<ProcessTime>2018-07-03T09:06:04</ProcessTime>
//<Count>2</Count>
//<MessageList>
//<MessageId>6410</MessageId>
//<Count>1</Count>
//<Amount>3.00</Amount>
//<Transaction>
//<TransId>1</TransId>
//<Result>1</Result>
//</Transaction>
//</MessageList>
//<MessageList>
//<MessageId>6462</MessageId>
//<Count>1</Count>
//<Amount>2.00</Amount>
//<Transaction>
//<TransId>1</TransId>
//<Result>1</Result>
//</Transaction>
//</MessageList>
//</Body>
//</Message>

//应答争议处理
type ResDisputeProcessMessage struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type ResDisputeProcessHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32    //消息传输的机制6
	MessageType  int32    //消息的应用类型7
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

type ResDisputeProcessBody struct {
	XMLName     xml.Name  `xml:"Body"`
	ContentType int       `xml:",attr"` //争议消息的ContentType始终为2
	MessageId   int64     //确认的消息id
	ProcessTime time.Time //处理时间
	Result      int       //执行结果
}

//执行结果：
//1.	消息已正常接收（用于Advice Response时含已接受建议）
//2.	消息头错误，如MessageClass或MessageType不符合定义，SenderId不存在等
//3.	消息格式不正确，即XML Schema验证未通过
//4.	消息格式正确但内容错误，包括数量不符，内容重复等
//5.	消息重复
//6.	消息正常接收，但不接受建议（仅用于Advice Response）
//7.	消息版本错误
