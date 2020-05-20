package service

import (
	"fmt"
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
		fmt.Println("现在", time.Now().Format("2006-01-02 15:04:05"))
		<-tiker.C
		//扫描receive 文件夹 读取文件
		//获取文件或目录相关信息
		pwd := "../receive/"
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

	}
	//定时器定期扫描receive文件夹

}
