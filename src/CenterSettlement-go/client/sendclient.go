package client

import (
	storage "CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func connsendFile(data []byte, fname string, connect net.Conn) error {

	connect.Write(data)
	//暂时通过客户端sleep 100毫秒解决粘包问题，还可以通过tcp重连解决，以后再用（包头+数据）封装数据包的方式解决
	time.Sleep(time.Millisecond * 100)

	path := "CenterSettlement-go/sendzipxml/" + fname + ".lz77"
	log.Println("path:=", path)
	file, oserr := os.Open(path)
	if oserr != nil {
		log.Println("os.Open error", oserr)
		return oserr
	}
	log.Println("文件conn发送开始")
	defer file.Close()
	total := 0
	buff := make([]byte, 1024*10)
	for {
		size, rerr := file.Read(buff)
		log.Println(size, rerr)
		if rerr != nil {
			log.Println("file.Read err:", rerr)
			break
		}
		if rerr == io.EOF {
			log.Println("文件读取完毕")
			log.Println(total)
			break
		}

		//data = append(data, buff[:size]...)
		//log.Println(string(data))
		_, werr := connect.Write(buff[:size])
		if werr != nil {
			log.Println("err", werr)
			return werr
		}
		total += size

	}
	log.Println("发送报文到联网中心成功")
	return nil
}

//发送
func Sendxml(sendStru *types.SendStru, conn net.Conn) {
	//把报文写给服务端
	data := []byte(sendStru.Massageid)
	length := []byte(sendStru.Xml_length)
	md5 := []byte(sendStru.Md5_str)
	data = append(data, length...)
	data = append(data, md5...)
	err := connsendFile(data, sendStru.Xml_msgName, conn)
	if err != nil {
		log.Println("connsendFile error:", err)
	}
	//发送成功
	Mid, _ := strconv.Atoi(sendStru.Massageid)
	//原始交易消息包发送成功更新 发送状态 发送时间 发送成功后消息包的文件路径
	err1, DBsj := storage.SendedUpdateYuansjyxx(int64(Mid), sendStru.Xml_msgName)
	if err1 != nil {
		log.Println("storage.SendedUpdateYuansjyxx  error:", err1)
	}
	//TCP发送记录
	has, serr, count := storage.GetTcpSendRecord(strconv.Itoa(Mid))
	if serr != nil {
		log.Println("查询TCP发送记录错误")
	}
	var record storage.BJsTcpqqjl

	if has == false {
		log.Println("此msgid属于第一次发送")
		record.FVcXiaoxxh = strconv.Itoa(Mid)
		record.FDtZuixsj = DBsj
		record.FNbChongfcs = 0
		record.FNbFasz = 1
		record.FNbMd5 = sendStru.Md5_str
		err2 := storage.TcpSendRecordInsert(record)
		if err2 != nil {
			log.Println("storage.TcpSendRecordInsert error:", err2)
		}
	}
	if has == true {
		log.Println("此msgid 已经发送")
		record.FVcXiaoxxh = strconv.Itoa(Mid)
		record.FNbChongfcs = count + 1
		err3 := storage.TcpSendRecordUpdate(record)
		if err3 != nil {
			log.Println("storage.TcpSendRecordUpdate error:", err3)
		}
	}
	log.Println("TCP发送记录 插入完成")

}
