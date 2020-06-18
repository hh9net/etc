package types

import (
	"encoding/xml"
	"time"
)

//清分统计处理消息结构
//	清分统计处理（可疑帐调整数据）
type ClearingMessage struct {
	XMLName xml.Name       `xml:"Message"`
	Header  ClearingHeader `xml:"Header"`
	Body    ClearingBody   `xml:"Body"`
}

type ClearingHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int      //消息传输的机制 5
	MessageType  int      //消息的应用类型 5
	SenderId     string   //Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

//清分统计body
type ClearingBody struct {
	XMLName         xml.Name `xml:"Body"`
	ContentType     int      `xml:",attr"` //争议消息的ContentType始终为2
	ClearTargetDate string   //清分目标日
	Count           int      //对应清分总金额的的交易记录数量，包含原始交易包中由发行方确认应付的交易数量和争议处理结果中确认应付的交易数量之和， 不包含争议处理结果中为坏帐的记录数量。
	Amount          string   //清分总金额(确认付款金额)
	ProcessTime     string   //处理时间 清分统计处理时间
	IssuerId        string   //发行服务机构Id，
	List            List
}

//ProcessTime之后的内容，或者包含IssuerId及0或多个List，或者仅包含一个List。
//前者是联网中心为发行方产生的清分统计结果，后者是为通行宝中心系统产生的清分统计结果。

type List struct {
	XMLName      xml.Name `xml:"List"`
	MessageCount int      `xml:",attr"` //本次清分包含的所有原始交易包数量
	FileCount    int      `xml:",attr"` //本次清分包含的所有争议处理结果包数量

	ServiceProviderId string //通行宝中心系统Id
	MessageId         []MessageId
	FileId            int //争议处理结果文件Id (含有可疑调整数据的有)
}

type MessageId struct {
	MessageId int64 //通行宝中心系统发的原始交易包Id 由于此处为messageid的明细，
	// 所以建议收费方上传打包时尽量在一个messageid中多打包原始记录（小于10000条），
	//减少此处的数据量，因为该包格式为单包
}

//List说明本次清分所包含的交易记录范围。该范围可以通过交易记录消息包MessageId和争议处理结果FileId确定。
//<?xml version="1.0" encoding="UTF-8"?>
//<Message>
//<Header>
//<Version>00010000</Version>
//<MessageClass>5</MessageClass>
//<MessageType>5</MessageType>
//<SenderId>0000000000000020</SenderId>
//<ReceiverId>00000000000000FD</ReceiverId>
//<MessageId>114462</MessageId>
//</Header>
//<Body ContentType="2">
//<ClearTargetDate>2018-06-30</ClearTargetDate>
//<Amount>0.00</Amount>
//<Count>0</Count>
//<ProcessTime>2018-07-02T09:18:05</ProcessTime>
//<List MessageCount="0" FileCount="0">
//<ServiceProviderId>00000000000000FD</ServiceProviderId>
//<MessageId/>
//</List>
//</Body>
//</Message>

//应答争议处理
type ResClearingMessage struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type ResClearingHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32    //消息传输的机制6
	MessageType  int32    //消息的应用类型5
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

type ResClearingBody struct {
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
//8.	清分统计对帐失败
