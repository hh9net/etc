package types

import (
	"encoding/xml"
	"time"
)

//通用确认消息结构
//ContentType 存在
type ResponseCTMessage struct {
	XMLName xml.Name       `xml:"Message"`
	Header  ResponseHeader `xml:"Header"`
	Body    ResponseCTBody `xml:"Body"`
}

type ResponseHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int      //消息传输的机制
	MessageType  int      //消息的应用类型
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节
}

type ResponseCTBody struct {
	XMLName     xml.Name `xml:"Body"`
	ContentType int      `xml:",attr"`
	MessageId   int64    //确认的消息Id （对应于发送方的header里的MessageId）从发送方的header取MessageId
	ProcessTime string   //处理时间
	Result      int      // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
}
type ResponseMessage struct {
	XMLName xml.Name       `xml:"Message"`
	Header  ResponseHeader `xml:"Header"`
	Body    ResponseBody   `xml:"Body"`
}
type ResponseBody struct {
	XMLName     xml.Name `xml:"Body"`
	MessageId   int64    //确认的消息Id （对应于发送方的header里的MessageId）从发送方的header取MessageId
	ProcessTime string   //处理时间
	Result      int      // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
}

//  通用重发请求消息结构
type ResendMessage struct {
	XMLName xml.Name     `xml:"Message"`
	Header  ResendHeader `xml:"Header"`
	Body    ResendBody   `xml:"Body"`
}

type ResendHeader struct {
	XMLName xml.Name `xml:"Header"`
	//Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32  //消息传输的机制 1，Request
	MessageType  int32  //消息的应用类型  请求重发的数据类型对应的MessageType
	SenderId     string // Hex(16位，不足补零) 发送方Id  请求重发消息方Id
	ReceiverId   string //Hex(16位，不足补零) 接收方Id    接收消息的参与方Id
	//MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节
}

//通用重发请求消息中没有更多的数据，其Body为空。
type ResendBody struct {
	XMLName         xml.Name  `xml:"Body"`
	MessageId       string    //确认的消息Id （对应于发送方的header里的MessageId）从发送方的header取MessageId
	ProcessTime     time.Time //处理时间
	Result          int       // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
	ClearTargetDate string    //ClearTargetDate  目标清分日
}
