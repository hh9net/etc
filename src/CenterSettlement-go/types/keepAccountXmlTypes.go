package types

import (
	"encoding/xml"
	"time"
)

//记账处理消息结构
type KeepAccountMessage struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type KeepAccountHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32    //消息传输的机制5
	MessageType  int32    //消息的应用类型5
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节
}

type KeepAccountBody struct {
	XMLName           xml.Name      `xml:"Body"`
	ContentType       int           `xml:",attr"` //始终为1
	ClearTargetDate   time.Time     //日期 如：2017-06-05 清分目标日期：取当前日期
	ServiceProviderId string        //通行宝中心系统Id，表示消息包中的交易是由收费方产生的
	IssuerId          string        //发行服务机构Id， 表示产生交易记录的发行服务机构。
	MessageId         int64         //交易消息包Id。配置文件中获得
	Count             int           //本消息包含的记录数量
	Amount            int           //交易总金额(元) 数据库为分【注意转换的小数问题】
	Transaction       []Transaction //交易原始数据

}
