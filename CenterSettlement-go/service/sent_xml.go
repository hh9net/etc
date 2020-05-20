package service

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/types"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//线程2 发送数据包
func HandleSendXml() {
	//从文件夹sendzipxml中扫描打包文件（判断这个文件夹下面有没有文件）
	//这里也可以从channel中读取线程1发送的数据
	tiker := time.NewTicker(time.Second * 2)
	for {
		fmt.Println(<-tiker.C)
		//准备数据   //read 后
		//扫描receive 文件夹 读取文件
		//获取文件或目录相关信息
		pwd := "../sendzipxml/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		for i := range fileInfoList {

			//判断文件的结尾名
			if strings.HasSuffix(fileInfoList[i].Name(), ".xml") {
				fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
			}

			//		解析文件
			//		应答数据包
			//		记账数据包
			//		争议数据包
			//		清分数据包
			//		退费数据包
		}

		//发送数据
		var sendStru types.SendStru
		sendStru.Massageid = "563462"
		sendStru.Md5_str = "3414"
		sendStru.Xml_length = "34"
		sendStru.Xml_msg = []byte("1233213213sadfdasfsdsd")
		//发送给联网中心
		ok := client.Sendxml1()
		if ok == "ok" {
			//调接口成功后 mv文件夹到另一个文件中
		}
		if ok == "no" {
			//	发送失败 触发重发机制
		}

	}

	//定时器定期扫描sendzipxml文件
	//	读取文件
	//	准备报文
	//	发送报文
	//	发送成功则mv消息包至sendsucceed
	//	发送失败 触发重发机制
}