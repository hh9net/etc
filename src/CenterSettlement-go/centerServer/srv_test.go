package centerServer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	//模拟联网中心，处理结算数据
	Server()
	//处理文件的解析
	//HandleFile()
	//处理文件的发送
	//CenterClient()
}

func TestHandleFile(t *testing.T) {
	//模拟联网中心，处理结算数据
	//Server()
	//处理文件的解析
	HandleFile()
	//处理文件的发送
	//CenterClient()
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

func TestGetFileMd5(t *testing.T) {
	f := "00000000000000010612.xml"
	s := GetFileMd5(f)
	log.Println(s)
}

func TestUnzip(t *testing.T) {
	f := "00000000000000010612.xml"
	s := UnZipLz77(f)
	log.Println(s)
}

func TestToABCD(t *testing.T) {
	s := "12asd456fdghfsh"
	fmt.Println("s =", s)
	str := strings.ToUpper(s)
	fmt.Println("str =", str)
}
