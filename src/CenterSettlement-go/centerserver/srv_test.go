package centerserver

import (
	"CenterSettlement-go/lz77zip/Cgo"
	"testing"
)

func TestServer(t *testing.T) {
	//模拟联网中心，处理结算数据
	//Server()
	//处理文件的解析
	HandleFile()
	//处理文件的发送
	//CenterClient()
}

func TestZip(t *testing.T) {
	Cgo.Zip("jz00001.xml")
}
func TestCenterClient(t *testing.T) {
	//模拟联网中心，发送数据
	CenterClient()
}

func TestNewDatabase(t *testing.T) {
	NewDatabase()
}
func TestDB_NewTable(t *testing.T) {
	db := new(DB)

	db.NewTable()
}
