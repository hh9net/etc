package service

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	"CenterSettlement-go/lz77zip"
	"CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//线程2 发送数据包
func HandleSendXml() {
	//从文件夹sendzipxml中扫描打包文件（判断这个文件夹下面有没有文件）
	tiker := time.NewTicker(time.Second * 15)
	for {
		log.Println(common.DateTimeFormat(<-tiker.C), "执行线程2 发送原始交易包")

		pwd := "./generatexml/" //先压缩后发送    go run main.go
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("该文件夹下有文件的数量 %d：", len(fileInfoList))

		for i := range fileInfoList {
			//判断文件的结尾名
			if strings.HasSuffix(fileInfoList[i].Name(), ".xml") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())
				//压缩文件   后面使用并发压缩与并发解压缩 可以提高速度近10倍
				zerr := lz77zip.ZipLz77(fileInfoList[i].Name())
				if zerr != nil {
					log.Println("发送文件时 压缩xml文件失败", zerr)
					return
				}
				//压缩成功 移动xml文件
				if zerr == nil {
					//  CenterSettlement-go  //goland
					s := "./generatexml/" + fileInfoList[i].Name()
					des := "./compressedXml/" + fileInfoList[i].Name()
					merr := client.MoveFile(s, des)
					if merr != nil {
						log.Println("client.MoveFile err : ", merr)
						return
					}

					//移动成功 应该更新原始交易包  文件路径

				}

				// 解析文件  获取TCP报文发送数据
				sendStru := ParsingXMLFiles(fileInfoList[i].Name())
				log.Println("连接联网中心服务器")
				//连接联网中心服务器
				address := conf.AddressConfigInit()
				Address := address.AddressIp + ":" + address.AddressPort
				log.Println("联网中心服务器 Address:", Address)
				conn, derr := net.Dial("tcp", Address)
				if derr != nil {
					log.Println("Dial 失败", derr)
					return
				}
				if conn != nil {
					log.Println("Dial 联网中心服务器 成功")
				}
				//发送
				client.Sendxml(sendStru, &conn)
				//读取即时应答 (buf[:n])
				buf := make([]byte, 1024)
				n, err2 := conn.Read(buf)
				if err2 != nil {
					log.Println("conn.Read err = ", err2)
					return
				}

				//对联网中心的接收应答处理
				resperr := ImmediateResponseProcessing(string(buf[:n]), fileInfoList[i].Name(), sendStru, &conn)
				if resperr != nil {
					return
				}
				conn.Close()
			}
		}
	}
}

//解析xml文件 fname为原始交易消息包xml文件的名字
func ParsingXMLFiles(fname string) *types.SendStru {
	var sendStru types.SendStru
	//1、获取消息包序号Massageid
	fnstr := strings.Split(fname, "_")
	idstr := strings.Split(fnstr[2], ".")
	sendStru.Massageid = idstr[0] //20位

	//2、获取压缩文件大小
	lengthstr := fmt.Sprintf("%06d", common.GetFileSize(fname))
	sendStru.Xml_length = lengthstr //
	log.Println("要发送压缩文件的大小（6位数）：", lengthstr)

	//3、获取xml文件md5 字母全大写【联网中心大小写忽略】
	sendStru.Md5_str = common.GetFileMd5(fname) //fname 是 xml
	if sendStru.Md5_str != "" {
		log.Println(" 文件大小写忽略 的 md5 为 ：", sendStru.Md5_str)
	} else {
		log.Fatal("获取文件md5 error ")
		return nil
	}

	//4、获得xml文件名
	sendStru.Xml_msgName = fname

	//5、更新数据    根据 包号 更新原始交易消息包的【发送状态   发送中】
	Mid, _ := strconv.Atoi(sendStru.Massageid)
	err := storage.UpdateYuansjyxx(int64(Mid))
	if err != nil {
		log.Println("根据 包号 更新原始交易消息包的发送状态  error: ", err)
	}
	return &sendStru
}

//即使应答处理
func ImmediateResponseProcessing(str string, name string, sendStru *types.SendStru, conn *net.Conn) error {
	log.Println("接收即使应答为：", str)
	b := []byte(str)
	if string(b[len(str)-1:]) == "1" {

		log.Println("收到联网中心的接收成功的即时应答")
		//TCP 记录联网中心的即时应答

		has, serr, reqcount := storage.GetTcpReqRecord(string(b[:len(str)-1]))
		if serr != nil {
			log.Println("查询TCP接收记录错误")
		}
		var resRecord storage.BJsTcpydjl

		if has == false {
			log.Println("此msgid属于第一次接收")
			resRecord.FVcXiaoxxh = string(b[:len(str)-1])
			resRecord.FNbFasz = 2
			resRecord.FNbChongfcs = 0
			resRecord.FDtZuixsj = time.Now()
			err2 := storage.TcpResponseRecordInsert(resRecord)
			if err2 != nil {
				log.Println("storage.TcpReqRecordInsert error:", err2)
			}
		}
		if has == true {
			log.Println("此msgid 已经接收")
			resRecord.FVcXiaoxxh = string(b[:len(str)-1])
			resRecord.FDtZuixsj = time.Now()
			resRecord.FNbChongfcs = reqcount + 1
			err3 := storage.TcpReqRecordUpdate(resRecord)
			if err3 != nil {
				log.Println("storage.TcpReqRecordUpdate error:", err3)
				return err3
			}
		}
		//发送成功后 mv 压缩的zip文件 到另一个文件夹中sendxmlsucceed        CenterSettlement-go

		s := "./sendzipxml/" + name + ".lz77"
		des := "./sendxmlsucceed/" + name + ".lz77"
		mverr := client.MoveFile(s, des)
		if mverr != nil {
			return mverr
		}

	}
	if string(b[len(str)-1:]) == "0" {
		//
		log.Println("联网中心的接收失败")
		log.Println("原始交易包发送失败,触发重发机制")
		//	发送失败 触发重发机制
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 5)
			client.Sendxml(sendStru, conn)
			log.Printf("原始交易包发送失败,触发重发机制,重发第%d次", 1+i)
		}
	}
	return nil
}
