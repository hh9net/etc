package client

import (
	storage "CenterSettlement-go/storage"
	"CenterSettlement-go/types"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func connsendFile(data []byte, fname string, connect net.Conn) error {

	path := "../sendzipxml/" + fname + ".lz77"
	log.Println("path:=", path)
	file, oserr := os.Open(path)
	if oserr != nil {
		log.Println("os.Open error", oserr)
		return oserr
	}
	log.Println("文件conn发送")
	defer file.Close()

	buff := make([]byte, 1024*10)
	for {
		size, rerr := file.Read(buff)
		if rerr == io.EOF {
			log.Println("文件读取完毕")
			return rerr
		}
		if rerr != nil {
			log.Println("file.Read err:", rerr)
			return rerr
		}
		data = append(data, buff[:size]...)
		log.Println(string(data))
		_, werr := connect.Write(data)
		if werr != nil {
			log.Println("err", werr)
			return werr
		}
		log.Println("发送报文到联网中心成功")
		return nil
	}
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
	var record storage.BJsTcpqqjl
	record.FVcXiaoxxh = strconv.Itoa(Mid)
	record.FDtZuixsj = DBsj
	record.FNbChongfcs = 0
	record.FNbFasz = 1
	record.FNbMd5 = sendStru.Md5_str
	err2 := storage.TcpSendRecordInsert(record)
	if err2 != nil {
		log.Println("storage.TcpSendRecordInsert error:", err2)
	}
	log.Println("TCP发送记录 成功")

}
