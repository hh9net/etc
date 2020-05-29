package types

//tcp通讯报文格式
type SendStru struct {
	Massageid   string //消息报文序号  20字节Char型，不足左补0
	Xml_length  string //压缩后XML消息包长度
	Md5_str     string //32字节MD5校验串   MD5值由<Message>
	Xml_msgName string //二进制压缩后的XML 消息包 .xml.lz77
}

//通讯报文应答格式
type ReplyStru struct {
	Massageid string //消息报文序号【20字节】
	Result    string //【1字节】	1 成功接收 0 接收超时、长度不符
}
