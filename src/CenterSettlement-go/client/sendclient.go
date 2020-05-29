package client

import (
	"CenterSettlement-go/types"
	"io"
	"log"
	"net"
	"os"
)

func connsendFile(data []byte, fname string, connect net.Conn) {

	path := "../sendzipxml/" + fname
	log.Println("path:=", path)
	file, oserr := os.Open(path)
	if oserr != nil {
		log.Println("os.Open error", oserr)
		return
	}
	log.Println("文件conn发送")
	defer file.Close()

	buff := make([]byte, 1024*10)
	for {
		size, rerr := file.Read(buff)
		if rerr == io.EOF {
			log.Println("文件读取完毕")
			return
		}
		if rerr != nil {
			log.Println("file.Read err:", rerr)
			return
		}
		data = append(data, buff[:size]...)
		log.Println(string(data))
		_, werr := connect.Write(data)
		if werr != nil {
			log.Println("err", werr)
			return
		}
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

	connsendFile(data, sendStru.Xml_msgName, conn)
}
