package server

import (
	"CenterSettlement-go/types"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type DataPacket struct {
	Type string
	Body string
}

//模拟联网中心，处理结算数据   业务 模拟
func Server() {
	//绑定端口
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8808")
	if err != nil {
		log.Println(err.Error())
	}
	//监听
	log.Println("监听客户端 127.0.0.1:8808")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err.Error())
	}
	defer listener.Close()

	//联网中心开始接收数据
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		go ReceiveHandler(conn)
	}
}

func ReceiveHandler(conn net.Conn) {
	defer conn.Close()
	//每次读取数据长度
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	data := string(buf[:n])
	log.Println("接收数据")
	log.Println(data)

	msgid := string(buf[:20])
	log.Println("消息包Massageid", msgid)

	msglength := string(buf[20:26])
	log.Println("消息包msglength", msglength)

	msgmd5 := string(buf[26:58])
	log.Println("消息包msgmd5：", msgmd5)

	msg := string(buf[58:])
	log.Println("消息包msg：")
	log.Println(msg)

	//保存为文件
	RevFile(msgid, buf[58:])

	//即时应答
	var replyStru types.ReplyStru
	replyStru.Massageid = "00000000000000100025"
	replyStru.Result = "1"
	m := []byte(replyStru.Massageid)
	r := []byte(replyStru.Result)
	d := append(m, r...)
	log.Println(string(d))
	// 返回ok
	conn.Write(d)

	//如果文件比较大要循环读取

	//把读取的内容保存成文件 存放在center中

	//判断是否接收完， 如果接收完毕，则回复应答
	//准备应答报文

	// 封装函数，去到服务器指定目录中找寻文件，存在打开写会给浏览器， 不存在报错
	//openSendFile(fileName, w)
}

func RevFile(fileName string, data []byte) {

	//err := ioutil.WriteFile("test.txt",data, 0644)
	err := ioutil.WriteFile("../receive/"+fileName+".xml.lz77", data, os.ModeAppend)
	if err != nil {
		log.Println("联网中心保存文件失败  ioutil.WriteFile  err =", err)
		return
	}

	//fs,err := os.Create("../receice/"+fileName)
	//defer fs.Close()
	//if err != nil {
	//	log.Println("os.Create err =",err)
	//	return
	//}
	//// 拿到数据
	////buf := make([]byte ,1024*10)
	//for {
	//	n,err := fs.Read(buf)
	//	if err != nil {
	//		log.Println("fs.Read error =",err)
	//		if err == io.EOF {
	//			log.Println("文件结束了",err)
	//		}
	//		return
	//	}
	//	if n == 0 {
	//		log.Println("文件结束了",err)
	//		return
	//	}
	//	n,werr :=fs.Write(buf[:n])
	//	if werr != nil {
	//		log.Println("联网中心保存文件失败",err)
	//	}
	//}
}
