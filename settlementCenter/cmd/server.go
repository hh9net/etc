package cmd

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
				fmt.Println("EOF",rerr)
			} else {
				fmt.Println("Read", rerr)
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

func main11() {
	//监听
	listen, lerr := net.Listen("tcp", "127.0.0.1:10240")
	if lerr != nil {
		fmt.Println("Listen", lerr)
		return
	}

	fmt.Println("等待客户端发送文件")
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
