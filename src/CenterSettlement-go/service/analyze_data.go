package service

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	"CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

//线程4 处理数据包  定期扫描 接收联网的接收数据的文件夹 receivexml，如果有文件就解压， 解压后分析包。
func AnalyzeDataPakage() {

	//定期检查文件夹receivexml
	tiker := time.NewTicker(time.Second * 10)
	for {
		log.Println("执行线程4 处理数据包")
		//1、处理文件解压，解压至receivexml文件夹 [已ok]

		//2、处理文件解析
		perr := ParseFile()
		if perr != nil {
			log.Println("ParseFile 失败 ", perr)
			return
		}

		log.Println("执行线程4 处理数据包", common.DateTimeFormat(<-tiker.C))

	}
}

//2、处理文件解析
func ParseFile() error {
	//扫描receivexml 文件夹 读取文件信息
	pwd := "./receivexml/"
	fileList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 1 {
		log.Println("该receivexml 文件夹下没有需要解析的xml文件")
		return nil
	}
	for i := range fileList {
		//判断文件的结尾名
		if strings.HasSuffix(fileList[i].Name(), ".xml") {
			log.Println("该receivexml 文件夹下需要解析的xml文件名字为:", fileList[i].Name())
			content, err := ioutil.ReadFile("./receivexml/" + fileList[i].Name())
			if err != nil {
				log.Println("读文件位置错误信息：", err)
				return err
			}

			//将xml文件转换为对象
			var result types.ReceiveMessage
			uerr := xml.Unmarshal(content, &result)
			if uerr != nil {
				log.Println("解析 receivexml文件夹中xml文件内容时，错误信息为：", uerr)
			}

			log.Println("result:", result.Header.MessageClass, result.Header.MessageType, result.Body.ContentType, result.Header.MessageId)

			//记账数据包
			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 1 {
				if pjzerr := ParseKeepAccountFile(result, fileList[i].Name()); pjzerr != nil {
					return pjzerr
				}
			}

			if result.Header.MessageClass == 5 && result.Header.MessageType == 7 && result.Body.ContentType == 2 {
				//争议数据包
				if pzyerr := ParseDisputeFile(result, fileList[i].Name()); pzyerr != nil {
					return pzyerr
				}
			}

			//清分数据包
			if result.Header.MessageClass == 5 && result.Header.MessageType == 5 && result.Body.ContentType == 2 {
				if pqferr := ParseClearlingFile(result, fileList[i].Name()); pqferr != nil {
					return pqferr
				}
			}

			//原始数据应答包
			if result.Header.MessageClass == 6 && result.Header.MessageType == 7 && result.Body.ContentType == 1 {
				if preqerr := ParseRespFile(result, fileList[i].Name()); preqerr != nil {
					return preqerr
				}
			}

			//		//退费数据包【不做】
			//if result.Header.MessageClass == 5 && result.Header.MessageType == 7 &&result.Body.ContentType==3{
			//
			//	return
			//}
		}

	}

	//log.Println()
	return nil
}

