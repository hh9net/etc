package server

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	"CenterSettlement-go/lz77zip"
	"CenterSettlement-go/service"
	"CenterSettlement-go/types"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

//模拟联网中心发送 记账数据、争议数据、清分数据
func CenterClient() {
	//连接通行包
	address := conf.ListeningAddressConfigInit()
	Address := address.ListeningAddressIp + ":" + address.ListeningAddressPort
	//Dial
	conn, derr := net.Dial("tcp", Address)
	if derr != nil {
		log.Println("Dial", derr)
		//return ""
	}
	if conn != nil {
		log.Println("Dial 成功")
	}

	pwd := "CenterSettlement-go/generatexml/"
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
	for i := range fileInfoList {
		//判断文件的结尾名
		if strings.HasSuffix(fileInfoList[i].Name(), ".xml") {
			log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())
			//压缩文件
			lz77zip.Lz77zip(fileInfoList[i].Name())
			//解析文件
			//		解析文件  获取数据
			sendStru := ParsingXMLFiles(fileInfoList[i].Name())

			//发送
			client.Sendxml(&sendStru, conn)

			buf := make([]byte, 1024)
			n, err2 := conn.Read(buf)
			if err2 != nil {
				log.Println("conn.Read err = ", err2)
				return
			}
			log.Println(string(buf[:n]))
		}
	}
}

//解析xml文件
func ParsingXMLFiles(fname string) types.SendStru {
	var sendStru types.SendStru
	//1、获取消息包序号Massageid
	fnstr := strings.Split(fname, "_")
	idstr := strings.Split(fnstr[2], ".")
	sendStru.Massageid = idstr[0]

	//2、获取文件大小
	lengthstr := fmt.Sprintf("%06d", service.GetFileSize(fname))

	sendStru.Xml_length = lengthstr
	log.Println("发送文件大小", lengthstr)

	//3、获取文件md5  "JZ_3301_00000000000000100094.xml.lz77"
	sendStru.Md5_str = common.GetFileMd5(fname)
	if sendStru.Md5_str != "" {
		log.Println("文件md5：", sendStru.Md5_str)
	} else {
		log.Println("获取文件md5 error ")
	}

	//4、获得文件名
	sendStru.Xml_msgName = fname
	//
	log.Println("报文信息：", sendStru)

	////5、更新数据    根据 包号 更新原始交易消息包的【发送状态   发送中】
	//Mid, _ := strconv.Atoi(sendStru.Massageid)
	//err := storage.UpdateYuansjyxx(int64(Mid))
	//if err != nil {
	//	log.Println("根据 包号 更新原始交易消息包的发送状态  error: ", err)
	//
	//}

	////移动已压缩的xml文件
	//xmls := "../generatexml/" + fname
	//xmldes := "../compressed_xml/" + fname
	//client.MoveFile(xmls, xmldes)
	return sendStru
}
