package server

import (
	"fmt"
	"io"
	"net"
	"os"
)

func ReadFile(fileName string, connect net.Conn) {
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Create", ferr);
		return
	}
	buff := make([]byte, 1024*4)
	for {
		size, rerr := connect.Read(buff)
		if rerr != nil {
			if rerr == io.EOF {
				fmt.Println("EOF error",rerr)
			} else {
				fmt.Println("Read error", rerr)
			}
			return
		}
		file.Write(buff[:size])
	}
}

func Response(connect net.Conn) {
	defer connect.Close()
	buff := make([]byte, 1024*4)
	size, rerr := connect.Read(buff)
	if rerr != nil {
		fmt.Println("Read", rerr)
		return
	}
	fileName := string(buff[:size])
	connect.Write([]byte("ok"))
	ReadFile(fileName, connect)
}
//接收联网中心发来数据包
func Receive() {
	//监听
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
		go Response(connect);
	}
	defer listen.Close()
}
