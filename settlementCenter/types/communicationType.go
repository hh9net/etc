package types

//tcp通讯报文格式
type SendStru struct {
	Massageid  [20]byte //消息报文序号
	Xml_length [6]byte  //压缩后XML消息包长度
	Md5_str    [32]byte //32字节MD5校验串
	Xml_msg    string   //二进制压缩后的XML 消息包
}

//通讯报文应答格式
type ReplyStru struct {
	Massageid [20]byte //消息报文序号【20字节】
	result    byte     //【1字节】	1 成功接收 0 接收超时、长度不符
}
