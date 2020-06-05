package types

import "encoding/xml"

type ReceiveMessage struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type ReceiveHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32    //消息传输的机制
	MessageType  int32    //消息的应用类型
	SenderId     string   //Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

//jibody
type ReceiveBody struct {
	XMLName     xml.Name `xml:"Body"`
	ContentType int      `xml:",attr"` //争议消息的ContentType始终为2
}

//
