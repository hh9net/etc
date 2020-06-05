package server

import (
	"CenterSettlement-go/lz77zip"
	"testing"
)

func TestServer(t *testing.T) {
	//模拟联网中心，处理结算数据
	Server()
}

func TestZip(t *testing.T) {
	lz77zip.Zip("jz00001.xml")
}
func TestCenterClient(t *testing.T) {
	//模拟联网中心，发送数据
	CenterClient()
}
