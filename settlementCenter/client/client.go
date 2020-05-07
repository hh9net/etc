package client

import (
	"fmt"
	"net"
)

func Client() {
	//1、客户端主动连接服务器 http://127.0.0.1:8808/ysjyxx0000011.xml
	conn, err := net.Dial("tcp", "127.0.0.1:8808")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}
	fmt.Println("tcp连接成功")
	defer conn.Close() //条件反射出来 延迟关闭

	//2、模拟浏览器，组织一个最简单的请求报文。包含请求行，请求头，空行即可。
	requestHttpHeader := "GET /ysjyxx0000011.xml HTTP/1.1\r\nHost:127.0.0.1:8808\r\n\r\n"
	//3、给服务器发送请求报文
	conn.Write([]byte(requestHttpHeader))

	//4、读取 服务器回复 响应报文
	// 读取响应缓冲区
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read err:", err)
		return
	}
	//5、打印观察
	fmt.Printf("#\n%s#", string(buf[:n]))
}
