package main

import (
	"fmt"
	"net"
	"strings"
)

// 用于与客户端数据通信
func handlerConnect(conn net.Conn)  {
	// 延迟关闭与客户端进行通信的 conn
	defer conn.Close()

	// 获取客户端的 IP+port

	cltAddr := conn.RemoteAddr()
	fmt.Println("[", cltAddr, "]:客户端连接成功!")

	// 创建存储数据的 缓冲区
	buf :=make([]byte, 4096)

	for {
		// 1. 读取客户端发送数据
		n, err := conn.Read(buf)
		fmt.Println("---读到:", buf[:n])
		if string(buf[:n]) == "exit\n" || string(buf[:n-2]) == "exit" {  //  \r\n
			fmt.Println("服务器接收客户端关闭请求，退出！")
			return
		}
		if n == 0 {
			fmt.Println("服务器检测到客户端 粗暴关闭， 本端也退出！")
			return
		}
		if err != nil {
			fmt.Println("read err:", err)
			return
		}
		fmt.Println("服务器读到：", string(buf[:n]))

		// 2. 处理数据 小——大写
		upperStr := strings.ToUpper(string(buf[:n]))
		// 3. 回发给客户端
		conn.Write([]byte(upperStr))
	}
}

func main() {
	// 创建监听 socket
	listener, err := net.Listen("tcp", "127.0.0.1:8003")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// 循环 阻塞监听客户端连接请求
	for {
		fmt.Println("服务器等待中...")
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// 创建子go程，专门用于与client 进行数据通信
		go handlerConnect(conn)
	}
}