//记账包的解析
func ParseKeepAccountFile(result types.ReceiveMessage, fname string) error {
	//记账包的确认应答
	gerr, filename, resmsg := GenerateRespMessage("jz", result)
	if gerr != nil {
		return gerr
	}
	log.Printf("记账包的确认应答包的名字", filename)
	//新增记账包的应答包记录
	inerr := storage.InsertResMessageData(resmsg)
	if inerr != nil {
		return inerr
	}

	var des string
	//记账数据包
	//1、修改文件名字  2、移动文件
	if result.Body.DisputedCount == 0 {
		des = "./keepAccountFile/" + "JZB-ok_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	}
	if result.Body.DisputedCount > 0 {
		des = "./keepAccountFile/" + "JZB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	}
	src := "./receivexml/" + fname
	jzfrerr := common.FileRename(src, des)
	if jzfrerr != nil {
		log.Println("记账数据包 修改文件名字错误：", jzfrerr)
		return jzfrerr
	}
	//log.Println("keepAccount result:", result)

	//解析xml数据 把数据导入数据库
	if result.Body.DisputedCount == 0 {
		jzpxerr := Parsexml("./keepAccountFile/", "JZB-ok"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
		if jzpxerr != nil {
			log.Println("记账数据包 解析xml数据 把数据导入数据库 时 错误 ：", jzpxerr)
			return jzpxerr
		}
	} else {
		jzpxerr := Parsexml("./keepAccountFile/", "JZB"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
		if jzpxerr != nil {
			log.Println("记账数据包 解析xml数据 把数据导入数据库 时 错误 ：", jzpxerr)
			return jzpxerr
		}
	}

	return nil
}

//争议包的解析
func ParseDisputeFile(result types.ReceiveMessage, fname string) error {
	//争议包的确认应答
	gerr, filename, resmsg := GenerateRespMessage("zy", result)
	if gerr != nil {
		return gerr
	}

	log.Printf("争议处理包的确认应答包的名字", filename)

	//新增记账包的应答包记录
	inerr := storage.DisputeProcessResInsert(resmsg)
	if inerr != nil {
		return inerr
	}

	//1、修改文件名字  2、移动文件
	src := "./receivexml/" + fname
	des := "./disputeProcessFile/" + "ZYB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	zyfrerr := common.FileRename(src, des)
	if zyfrerr != nil {
		log.Println("争议数据包 修改文件名字错误：", zyfrerr)
		return zyfrerr
	}
	//解析xml数据 把数据导入数据库
	zypxerr := Parsexml("./disputeProcessFile/", "ZYB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
	if zypxerr != nil {
		log.Println("争议数据包 解析xml数据 把数据导入数据库 时 错误 ：", zypxerr)
		return zypxerr
	}
	return nil
}

//清分包的解析
func ParseClearlingFile(result types.ReceiveMessage, fname string) error {
	//清分包的确认应答
	gerr, filename, resmsg := GenerateRespMessage("qf", result)
	if gerr != nil {
		return gerr
	}
	log.Printf("清分处理包的确认应答包的名字", filename)

	//新增清分包的应答包记录
	inerr := storage.ClearingYDInsert(resmsg)
	if inerr != nil {
		return inerr
	}

	var des string
	//1、修改文件名字  2、移动文件
	if result.Body.List.FileCount == 0 {
		des = "./clearing/" + "QFB-ok_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	}
	if result.Body.List.FileCount > 0 {
		des = "./clearing/" + "QFB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	}
	src := "./receivexml/" + fname
	qffrerr := common.FileRename(src, des)
	if qffrerr != nil {
		log.Println("清分数据包 修改文件名字错误：", qffrerr)
		return qffrerr
	}
	if result.Body.List.FileCount == 0 {
		//解析xml数据 把数据导入数据库
		qfpxerr := Parsexml("./clearing/", "QFB-ok_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
		if qfpxerr != nil {
			log.Println("清分数据包 解析xml数据 把数据导入数据库 时 错误 ：", qfpxerr)
			return qfpxerr
		}
	} else {
		//解析xml数据 把数据导入数据库
		qfpxerr := Parsexml("./clearing/", "QFB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
		if qfpxerr != nil {
			log.Println("清分数据包 解析xml数据 把数据导入数据库 时 错误 ：", qfpxerr)
			return qfpxerr
		}
	}

	return nil
}

//解析原始记录应答包
func ParseRespFile(result types.ReceiveMessage, fname string) error {
	//1、修改文件名字  2、移动文件
	src := "./receivexml/" + fname
	des := "./respfile/" + "YSYDB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml"
	ydfrerr := common.FileRename(src, des)
	if ydfrerr != nil {
		log.Println("原始记录应答包 修改文件名字错误：", ydfrerr)
		return ydfrerr
	}
	//解析xml数据 把数据导入数据库
	ydpxerr := Parsexml("./respfile/", "YSYDB_"+fmt.Sprintf("%020d", result.Header.MessageId)+".xml")
	if ydpxerr != nil {
		log.Println("原始数据应答包 解析xml数据 把数据导入数据库 时 错误 ：", ydpxerr)
		return ydpxerr
	}
	return nil
}

//到对应的文件夹下，使用对应的xml结构体，获取数据，插入数据
//解析xml数据 把数据导入数据库
func Parsexml(filePath string, fname string) error {
	log.Println("解析xml数据：filePath+fname， 把数据导入数据库", filePath+fname)
	fileInfoList, rerr := ioutil.ReadDir(filePath)
	if rerr != nil {
		log.Fatal(rerr)
		return rerr
	}
	for i := range fileInfoList {
		//判断文件的结尾名
		if fileInfoList[i].Name() == fname {
			log.Println("消息包要入数据库的xml文件名：", fileInfoList[i].Name())

			//解析xml文件 获取xml文件位置
			content, rderr := ioutil.ReadFile(filePath + fname)
			if rderr != nil {
				log.Println("读文件位置错误信息：", rderr)
				return rderr
			}
			f := strings.Split(fname, "_")
			switch f[0] {
			//记账包
			case "JZB-ok":
				//将文件转换为对象
				var result types.KeepAccountokMessage
				umerr := xml.Unmarshal(content, &result)
				if umerr != nil {
					log.Println("解析 KeepAccount 文件夹中xml文件内容的错误信息：", umerr)
					return umerr
				}
				log.Println("解析 KeepAccount 文件夹中xml文件内容result为：", result)

				//将数据存储数据库
				//新增记账包消息
				seterr := KeepAccountMessageInsert(result)
				if seterr != nil {
					return seterr
				}
				//新增 记账包明细
				setmxerr := KeepAccountMessageMxInsert(result)
				if setmxerr != nil {
					return setmxerr
				}
				//更新结算数据 记账状态为 已记账
				uperr := KeepAccountUpdate(result)
				if uperr != nil {
					return uperr
				}
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
				//新增记账包消息【有争议】
				disseterr := KeepAccountDisputeMessageInsert(result)
				if disseterr != nil {
					return disseterr
				}
				//新增 记账包明细【有争议】
				dissetmxerr := KeepAccountDisputeMessageMxInsert(result)
				if dissetmxerr != nil {
					return dissetmxerr
				}
				////更新结算数据 记账状态为 没有争议的：已记账【有争议】
				disuperr := KeepAccountDisputeUpdate(result)
				if disuperr != nil {
					return disuperr
				}

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
				diserr := DisputeProcessMessageInsert(result)
				if diserr != nil {
					log.Println("新增争议包消息错误 ：", diserr)
				}
				//新增争议包明细
				dismxerr := DisputeProcessMessageMxInsert(result)
				if dismxerr != nil {
					log.Println("新增争议包消息明细错误 ：", dismxerr)
					return dismxerr
				}

			//更新结算数据  此处需要确认

			case "QFB-ok":
				var result types.ClearingokMessage
				qfumerr := xml.Unmarshal(content, &result)
				if qfumerr != nil {
					log.Println("解析 Clearing 文件夹中xml文件内容的错误信息：", qfumerr)
				}
				log.Println("清分包数据：", result)

				//将数据存储数据库
				//新增清分包消息
				cerr := ClearingokMessageInsert(result)
				if cerr != nil {
					return cerr
				}

				//新增清分包明细
				cmxerr := ClearingokMessageMxInsert(result)
				if cmxerr != nil {
					return cmxerr
				}

			//更新结算数据  更新已数据清分

			case "QFB":
				//有可疑数据 存入fileid
				var result types.ClearingMessage
				qfumerr := xml.Unmarshal(content, &result)
				if qfumerr != nil {
					log.Println("解析 Clearing 文件夹中xml文件内容的错误信息：", qfumerr)
				}
				log.Println("清分包数据：", result)

				//将数据存储数据库
				//新增清分包消息
				cerr := ClearingMessageInsert(result)
				if cerr != nil {
					return cerr
				}

				//新增清分包明细
				cmxerr := ClearingMessageMxInsert(result)
				if cmxerr != nil {
					return cmxerr
				}

			//更新结算数据  更新已数据清分

			//原始交易应答包
			case "YSYDB":
				var result types.RespMessage
				qfumerr := xml.Unmarshal(content, &result)
				if qfumerr != nil {
					log.Println("解析 respfile 文件夹中xml文件内容的错误信息：", qfumerr)
				}
				log.Println("原始交易应答包数据：", result)

				//将数据存储数据库
				//新增原始交易应答包消息

				resperr := RespMessageInsert(result)
				if resperr != nil {
					return resperr
				}
			}
		}
		//这里加一个跳转

	}
	return nil
}

//封装函数处理数据入库的问题

//新增记账包消息
func KeepAccountMessageInsert(result types.KeepAccountokMessage) error {
	jzshuju := new(types.BJsJizclxx)
	//赋值
	jzshuju.FVcBanbh = result.Header.Version        //F_VC_BANBH	版本号	VARCHAR(32)
	jzshuju.FNbXiaoxlb = result.Header.MessageClass //F_NB_XIAOXLB	消息类别	INT
	jzshuju.FNbXiaoxlx = result.Header.MessageType  //F_NB_XIAOXLX	消息类型	INT
	jzshuju.FVcFaszid = result.Header.SenderId      //F_VC_FASZID	发送者ID	VARCHAR(32)
	jzshuju.FVcJieszid = result.Header.ReceiverId   //F_VC_JIESZID	接收者ID	VARCHAR(32)
	jzshuju.FNbXiaoxxh = result.Header.MessageId    //F_NB_XIAOXXH	消息序号	BIGINT
	jzshuju.FDtJiessj = common.DateTimeNowFormat()  //F_DT_JIESSJ	接收时间	DATETIME
	jzshuju.FNbYuansjyxxxh = result.Body.MessageId  //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	jzshuju.FNbJilsl = result.Body.Count            //F_NB_JILSL	记录数量	INT

	amount, _ := strconv.Atoi(result.Body.Amount)
	jzshuju.FNbZongje = amount                      //F_NB_ZONGJE	总金额	INT
	jzshuju.FNbZhengysl = result.Body.DisputedCount //F_NB_ZHENGYSL	争议数量	INT
	jzshuju.FNbZhixjg = 1                           //F_NB_ZHIXJG	执行结果	INT 1:消息已正常接收
	t := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))
	jzshuju.FDtChulsj = t //F_DT_CHULSJ	处理时间	DATETIME

	jzshuju.FVcXiaoxwjlj = "keepAccountFile/" + "JZB-ok_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
	log.Println("新增记账包消息", jzshuju)
	jzerr := storage.InsertMessageData(jzshuju)
	if jzerr != nil {
		log.Println("新增记账包消息错误 ：", jzerr)
		return jzerr
	}
	log.Println("新增记账包消息成功")
	return nil
}

//新增 记账包明细
func KeepAccountMessageMxInsert(result types.KeepAccountokMessage) error {

	jzshujuMX := new(types.BJsJizclmx)
	//全部记账，没有争议数据
	//查询原始交易包mx 原始交易序号 包内序号
	yuanshiMx := storage.QueryYuanshiMx(result.Body.MessageId)
	if len(*yuanshiMx) == 0 {
		return errors.New("结算数据表没有此包中的数据")
	}

	for _, mx := range *yuanshiMx {
		//赋值  调用函数
		jzshujuMX.FNbYuansjyxxxh = result.Body.MessageId //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
		jzshujuMX.FNbBaonxh = mx.FNbBaonxh               //F_NB_BAONXH	包内序号	INT
		jzshujuMX.FNbChuljg = 0                          //F_NB_CHULJG	处理结果	INT

		jzmxerr := storage.InsertMessageMXData(jzshujuMX)
		if jzmxerr != nil {
			log.Println("新增记账包明细 error ：", jzmxerr)
			return jzmxerr
		}
	}
	log.Println("新增记账包明细 成功")
	return nil
}

//更新结算数据 记账状态为 已记账
func KeepAccountUpdate(result types.KeepAccountokMessage) error {
	//全部记账，没有争议数据
	//查询原始交易包mx 原始交易序号 包内序号
	yuanshiMx := storage.QueryYuanshiMx(result.Body.MessageId)

	for _, mx := range *yuanshiMx {

		//更新结算数据 记账状态为 已记账
		uperr := storage.KeepAccountUpdate(result.Body.MessageId, mx.FNbBaonxh, 1)
		if uperr != nil {
			log.Println("更新结算数据 记账状态为 已记账 error ：", uperr)
			return uperr
		} else {
			log.Println("更新结算数据 记账状态为 已记账成功")
		}
	}
	return nil
}

//新增记账包消息 有争议
func KeepAccountDisputeMessageInsert(result types.KeepAccountMessage) error {
	jzshuju := new(types.BJsJizclxx)
	//赋值
	jzshuju.FVcBanbh = result.Header.Version        //F_VC_BANBH	版本号	VARCHAR(32)
	jzshuju.FNbXiaoxlb = result.Header.MessageClass //F_NB_XIAOXLB	消息类别	INT
	jzshuju.FNbXiaoxlx = result.Header.MessageType  //F_NB_XIAOXLX	消息类型	INT
	jzshuju.FVcFaszid = result.Header.SenderId      //F_VC_FASZID	发送者ID	VARCHAR(32)
	jzshuju.FVcJieszid = result.Header.ReceiverId   //F_VC_JIESZID	接收者ID	VARCHAR(32)
	jzshuju.FNbXiaoxxh = result.Header.MessageId    //F_NB_XIAOXXH	消息序号	BIGINT
	jzshuju.FDtJiessj = common.DateTimeNowFormat()  //F_DT_JIESSJ	接收时间	DATETIME
	jzshuju.FNbYuansjyxxxh = result.Body.MessageId  //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
	jzshuju.FNbJilsl = result.Body.Count            //F_NB_JILSL	记录数量	INT

	str := strings.Split(result.Body.Amount, ".")
	je, _ := strconv.Atoi(str[0] + str[1])

	jzshuju.FNbZongje = je                          //F_NB_ZONGJE	总金额	INT
	jzshuju.FNbZhengysl = result.Body.DisputedCount //F_NB_ZHENGYSL	争议数量	INT
	jzshuju.FNbZhixjg = 1                           //F_NB_ZHIXJG	执行结果	INT 1:消息已正常接收

	t := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))
	jzshuju.FDtChulsj = t //F_DT_CHULSJ	处理时间	DATETIME

	jzshuju.FVcXiaoxwjlj = "keepAccountFile/" + "JZB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)
	log.Println(jzshuju)
	jzerr := storage.InsertMessageData(jzshuju)
	if jzerr != nil {
		log.Println("新增记账包消息错误 ：", jzerr)
		return jzerr
	}
	log.Println("新增记账包消息成功")
	return nil
}

//新增 记账包明细 有争议
func KeepAccountDisputeMessageMxInsert(result types.KeepAccountMessage) error {
	jzshujuMX := new(types.BJsJizclmx)
	//查询原始交易包mx 原始交易序号 包内序号
	yuanshiMx := storage.QueryYuanshiMx(result.Body.MessageId)
	if len(*yuanshiMx) == 0 {
		return errors.New("结算数据表没有此包中的数据")
	}

	jzshujuMX.FNbYuansjyxxxh = result.Body.MessageId //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT

	for _, mx := range *yuanshiMx {
		for _, dis := range result.Body.DisputedRecord {
			//	把有争议的
			if mx.FNbBaonxh == dis.TransId {
				jzshujuMX.FNbBaonxh = mx.FNbBaonxh //F_NB_BAONXH	包内序号	INT
				jzshujuMX.FNbChuljg = dis.Result   //F_NB_CHULJG	处理结果	INT

				//正常的处理结果为0
			} else {
				jzshujuMX.FNbBaonxh = mx.FNbBaonxh //F_NB_BAONXH	包内序号	INT
				jzshujuMX.FNbChuljg = 0            //F_NB_CHULJG	处理结果	INT  正常记账
			}
		}

		jzmxerr := storage.InsertMessageMXData(jzshujuMX)
		if jzmxerr != nil {
			log.Println("新增记账包明细 error ：", jzmxerr)
			return jzmxerr
		}
	}
	log.Println("新增有争议的记账包明细 成功")
	return nil
}

//更新结算数据 记账状态为 已记账 有争议
func KeepAccountDisputeUpdate(result types.KeepAccountMessage) error {
	//全部记账，没有争议数据
	//查询原始交易包mx 原始交易序号 包内序号
	yuanshiMx := storage.QueryYuanshiMx(result.Body.MessageId)
	for _, mx := range *yuanshiMx {

		for _, dis := range result.Body.DisputedRecord {
			//	把有争议的 记账状态为 争议数据
			if mx.FNbBaonxh == dis.TransId {
				//更新结算数据 记账状态为 已记账  1：已记账、2：争议数据'
				uperr := storage.KeepAccountUpdate(result.Body.MessageId, mx.FNbBaonxh, 2)
				if uperr != nil {
					log.Println("更新结算数据 记账状态 为 争议数据 error ：", uperr)
					return uperr
				} else {
					log.Println("更新结算数据 记账状态 为 争议数据 成功")
				}

				//正常的  记账状态为 已记账
			} else {
				//更新结算数据 记账状态为 已记账  1：已记账、2：争议数据
				uperr := storage.KeepAccountUpdate(result.Body.MessageId, mx.FNbBaonxh, 1)
				if uperr != nil {
					log.Println("更新结算数据 记账状态为 已记账 error ：", uperr)
					return uperr
				} else {
					log.Println("更新结算数据 记账状态为 已记账成功")
				}
			}
		}
	}
	return nil
}

//新增 争议数据包
func DisputeProcessMessageInsert(result types.DisputeProcessMessage) error {
	zyshuju := new(types.BJsZhengyjyclxx)
	//赋值
	zyshuju.FVcBanbh = result.Header.Version             //F_VC_BANBH	版本号	VARCHAR(32)
	zyshuju.FNbXiaoxlb = result.Header.MessageClass      //F_NB_XIAOXLB	消息类别	INT
	zyshuju.FNbXiaoxlx = result.Header.MessageType       //F_NB_XIAOXLX	消息类型	INT
	zyshuju.FVcFaszid = result.Header.SenderId           //F_VC_FASZID	发送者ID	VARCHAR(32)
	zyshuju.FVcJieszid = result.Header.ReceiverId        //F_VC_JIESZID	接收者ID	VARCHAR(32)
	zyshuju.FNbXiaoxxh = result.Header.MessageId         //F_NB_XIAOXXH	消息序号	BIGINT
	zyshuju.FDtJiessj = time.Now()                       //F_DT_JIESSJ	接收时间	DATETIME
	zyshuju.FVcQingffid = result.Body.ClearingOperatorId //F_VC_QINGFFID	清分方ID	VARCHAR(32)
	zyshuju.FVcLianwzxid = result.Body.IssuerId          //F_VC_LIANWZXID	联网中心ID	VARCHAR(32)
	zyshuju.FVcFaxfid = result.Body.ServiceProviderId    //F_VC_FAXFID	发行方ID	VARCHAR(32)
	zyshuju.FVcZhengyjgwjid = result.Body.FileId         //F_VC_ZHENGYJGWJID	争议结果文件ID	INT

	t := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))
	zyshuju.FDtZhengyclsj = t //F_DT_ZHENGYCLSJ	争议处理时间	DATETIME

	zyshuju.FNbZhengysl = result.Body.Count //F_NB_ZHENGYSL	争议数量	INT

	str := strings.Split(result.Body.Amount, ".")
	amount, _ := strconv.Atoi(str[0] + str[1])
	zyshuju.FNbQuerxyjzdzje = amount                                                                               //F_NB_QUERXYJZDZJE	确认需要记账的总金额	INT
	zyshuju.FVcXiaoxwjlj = "disputeProcessFile/" + "ZYB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)

	zyerr := storage.DisputeProcessXXInsert(zyshuju)
	if zyerr != nil {
		log.Println("新增争议包消息错误 ：", zyerr)
		return zyerr
	}

	return nil
}

//新增 争议数据包明细
func DisputeProcessMessageMxInsert(result types.DisputeProcessMessage) error {
	zyshujuMX := new(types.BJsZhengyjyclmx)
	//赋值
	zyshujuMX.FNbZhengyjyclxxxh = result.Header.MessageId //F_NB_ZHENGYJYCLXXXH	争议交易处理消息序号	BIGINT

	//根据包号、包内序号新增数据记录
	for i, mxList := range result.Body.MessageList {

		zyshujuMX.FNbFenzxh = i                     //F_NB_FENZXH	分组序号	INT  也就是在争议包里的序号
		zyshujuMX.FNbYuansjyxxxh = mxList.MessageId //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
		zyshujuMX.FNbZunjlsl = mxList.Count         //F_NB_ZUNJLSL	组内记录数量	INT

		str := strings.Split(mxList.Amount, ".")
		amount, _ := strconv.Atoi(str[0] + str[1])
		zyshujuMX.FNbZunjezh = amount                       //F_NB_ZUNJEZH	组内金额总和	INT
		zyshujuMX.FNbYuansbnxh = mxList.Transaction.TransId //F_NB_YUANSBNXH	原始包内序号	INT
		zyshujuMX.FNbChuljg = mxList.Transaction.Result     //F_NB_CHULJG	处理结果	INT
	}

	zyerr := storage.DisputeProcessMxInsert(zyshujuMX)
	if zyerr != nil {
		log.Println("新增争议包消息明细错误 ：", zyerr)
		return zyerr
	}
	log.Println("新增争议包消息明细 成功")
	return nil
}

//争议数据存储后，更新结算数据
func DisputeProcessUpdate(result types.DisputeProcessMessage) error {

	return nil
}

//新增清分包消息 无争议
func ClearingokMessageInsert(result types.ClearingokMessage) error {

	qfshuju := new(types.BJsQingftjxx)
	//赋值
	qfshuju.FVcBanbh = result.Header.Version
	qfshuju.FNbXiaoxlb = result.Header.MessageClass
	qfshuju.FNbXiaoxlx = result.Header.MessageType
	qfshuju.FVcFaszid = result.Header.SenderId
	qfshuju.FVcJieszid = result.Header.ReceiverId
	qfshuju.FNbXiaoxxh = result.Header.MessageId

	qfshuju.FDtJiessj = time.Now()                    //接收时间
	qfshuju.FVcQingfmbr = result.Body.ClearTargetDate //清分目标日

	str := strings.Split(result.Body.Amount, ".")
	jine, _ := strconv.Atoi(str[0] + str[1])
	qfshuju.FNbQingfzje = jine             //总金额
	qfshuju.FNbQingfsl = result.Body.Count //清分数量

	qfsj := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))
	qfshuju.FDtQingftjclsj = qfsj            //清分统计处理时间
	qfshuju.FNbYuansjysl = result.Body.Count //原始交易数量
	lists := len(result.Body.List.MessageId)
	qfshuju.FNbZhengycljgbsl = lists //争议处理结果包数量

	qfshuju.FDtChulsj = time.Now() //处理时间    入库时间

	qfshuju.FVcXiaoxwjlj = "clearing/" + "QFB-ok_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //消息文件路径

	qferr := storage.ClearingInsert(qfshuju)
	if qferr != nil {
		log.Println("新增清分包消息错误 ：", qferr)
	}

	return nil
}

//新增清分包消息明细 无争议
func ClearingokMessageMxInsert(result types.ClearingokMessage) error {

	qfmxshuju := new(types.BJsQingftjmx)
	//赋值
	qfmxshuju.FNbQingftjxxxh = result.Header.MessageId             //F_NB_QINGFTJXXXH	清分统计消息序号	BIGINT
	qfmxshuju.FVcTongxbzxxtid = result.Body.List.ServiceProviderId //F_VC_TONGXBZXXTID	通行宝中心系统ID	VARCHAR(32)
	//qfmxshuju.FNbZhengycljgwjid =result //F_NB_ZHENGYCLJGWJID	争议处理结果文件ID	INT

	for i, msgid := range result.Body.List.MessageId {
		qfmxshuju.FNbFenzxh = i + 1      //F_NB_FENZXH	分组序号	INT
		qfmxshuju.FNbYuansjyxxxh = msgid //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
		qfmxerr := storage.ClearingMXInsert(qfmxshuju)
		if qfmxerr != nil {
			log.Println("新增清分包消息错误 ：", qfmxerr)
		}
	}

	return nil
}

//新增清分包消息 有争议
func ClearingMessageInsert(result types.ClearingMessage) error {

	qfshuju := new(types.BJsQingftjxx)
	//赋值
	qfshuju.FVcBanbh = result.Header.Version
	qfshuju.FNbXiaoxlb = result.Header.MessageClass
	qfshuju.FNbXiaoxlx = result.Header.MessageType
	qfshuju.FVcFaszid = result.Header.SenderId
	qfshuju.FVcJieszid = result.Header.ReceiverId
	qfshuju.FNbXiaoxxh = result.Header.MessageId

	qfshuju.FDtJiessj = time.Now()                    //接收时间
	qfshuju.FVcQingfmbr = result.Body.ClearTargetDate //清分目标日

	str := strings.Split(result.Body.Amount, ".")
	jine, _ := strconv.Atoi(str[0] + str[1])
	qfshuju.FNbQingfzje = jine             //总金额
	qfshuju.FNbQingfsl = result.Body.Count //清分数量

	qfsj := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))
	qfshuju.FDtQingftjclsj = qfsj            //清分统计处理时间
	qfshuju.FNbYuansjysl = result.Body.Count //原始交易数量
	lists := len(result.Body.List.MessageId)
	qfshuju.FNbZhengycljgbsl = lists //争议处理结果包数量

	qfshuju.FDtChulsj = time.Now() //处理时间    入库时间

	qfshuju.FVcXiaoxwjlj = "clearing/" + "QFB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //消息文件路径

	qferr := storage.ClearingInsert(qfshuju)
	if qferr != nil {
		log.Println("新增清分包消息错误 ：", qferr)
	}

	return nil
}

