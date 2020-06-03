package client

import (
	storage "CenterSettlement-go/storages"
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
		log.Println("connsendFile err:", err)

	}
	//发送成功
	Mid, _ := strconv.Atoi(sendStru.Massageid)
	err1 := storage.SendedUpdateYuansjyxx(int64(Mid), sendStru.Xml_msgName)
	if err1 != nil {
		log.Println("storage.SendedUpdateYuansjyxx  err:", err1)

	}
}
