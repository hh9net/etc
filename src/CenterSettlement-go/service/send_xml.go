package service

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/types"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

//线程2 发送数据包
func HandleSendXml() {
	//从文件夹sendzipxml中扫描打包文件（判断这个文件夹下面有没有文件）
	tiker := time.NewTicker(time.Second * 3)
	for {
		log.Println(<-tiker.C)
		log.Println("执行线程2")

		//扫描receive 文件夹 读取文件
		//获取文件或目录相关信息
		pwd := "../sendzipxml/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		for i := range fileInfoList {
			//判断文件的结尾名
			if strings.HasSuffix(fileInfoList[i].Name(), ".xml.lz77") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())

				//		解析文件  获取数据
				var sendStru types.SendStru
				sendStru.Massageid = "00000000000000100025"
				sendStru.Md5_str = "@234567812345678123456781234567@"
				sendStru.Xml_length = "#1122#"
				sendStru.Xml_msgName = fileInfoList[i].Name()

				//连接联网中心服务器
				conn, derr := net.Dial("tcp", "127.0.0.1:8808")
				if derr != nil {
					log.Println("Dial", derr)
					//return ""
				}
				if conn != nil {
					log.Println("Dial 成功")
				}
				client.Sendxml(&sendStru, conn)
				buf := make([]byte, 1024)
				n, err2 := conn.Read(buf)
				if err2 != nil {
					log.Println("conn.Read err = ", err2)
					return
				}
				//str := string(buf[:n])
				//对联网中心的接收应答处理
				ImmediateResponseProcessing(string(buf[:n]), fileInfoList[i].Name())

				if "no" == string(buf[:n]) {
					log.Println("发送失败")

				}
				conn.Close()
			}
		}
	}
}

//即使应答处理
func ImmediateResponseProcessing(str string, name string) {
	log.Println(str)
	b := []byte(str)
	if string(b[len(str)-1:]) == "1" {

		log.Println("收到联网中心的接收成功的即时应答")
		log.Println("发送成功")
		//成功后 mv文件夹到另一个文件中
		s := "../sendzipxml/" + name
		des := "../sendxmlsucceed/" + name
		client.MoveFile(s, des)

	}
	if string(b[len(str)-1:]) == "0" {
		//
		log.Println("联网中心的接收失败")
		log.Println("发送失败")
		//	发送失败 触发重发机制
	}

}
