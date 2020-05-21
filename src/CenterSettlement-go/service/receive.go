package service

import (
	"fmt"
	"io"
	"net"
	"os"
)

//保存联网中心的数据
func SaveFile(fileName string, connect net.Conn) {
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Create", ferr)
		return
	}
	defer file.Close()
	buff := make([]byte, 1024*4)
	for {
		size, rerr := connect.Read(buff)
		if rerr != nil {
			if rerr == io.EOF {
				fmt.Println("EOF error", rerr)
			} else {
				fmt.Println("Read error", rerr)
			}
			return
		}
		file.Write(buff[:size])
	}
}

//响应联网中心的应答
func Response(connect net.Conn) {
	connect.Write([]byte("ok"))
}

//接收联网中心发来数据包
func Receive() {
	//监听联网中心数据端口
	listen, lerr := net.Listen("tcp", "127.0.0.1:8808")
	if lerr != nil {
		fmt.Println("Listen", lerr)
		return
	}
	fmt.Println("等待联网中心发送文件")
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
	fileName := ""
	Response(conn)
	//把读到的数据 以文件记录
	SaveFile(fileName, conn)
	//这里可以选择不存储，直接使用channnel将数据发送给线程4
	//监听channel数据，将接收到的数据存储到指定文件夹
	conn.Write([]byte("ok"))
	//监听联网中心端口
	//	接收数据包
	//	存储数据包至receive文件夹
}
