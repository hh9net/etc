package service

import (
	"CenterSettlement-go/types"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//线程4 处理数据包  定期扫描 接收联网的接收数据的文件夹 receive，如果有文件就解压， 解压后分析包。
func AnalyzeDataPakage() {

	//定期检查文件夹receive    解压后
	tiker := time.NewTicker(time.Second * 3)
	for {
		log.Println("执行线程4")
		log.Println("现在", time.Now().Format("2006-01-02 15:04:05"))
		<-tiker.C
		//扫描receive 文件夹 读取文件x信息
		//获取文件或目录相关信息
		//pwd := "CenterSettlement-go/receive/"
		pwd := "../receive/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		for i := range fileInfoList {

			//判断文件的结尾名
			if strings.HasSuffix(fileInfoList[i].Name(), ".xml") {
				log.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名

				//		解析xml文件
				//获取xml文件位置
				content, err := ioutil.ReadFile("../receive/" + fileInfoList[i].Name())
				if err != nil {
					log.Println("读文件位置错误信息：", err)
				}

				//将文件转换为对象
				var result types.Message
				err = xml.Unmarshal(content, &result)
				if err != nil {
					log.Println("读receive文件夹中xml文件的内容的错误信息：", err)
				}

				log.Println("result:", result)
				//原始交易数据
				if result.Header.MessageClass == 5 && result.Header.MessageType == 7 {
					//
					return
				}
				//		记账数据包
				if result.Header.MessageClass == 5 && result.Header.MessageType == 5 {
					//

				}
				//		争议数据包
				if result.Header.MessageClass == 5 && result.Header.MessageType == 7 {
					//

				}
				//		清分数据包
				if result.Header.MessageClass == 5 && result.Header.MessageType == 5 {
					//

				}

				//		退费数据包【不做】
			}

		}

	}

}
