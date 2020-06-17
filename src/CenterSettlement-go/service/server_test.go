package service

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

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
	fname := ""
	fileInfoList, rerr := ioutil.ReadDir("../keepAccountFile/")
	if rerr != nil {
		log.Fatal(rerr)
	}
	for i := range fileInfoList {
		fname = fileInfoList[i].Name()

		path := "../keepAccountFile/"

		Parsexml(path, fname)
	}
	//{{ Message} {{ Header} 00020000 5 5 0000000000000020 00000000000000FD 181196} {{ Body} 1 00000000000000FD 0000000000000020 221908 2020-06-08T20:45:40 1 4.00 0}}
	//2020/06/16 18:28:10 <nil>

	//
	//f1 := "JZB-ok_00000000000000181191.xml"
	//Parsexml(path, f1)
	//
	//f2 := "JZB-ok_00000000000000181192.xml"
	//Parsexml(path, f2)
	//
	//f3 := "JZB-ok_00000000000000181193.xml"
	//Parsexml(path, f3)
	//
	//f4:= "JZB-ok_00000000000000181194.xml"
	//Parsexml(path, f4)
	//f5:= "JZB-ok_00000000000000181195.xml"
	//Parsexml(path, f5)

}
