package types

import (
	"encoding/xml"
	"time"
)

//记账处理消息结构
//接收联网中心的记账处理xml消息
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
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

//如果记账没有争议
type KeepAccountOkBody struct {
	XMLName           xml.Name  `xml:"Body"`
	ContentType       int       `xml:",attr"` //记帐消息的ContentType始终为1
	ServiceProviderId string    //通行宝中心系统Id，
	IssuerId          string    //发行服务机构Id， 记账消息哪一个发行方
	MessageId         int64     //交易消息包Id。原始交易包消息中的messageid
	ProcessTime       time.Time //处理时间
	Count             int       //本消息对应的原始交易包中交易记录的数量
	Amount            int       //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
	DisputedCount     int       //本消息包含的争议交易数量 （可疑帐笔数）
}

//如果记账有争议
type KeepAccountBody struct {
	XMLName           xml.Name         `xml:"Body"`
	ContentType       int              `xml:",attr"` //记帐消息的ContentType始终为1
	ServiceProviderId string           //通行宝中心系统Id，
	IssuerId          string           //发行服务机构Id， 记账消息哪一个发行方
	MessageId         int64            //交易消息包Id。原始交易包消息中的messageid
	ProcessTime       time.Time        //处理时间
	Count             int              //本消息对应的原始交易包中交易记录的数量
	Amount            int              //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
	DisputedCount     int              //本消息包含的争议交易数量 （可疑帐笔数）
	DisputedRecord    []DisputedRecord //争议记录内容
}

//记帐处理结果仅返回有争议（可疑）的交易记录明细。
//未包含在争议交易记录明细中的交易，均默认为发行方已确认可以付款。

//记录内容
type DisputedRecord struct {
	XMLName xml.Name `xml:"DisputedRecord"`
	TransId int      //表示该条交易在原始数据包中的交易记录Id 是由通行宝中心系统产生的该包内顺序Id，从1开始递增。
	Result  int      //处理结果

}

/*
<?xml version="1.0" encoding="UTF-8"?>
<Message>
<Header>
<Version>00020000</Version>
<MessageClass>5</MessageClass>
<MessageType>5</MessageType>
<SenderId>0000000000000020</SenderId>
<ReceiverId>00000000000000FD</ReceiverId>
<MessageId>114509</MessageId>
</Header>

<Body ContentType="1">
<ServiceProviderId>00000000000000FD</ServiceProviderId>
<IssuerId>0000000000000020</IssuerId>
<MessageId>6474</MessageId>
<ProcessTime>2018-07-03T10:15:36</ProcessTime>
<Count>29</Count>
<Amount>119.00</Amount>


//如果DisputedCount不为0
<DisputedCount>2</DisputedCount>
<DisputedRecord>
<TransId>7</TransId>
<Result>3</Result>
</DisputedRecord>
<DisputedRecord>
<TransId>9</TransId>
<Result>3</Result>
</DisputedRecord>
</Body>
</Message>*/

//记帐处理结果Result定义：
//取值	说明
//1	验证未通过（如：TAC错误）
//2	重复的交易信息
//3	用户状态变化
//4	无效交易类型
//5	逾期超过设定值
//6	交易数据域错
//7	超过最大交易限额
//8	卡号不存在
//9	卡状态不匹配
//10	卡超过有效期
//11	不允许的交易
//12	卡片CSN不匹配
//13	测试交易
//14	卡帐不符（仅用于储值卡）
//15	无效卡类型
//16	车道对时错误
//17	OBU号不存在
//100	其它

//应答记账处理
type ResKeepAccountMessage struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

type ResKeepAccountHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int32    //消息传输的机制6
	MessageType  int32    //消息的应用类型5
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    int64    //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

type ResKeepAccountBody struct {
	XMLName     xml.Name  `xml:"Body"`
	ContentType int       `xml:",attr"` //记帐消息的ContentType始终为1
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
