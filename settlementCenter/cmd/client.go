package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
)

func SendFile(path string, connect net.Conn) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}
	defer  file.Close()
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
func Send() {
	path := "../sendfilexml"

	info, serr := os.Stat(path)
	if serr != nil {
		fmt.Println("Stat error", serr)
		return
	}

	conn, derr := net.Dial("tcp", "127.0.0.1:8081")
	if derr != nil {
		fmt.Println("Dial", derr)
		return
	}

	_, werr := conn.Write([]byte(info.Name()))
	if werr != nil {
		fmt.Println("Write error", werr)
		return
	}

	buff := make([]byte, 4096)
	size, rerr := conn.Read(buff)
	if rerr != nil {
		fmt.Println("Read error", rerr)
		return
	}

	if "ok" == string(buff[:size]) {
		SendFile(path, conn)
	}

	defer conn.Close()
}