//新增清分包消息明细  有争议
func ClearingMessageMxInsert(result types.ClearingMessage) error {

	qfmxshuju := new(types.BJsQingftjmx)
	//赋值
	qfmxshuju.FNbQingftjxxxh = result.Header.MessageId             //F_NB_QINGFTJXXXH	清分统计消息序号	BIGINT
	qfmxshuju.FVcTongxbzxxtid = result.Body.List.ServiceProviderId //F_VC_TONGXBZXXTID	通行宝中心系统ID	VARCHAR(32)
	//qfmxshuju.FNbZhengycljgwjid =result //F_NB_ZHENGYCLJGWJID	争议处理结果文件ID	INT

	for i, msgid := range result.Body.List.MessageId {
		qfmxshuju.FNbFenzxh = i + 1                //F_NB_FENZXH	分组序号	INT
		qfmxshuju.FNbYuansjyxxxh = msgid.MessageId //F_NB_YUANSJYXXXH	原始交易消息序号	BIGINT
		qfmxerr := storage.ClearingMXInsert(qfmxshuju)
		if qfmxerr != nil {
			log.Println("新增清分包消息错误 ：", qfmxerr)
		}
	}

	return nil
}

//新增原始交易应答包消息
func RespMessageInsert(result types.RespMessage) error {

	ydshuju := new(types.BJsYuansjyydxx)

	//赋值

	ydshuju.FVcBanbh = result.Header.Version        //F_VC_BANBH	版本号	VARCHAR(32)
	ydshuju.FNbXiaoxlb = result.Header.MessageClass //F_NB_XIAOXLB	消息类别	INT
	ydshuju.FNbXiaoxlx = result.Header.MessageType  //F_NB_XIAOXLX	消息类型	INT
	ydshuju.FVcFaszid = result.Header.SenderId      //F_VC_FASZID	发送者ID	VARCHAR(32)
	ydshuju.FVcJieszid = result.Header.ReceiverId   //F_VC_JIESZID	接收者ID	VARCHAR(32)
	ydshuju.FNbXiaoxxh = result.Header.MessageId    //F_NB_XIAOXXH	消息序号	BIGINT
	ydshuju.FNbQuerdxxxh = result.Body.MessageId    //F_NB_QUERDXXXH	确认的消息序号	BIGINT

	t := common.StrTimeTotime(common.DataTimeFormatHandle(result.Body.ProcessTime))

	ydshuju.FDtChulsj = t                             //F_DT_CHULSJ	处理时间	DATETIME
	ydshuju.FNbZhixjg = result.Body.Result            //F_NB_ZHIXJG	执行结果	INT
	ydshuju.FVcQingfmbr = result.Body.ClearTargetDate //F_VC_QINGFMBR	清分目标日	VARCHAR(32)

	ydshuju.FVcXiaoxwjlj = "respfile/" + "YSYDB_" + fmt.Sprintf("%020d", result.Header.MessageId) + ".xml" //F_VC_XIAOXWJLJ	消息文件路径	VARCHAR(512)

	qferr := storage.PackagingRespRecordInsert(ydshuju)
	if qferr != nil {
		log.Println("新增应答包消息错误 ：", qferr)
	}
	log.Println("新增应答包消息成功")
	return nil
}

