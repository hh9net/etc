package cmd

import (
	"fmt"
	"net"
	"os"
)

var s = ""

func Client() {
	// 创建 用于数据通信 conn
	conn, err := net.Dial("tcp", "127.0.0.1:8003")
	if err != nil {
		panic(err)
	}
	var i = 10
	fmt.Println(i)
	// 创建 子go程 读用户键盘输入
	go func() {
		for {
			//fmt.Scan()		// 遇见空格、\n 会自动终止读取。
			str := make([]byte, 4096)
			n, err := os.Stdin.Read(str)
			if err != nil {
				panic(err)
			}
			// 将读到的数据发送给服务器, 读多少写多少。
			_, err = conn.Write(str[:n])
			if err != nil {
				panic(err)
			}
		}
	}()

	buf := make([]byte, 4096)
	// 主go程 ，循环， 读取服务器回发数据
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			fmt.Println("客户端检测到服务器关闭")
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("读到服务器回发：", string(buf[:n]))
	}
}
