package service

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

//线程4 处理数据包  定期扫描 接收联网的接收数据的文件夹 receivexml，如果有文件就解压， 解压后分析包。
func AnalyzeDataPakage() {

	//定期检查文件夹receive    解压后
	tiker := time.NewTicker(time.Second * 3)
	for {
		log.Println("执行线程4", <-tiker.C)
		log.Println("现在执行线程4", common.DateTimeNowFormat())

		//1、处理文件解压，解压至receivexml文件夹

		//2、处理文件解析
		ParseFile()
	}
}

//2、处理文件解析
func ParseFile() {
	//扫描receivexml 文件夹 读取文件信息
	//获取文件或目录相关信息
	pwd := "../receivexml/"
	//pwd := "CenterSettlement-go/receivexml/"
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
	if len(fileInfoList) == 0 {
		log.Println("该receivexml 文件夹下没有需要文件")
		return
	}
	for i := range fileInfoList {

		//判断文件的结尾名
		if strings.HasSuffix(fileInfoList[i].Name(), ".xml") {
			log.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名

			//解析xml文件
			//获取xml文件位置
			content, err := ioutil.ReadFile("../receivexml/" + fileInfoList[i].Name())
			if err != nil {
				log.Println("读文件位置错误信息：", err)
				return
			}

			//将文件转换为对象
			var result types.ReceiveMessage
			err = xml.Unmarshal(content, &result)
			if err != nil {
				log.Println("解析 receive文件夹中xml文件内容的错误信息：", err)
			}

			log.Println("result:", result.Header.MessageClass, result.Header.MessageType, result.Body.ContentType, result.Header.MessageId)
			////原始交易数据
			//if result.Header.MessageClass == 5 && result.Header.MessageType == 7 &&result.Body.ContentType==1{
			//	//
			//	return
			//}
			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 1 {
				//记账数据包
				//1、修改文件名字  2、移动文件
				src := "../receivexml/" + fileInfoList[i].Name()
				des := "../keepAccountFile/" + "JZB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				frerr := common.FileRename(src, des)
				if frerr != nil {
					log.Println("记账数据包 修改文件名字错误：", frerr)
					return
				}

				//解析xml数据 把数据导入数据库
				Parsexml("../keepAccountFile/", "JZB"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
			}
			if result.Header.MessageClass == 5 && result.Header.MessageType == 7 && result.Body.ContentType == 2 {
				//争议数据包
				//1、修改文件名字  2、移动文件
				src := "../receivexml/" + fileInfoList[i].Name()
				des := "../disputeProcessFile/" + "ZYB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				frerr := common.FileRename(src, des)
				if frerr != nil {
					log.Println("记账数据包 修改文件名字错误：", frerr)
					return
				}
				//解析xml数据 把数据导入数据库
				Parsexml("../disputeProcessFile/", "ZYB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
			}
			//清分数据包
			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 2 {
				//1、修改文件名字  2、移动文件
				src := "../receivexml/" + fileInfoList[i].Name()
				des := "../clearling/" + "QFB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				frerr := common.FileRename(src, des)
				if frerr != nil {
					log.Println("记账数据包 修改文件名字错误：", frerr)
					return
				}
				//解析xml数据 把数据导入数据库
				Parsexml("../clearling/", "QFB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
			}
			//		退费数据包【不做】
			//if result.Header.MessageClass == 5 && result.Header.MessageType == 7 &&result.Body.ContentType==3{
			//
			//	return
			//}
		}
	}

}

//解析xml数据 把数据导入数据库
func Parsexml(filepath string, fname string) {

	fileInfoList, err := ioutil.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := range fileInfoList {
		//判断文件的结尾名
		if fileInfoList[i].Name() == fname {
			log.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名

			//解析xml文件
			//获取xml文件位置
			content, err := ioutil.ReadFile(filepath + fname)
			if err != nil {
				log.Println("读文件位置错误信息：", err)
				return
			}
			f := strings.Split(fname, "_")
			switch f[0] {
			//
			case "JZB":
				//将文件转换为对象
				var result types.KeepAccountMessage
				err = xml.Unmarshal(content, &result)
				if err != nil {
					log.Println("解析 KeepAccount 文件夹中xml文件内容的错误信息：", err)
				}
				log.Println("", result)

				//将数据存储数据库
				//新增记账包消息
				jzshuju := new(types.BJsJizclxx)
				//赋值
				jzshuju.FNbXiaoxxh = 123
				jzerr := storage.InsertMessageData(jzshuju)
				if jzerr != nil {
					log.Println("新增记账包消息错误 ：", jzerr)
				}

				//新增记账包明细
				jzshujuMX := new(types.BJsJizclmx)
				//赋值  调用函数
				jzshujuMX.FNbYuansjyxxxh = 12
				jzmxerr := storage.InsertMessageMXData(jzshujuMX)
				if jzmxerr != nil {
					log.Println("新增记账包明细 error ：", jzmxerr)
				}
				//更新结算数据为已记账

			case "ZYB":
				var result types.DisputeProcessMessage
				err = xml.Unmarshal(content, &result)
				if err != nil {
					log.Println("解析 DisputeProcess文件夹中xml文件内容的错误信息：", err)
				}
				log.Println("", result)

				//将数据存储数据库

				//新增争议包消息
				zyshuju := new(types.BJsZhengyclxx)
				//赋值
				zyshuju.FNbXiaoxxh = 123
				jzerr := storage.DisputeProcessXXInsert(zyshuju)
				if jzerr != nil {
					log.Println("新增争议包消息错误 ：", jzerr)
				}

				//新增争议包明细
				//新增争议包消息
				zyshujuMX := new(types.BJsZhengyjyclmx)
				//赋值
				zyshujuMX.FNbZhengyjyclxxxh = 123
				zyerr := storage.DisputeProcessMxInsert(zyshujuMX)
				if zyerr != nil {
					log.Println("新增争议包消息错误 ：", zyerr)
				}

				//更新结算数据

			case "QFB":
				var result types.ClearingMessage
				err = xml.Unmarshal(content, &result)
				if err != nil {
					log.Println("解析 Clearing 文件夹中xml文件内容的错误信息：", err)
				}
				log.Println("", result)

				//将数据存储数据库

				//新增清分包消息
				qfshuju := new(types.BJsQingftjxx)
				//赋值
				qfshuju.FNbXiaoxxh = 123
				qferr := storage.ClearingInsert(qfshuju)
				if qferr != nil {
					log.Println("新增清分包消息错误 ：", qferr)
				}

				//新增清分包明细
				qfshujumx := new(types.BJsQingftongjimx)
				//赋值
				qfshujumx.FNbQingftjxxxh = 123
				qfmxerr := storage.ClearingInsert(qfshuju)
				if qfmxerr != nil {
					log.Println("新增清分包消息错误 ：", qfmxerr)
				}

				//更新结算数据

			}

		}

	}
}
