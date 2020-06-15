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
	//定期检查文件夹receivexml

	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("执行线程4", <-tiker.C)
		log.Println("现在执行线程4", common.DateTimeNowFormat())

		//1、处理文件解压，解压至receivexml文件夹 [已ok]

		//2、处理文件解析
		perr := ParseFile()
		if perr != nil {
			log.Println("ParseFile失败 ")
			return
		}
	}
}

//2、处理文件解析
func ParseFile() error {
	//扫描receivexml 文件夹 读取文件信息
	//获取文件或目录相关信息
	pwd := "CenterSettlement-go/receivexml/"
	//pwd := "../receivexml/"
	fileList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 0 {
		log.Println("该receivexml 文件夹下没有需要解析的文件")
		return nil
	}
	for i := range fileList {
		log.Println("该receivexml 文件夹下需要解析的xml文件名字为:", fileList[i].Name()) //打印当前文件或目录下的文件或目录名
		//判断文件的结尾名
		if strings.HasSuffix(fileList[i].Name(), ".xml") {

			//content, err := ioutil.ReadFile("../receivexml/" + fileInfoList[i].Name())
			content, err := ioutil.ReadFile("CenterSettlement-go/receivexml/" + fileList[i].Name())
			if err != nil {
				log.Println("读文件位置错误信息：", err)
				return err
			}

			//将xml文件转换为对象
			var result types.ReceiveMessage
			err = xml.Unmarshal(content, &result)
			if err != nil {
				log.Println("解析 receivexml文件夹中xml文件内容时，错误信息为：", err)
			}

			log.Println("result:", result.Header.MessageClass, result.Header.MessageType, result.Body.ContentType, result.Header.MessageId)

			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 1 {
				//记账数据包
				//1、修改文件名字  2、移动文件
				src := "CenterSettlement-go/receivexml/" + fileList[i].Name()
				des := "CenterSettlement-go/keepAccountFile/" + "JZB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				jzfrerr := common.FileRename(src, des)
				if jzfrerr != nil {
					log.Println("记账数据包 修改文件名字错误：", jzfrerr)
					return jzfrerr
				}

				//解析xml数据 把数据导入数据库
				jzpxerr := Parsexml("CenterSettlement-go/keepAccountFile/", "JZB"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
				if jzpxerr != nil {
					log.Println("记账数据包 解析xml数据 把数据导入数据库：", jzpxerr)
					return jzpxerr
				}
			}
			if result.Header.MessageClass == 5 && result.Header.MessageType == 7 && result.Body.ContentType == 2 {
				//争议数据包
				//1、修改文件名字  2、移动文件
				src := "CenterSettlement-go/receivexml/" + fileList[i].Name()
				des := "CenterSettlement-go/disputeProcessFile/" + "ZYB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				zyfrerr := common.FileRename(src, des)
				if zyfrerr != nil {
					log.Println("争议数据包 修改文件名字错误：", zyfrerr)
					return zyfrerr
				}
				//解析xml数据 把数据导入数据库
				zypxerr := Parsexml("CenterSettlement-go/disputeProcessFile/", "ZYB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
				if zypxerr != nil {
					log.Println("记账数据包 解析xml数据 把数据导入数据库：", zypxerr)
					return zypxerr
				}
			}
			//清分数据包
			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 2 {
				//1、修改文件名字  2、移动文件
				src := "CenterSettlement-go/receivexml/" + fileList[i].Name()
				des := "CenterSettlement-go/clearlings/" + "QFB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				qffrerr := common.FileRename(src, des)
				if qffrerr != nil {
					log.Println("清分数据包 修改文件名字错误：", qffrerr)
					return qffrerr
				}
				//解析xml数据 把数据导入数据库
				qfpxerr := Parsexml("CenterSettlement-go/clearlings/", "QFB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
				if qfpxerr != nil {
					log.Println("记账数据包 解析xml数据 把数据导入数据库：", qfpxerr)
					return qfpxerr
				}
			}

			//原始数据应答包
			if result.Header.MessageClass == 6 && result.Header.MessageType == 7 && result.Body.ContentType == 1 {
				//1、修改文件名字  2、移动文件
				src := "CenterSettlement-go/receivexml/" + fileList[i].Name()
				des := "CenterSettlement-go/reqfile/" + "REQ_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
				ydfrerr := common.FileRename(src, des)
				if ydfrerr != nil {
					log.Println("清分数据包 修改文件名字错误：", ydfrerr)
					return ydfrerr
				}
				//解析xml数据 把数据导入数据库
				ydpxerr := Parsexml("CenterSettlement-go/reqfile/", "YDB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
				if ydpxerr != nil {
					log.Println("记账数据包 解析xml数据 把数据导入数据库：", ydpxerr)
					return ydpxerr
				}
			}

			//		退费数据包【不做】
			//if result.Header.MessageClass == 5 && result.Header.MessageType == 7 &&result.Body.ContentType==3{
			//
			//	return
			//}
		} else {
			log.Println(fileList[i].Name()) //不是xml文件
			return nil
		}
	}
	return nil
}

//到对应的文件夹下，使用对应的xml结构体，获取数据，插入数据
//解析xml数据 把数据导入数据库
func Parsexml(filePath string, fname string) error {

	fileInfoList, rerr := ioutil.ReadDir(filePath)
	if rerr != nil {
		log.Fatal(rerr)
		return rerr
	}

	for i := range fileInfoList {
		//判断文件的结尾名
		if fileInfoList[i].Name() == fname {
			log.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名

			//解析xml文件
			//获取xml文件位置
			content, rderr := ioutil.ReadFile(filePath + fname)
			if rderr != nil {
				log.Println("读文件位置错误信息：", rderr)
				return rderr
			}
			f := strings.Split(fname, "_")
			switch f[0] {
			//记账包
			case "JZB":
				//将文件转换为对象
				var result types.KeepAccountMessage
				umerr := xml.Unmarshal(content, &result)
				if umerr != nil {
					log.Println("解析 KeepAccount 文件夹中xml文件内容的错误信息：", umerr)
					return umerr
				}
				log.Println("解析 KeepAccount 文件夹中xml文件内容result为：", result)

				//将数据存储数据库
				//新增记账包消息
				jzshuju := new(types.BJsJizclxx)
				//赋值
				jzshuju.FNbXiaoxxh = result.Header.MessageId
				jzerr := storage.InsertMessageData(jzshuju)
				if jzerr != nil {
					log.Println("新增记账包消息错误 ：", jzerr)
					return jzerr
				}

				//新增记账包明细
				jzshujuMX := new(types.BJsJizclmx)
				//赋值  调用函数
				jzshujuMX.FNbYuansjyxxxh = 12
				jzmxerr := storage.InsertMessageMXData(jzshujuMX)
				if jzmxerr != nil {
					log.Println("新增记账包明细 error ：", jzmxerr)
					return jzmxerr
				}
				//更新结算数据为已记账

			//争议包
			case "ZYB":
				var result types.DisputeProcessMessage
				zyumerr := xml.Unmarshal(content, &result)
				if zyumerr != nil {
					log.Println("解析 DisputeProcess文件夹中xml文件内容的错误信息：", zyumerr)
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
				qfumerr := xml.Unmarshal(content, &result)
				if qfumerr != nil {
					log.Println("解析 Clearing 文件夹中xml文件内容的错误信息：", qfumerr)
				}
				log.Println("清分包数据：", result)

				//将数据存储数据库

				//新增清分包消息
				qfshuju := new(types.BJsQingftjxx)
				//赋值
				qfshuju.FVcBanbh = result.Header.Version
				qfshuju.FNbXiaoxlb = result.Header.MessageClass
				qfshuju.FNbXiaoxlx = result.Header.MessageType
				qfshuju.FVcFaszid = result.Header.SenderId
				qfshuju.FVcJieszid = result.Header.ReceiverId
				qfshuju.FNbXiaoxxh = result.Header.MessageId

				qfshuju.FDtJiessj = time.Now() //接收时间
				qfshuju.FVcQingfmbr = result.Body.ClearTargetDate

				//qfshuju.FNbQingfzje=result.Body.Amount  //单位为元 希望以string来存
				qfshuju.FNbQingfsl = result.Body.Count

				//qfshuju.FDtQingftjclsj=//清分统计处理时间
				//qfshuju.FNbYuansjysl=//原始交易数量
				//qfshuju.FNbZhengycljgbsl=//争议处理结果包数量
				//qfshuju.FDtChulsj=//处理时间
				qfshuju.FVcXiaoxwjlj = "" //消息文件路径
				log.Println("清分数据", qfshuju)
				qferr := storage.ClearingInsert(qfshuju)
				if qferr != nil {
					log.Println("新增清分包消息错误 ：", qferr)
				}

				//新增清分包明细
				qfshujumx := new(types.BJsQingftjmx)
				//赋值
				qfshujumx.FNbQingftjxxxh = result.Header.MessageId
				qfshujumx.FVcTongxbzxxtid = result.Header.ReceiverId //通行宝中心系统ID
				//qfshujumx.FNbFenzxh=//分组序号 入库者自行生成，可取数组下标
				//qfshujumx.FNbYuansjyxxxh=//原始交易消息序号
				//qfshujumx.FNbZhengycljgwjid=//争议处理结果文件ID

				qfmxerr := storage.ClearingMXInsert(qfshujumx)
				if qfmxerr != nil {
					log.Println("新增清分包消息错误 ：", qfmxerr)
				}

				//更新结算数据  更新已数据清分

			}

		}

	}
}
