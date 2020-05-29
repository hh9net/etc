package server

import (
	"CenterSettlement-go/types"
	"log"
	"net"
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
	log.Println("接收数据", data)
	log.Println(data)
	var replyStru types.ReplyStru
	replyStru.Massageid = "00000000000000100025"
	replyStru.Result = "1"
	m := []byte(replyStru.Massageid)
	r := []byte(replyStru.Result)
	d := append(m, r...)
	log.Println(string(d))
	// 返回ok
	conn.Write(d)
	//即时应答
	//conn.Write([]byte ("ok"))

	//result, Body := check(buf)
	//if result {
	//	log.Printf("接收到报文内容:{ %s }\n", hex.EncodeToString(Body))
	//}

	//如果文件比较大要循环读取

	//把读取的内容保存成文件 存放在center中

	//判断是否接收完， 如果接收完毕，则回复应答
	//准备应答报文

	// 封装函数，去到服务器指定目录中找寻文件，存在打开写会给浏览器， 不存在报错
	//openSendFile(fileName, w)
}
