package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

func sendFile(path string, connect net.Conn) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}
	defer file.Close()
	buff := make([]byte, 1024*4)
	for {
		size, rerr := file.Read(buff)
		if rerr != nil {
			if rerr == io.EOF {
				fmt.Println("EOF", rerr)

			} else {
				fmt.Println("Read:", rerr)
			}
			return
		}
		_, err := connect.Write(buff[:size])
		if err != nil {
			fmt.Println("err", err)
			return
		}
	}
}

//发送
func Sendxml() string {
	//如果发送成功 "ok" 返回一个成功 触发 文件的移动
	path := "../sendfilexml/"
	info, serr := os.Stat(path)
	if serr != nil {
		fmt.Println("Stat error", serr)
		return ""
	}
	//连接服务器
	conn, derr := net.Dial("tcp", "127.0.0.1:8081")
	if derr != nil {
		fmt.Println("Dial", derr)
		return ""
	}
	defer conn.Close()

	_, w_err := conn.Write([]byte(info.Name()))

	if w_err != nil {
		fmt.Println("Write error", w_err)
		return ""
	}

	buff := make([]byte, 4096)
	size, r_err := conn.Read(buff)

	if r_err != nil {
		fmt.Println("Read error", r_err)
		return ""
	}

	if "ok" == string(buff[:size]) {
		//如果收到ok应答，发送文件
		sendFile(path, conn)
	}
	return "ok"
}
