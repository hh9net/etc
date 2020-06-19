package centerServer

import "encoding/xml"

//记账处理消息结构
//接收联网中心的记账处理xml消息
type KeepAccountokMessage struct {
	XMLName xml.Name          `xml:"Message"`
	Header  KeepAccountHeader `xml:"Header"`
	Body    KeepAccountOkBody `xml:"Body"`
}

type KeepAccountHeader struct {
	XMLName      xml.Name `xml:"Header"`
	Version      string   //统一 00010000 Hex(8) Header
	MessageClass int      //消息传输的机制5
	MessageType  int      //消息的应用类型5
	SenderId     string   // Hex(16位，不足补零) 发送方Id
	ReceiverId   string   //Hex(16位，不足补零) 接收方Id
	MessageId    string   //消息序号，从1开始，逐1递增 ，8字节  记账包的消息id
}

//如果记账没有争议
type KeepAccountOkBody struct {
	XMLName           xml.Name `xml:"Body"`
	ContentType       int      `xml:",attr"` //记帐消息的ContentType始终为1
	ServiceProviderId string   //通行宝中心系统Id，
	IssuerId          string   //发行服务机构Id， 记账消息哪一个发行方
	MessageId         int64    //交易消息包Id。原始交易包消息中的messageid
	ProcessTime       string   //处理时间
	Count             int      //本消息对应的原始交易包中交易记录的数量
	Amount            string   //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
	DisputedCount     int      //本消息包含的争议交易数量 （可疑帐笔数）
}

type KeepAccountMessage struct {
	XMLName xml.Name          `xml:"Message"`
	Header  KeepAccountHeader `xml:"Header"`
	Body    KeepAccountBody   `xml:"Body"`
}

//如果记账有争议
type KeepAccountBody struct {
	XMLName           xml.Name         `xml:"Body"`
	ContentType       int              `xml:",attr"` //记帐消息的ContentType始终为1
	ServiceProviderId string           //通行宝中心系统Id，
	IssuerId          string           //发行服务机构Id， 记账消息哪一个发行方
	MessageId         int64            //交易消息包Id。原始交易包消息中的messageid
	ProcessTime       string           //处理时间
	Count             int              //本消息对应的原始交易包中交易记录的数量
	Amount            string           //确认记帐总金额 交易总金额(元) 数据库为分【注意转换的小数问题】
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
	MessageId         []Messageid
	FileId            int //争议处理结果文件Id (含有可疑调整数据的有)
}

type Messageid struct {
	XMLName   xml.Name `xml:"MessageId"`
	MessageId int64    `xml:"MessageId"` //通行宝中心系统发的原始交易包Id 由于此处为messageid的明细，
	// 所以建议收费方上传打包时尽量在一个messageid中多打包原始记录（小于10000条），
	//减少此处的数据量，因为该包格式为单包
}
