package server

import "testing"

func TestReceive(t *testing.T) {
	Receive()
}
func TestHandleTable(t *testing.T) {
	//线程1 数据打包 压缩
	HandleTable()
}

func TestAnalyzeDataPakage(t *testing.T) {
	//线程4   解析
	AnalyzeDataPakage()
}
