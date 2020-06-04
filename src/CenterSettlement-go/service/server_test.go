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