//
//生成对应的确认消息应答包
func GenerateRespMessage(lx string, result types.ReceiveMessage) (error, string, *types.ResponseCTMessage) {
	// 记账包应答、争议包应答、清分包应答
	var msgclass, msgtype int
	var lxstr string
	Messageid := conf.GenerateYingdMessageId()
	Filenameid = fmt.Sprintf("%020d", Messageid)

	switch lx {
	case "jz":
		msgclass = 6
		msgtype = 5
		lxstr = "JZ_YDB"

	case "zy":
		msgclass = 6
		msgtype = 7
		lxstr = "ZY_YDB"
	case "qf":
		msgclass = 6
		msgtype = 5
		lxstr = "QF_YDB"
	}
	var ydmsgct *types.ResponseCTMessage
	if result.Body.ContentType > 0 {
		//获取数据
		ydmsgct = &types.ResponseCTMessage{
			Header: types.ResponseHeader{
				Version:      "00010000",
				MessageClass: msgclass,
				MessageType:  msgtype,
				SenderId:     "00000000000000FD",
				ReceiverId:   "0000000000000020",
				MessageId:    Messageid},
			Body: types.ResponseCTBody{
				ContentType: result.Body.ContentType,
				ProcessTime: common.DateTimeNowFormat(), //处理时间
				Result:      1,                          // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
				MessageId:   result.Header.MessageId,
			}}

		//生成xml
		//使用MarshalIndent函数，生成的XML格式有缩进
		outputxml, err := xml.MarshalIndent(ydmsgct, "  ", " ")
		if err != nil {
			log.Printf("error: %v\n", err)
			return err, "", nil
		}

		fw, f_werr := os.Create("./generatexml/" + lxstr + "_" + Filenameid + ".xml") //go run main.go
		if f_werr != nil {
			log.Fatal("Read:", f_werr)
			return f_werr, "", nil
		}
		//加入XML头
		headerBytes := []byte(xml.Header)
		//拼接XML头和实际XML内容
		xmlOutPutData := append(headerBytes, outputxml...)

		_, ferr := fw.Write((xmlOutPutData))
		if ferr != nil {
			log.Printf("Write xml file error: %v\n", ferr)
			return ferr, "", nil
		}
		//更新消息包信息
		fw.Close()

		return nil, lxstr + "_" + Filenameid + ".xml", ydmsgct
	}

	ydmsg := &types.ResponseMessage{
		Header: types.ResponseHeader{
			Version:      "00010000",
			MessageClass: msgclass,
			MessageType:  msgtype,
			SenderId:     "00000000000000FD",
			ReceiverId:   "0000000000000020",
			MessageId:    Messageid},
		Body: types.ResponseBody{
			ProcessTime: common.DateTimeNowFormat(), //处理时间
			Result:      1,                          // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
			MessageId:   result.Header.MessageId,
		}}

	//生成xml
	//使用MarshalIndent函数，生成的XML格式有缩进
	outputxml, err := xml.MarshalIndent(ydmsg, "  ", " ")
	if err != nil {
		log.Printf("error: %v\n", err)
		return err, "", nil
	}

	fw, f_werr := os.Create("./generatexml/" + lxstr + "_" + Filenameid + ".xml") //go run main.go
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
		return f_werr, "", nil
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, outputxml...)

	_, ferr := fw.Write((xmlOutPutData))
	if ferr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
		return ferr, "", nil
	}
	//更新消息包信息
	fw.Close()
	newydmsg := &types.ResponseCTMessage{
		Header: types.ResponseHeader{
			Version:      "00010000",
			MessageClass: msgclass,
			MessageType:  msgtype,
			SenderId:     "00000000000000FD",
			ReceiverId:   "0000000000000020",
			MessageId:    Messageid},
		Body: types.ResponseCTBody{
			ContentType: 0,
			ProcessTime: common.DateTimeNowFormat(), //处理时间
			Result:      1,                          // int8  1.消息已正常接收（用于Advice Response时含已接受建议）
			MessageId:   result.Header.MessageId,
		}}

	return nil, lxstr + "_" + Filenameid + ".xml", newydmsg
}
