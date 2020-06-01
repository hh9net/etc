package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//保存联网中心的数据
func SaveFile(connect net.Conn) {
	var fileName string

	buf := make([]byte, 1024*4)
	n, rerr := connect.Read(buf)
	if rerr != nil {
		if rerr == io.EOF {
			fmt.Println("EOF error", rerr)
		} else {
			fmt.Println("Read error", rerr)
		}
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
	fileName = msgid
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Create", ferr)
		return
	}
	defer file.Close()

	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, buf[58:]...)
	//这里可以不写，直接使用channel发送给线程2
	//写入文件
	_, fwerr := file.Write((xmlOutPutData))
	if fwerr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
	}
	//更新消息包信息
	file.Close()

}

//响应联网中心的应答
func Response(connect net.Conn) {
	connect.Write([]byte("ok"))
}

//接收联网中心发来数据包
func Receive() {
	log.Println("执行线程4")
	//监听联网中心数据端口
	listen, lerr := net.Listen("tcp", "127.0.0.1:8808")
	if lerr != nil {
		fmt.Println("Listen", lerr)
		return
	}
	log.Println("等待联网中心发送文件")
	for {
		connect, cerr := listen.Accept()
		if cerr != nil {
			fmt.Println("Accept", cerr)
			return
		}
		go HandleTask(connect)
	}
	defer listen.Close()
}

//处理任务
func HandleTask(conn net.Conn) {
	go HandleMessage(conn)
}

//线程3  接收数据包
func HandleMessage(conn net.Conn) {
	defer conn.Close()
	//fileName := ""
	//把读到的数据 以文件记录
	SaveFile(conn)
	//接收数据即时应答
	Response(conn)

	//这里可以选择不存储，直接使用channnel将数据发送给线程4
	//监听channel数据，将接收到的数据存储到指定文件夹
	conn.Write([]byte("ok"))
	//监听联网中心端口
	//	接收数据包
	//	存储数据包至receive文件夹
}
