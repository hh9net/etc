package service

import "testing"

//测试 原始交易数据打包  xml生成
func TestHandleGeneratexml(t *testing.T) {
	//线程1 数据打包
	HandleGeneratexml()
}

//测试 xml 发送
func TestHandleSendXml(t *testing.T) {
	//线程2
	HandleSendXml()
}

func TestReceive(t *testing.T) {
	//线程 3
	Receive()
}

func TestAnalyzeDataPakage(t *testing.T) {
	//线程4   解析文件
	AnalyzeDataPakage()
}

//测试xml文件解析
func TestParseFile(t *testing.T) {
	ParseFile()
}

//测试清分包xml文件解析
func TestParsexml(t *testing.T) {
	path := "../clearlings/"
	fname := "QFB_00000000000000114462.xml"
	Parsexml(path, fname)

}
