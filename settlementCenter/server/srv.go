package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

//模拟联网中心

func Server() {
	//创建用与监听客户端连接请求的socket
	listener, err := net.Listen("tcp", "127.0.0.1:8808")
	if err != nil {
		fmt.Println("Listen err:", err)
		return
	}
	defer listener.Close()
	//阻塞等待客户端（浏览器）连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept err:", err)
		return
	}
	defer conn.Close()
	//调用函数 获取连接客户端的网络地址，组织连接成功
	fmt.Println(conn.RemoteAddr().String(), "连接成功")
	//读取客户端发来的数据
	//数据缓冲区
	buf := make([]byte, 4096)
	n, err := conn.Read(buf) //n 接收数据的长度
	if err != nil {
		fmt.Println("Read err:", err)
		return
	}
	result := buf[:n] //切片截取
	fmt.Printf("#\n%s#", string(result))
}

// 浏览器访问时，该函数被回调
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello http"))
}
func Server2() {
	// 注册处理函数
	http.HandleFunc("/hello", handler)
	// 绑定服务器监听地址
	http.ListenAndServe("127.0.0.1:8000", nil)
}

//打开文件  http://127.0.0.1:8808/00000000000000079766.xml
func openSendFile(fileName string, w http.ResponseWriter) {
	filePath := "./jnp" + fileName
	// 只读打开文件：
	f, err := os.Open(filePath)
	// 说明 文件不存在。写错误页面
	if err != nil {
		errHtml := "<html><head><title>404 Not Found</title></head><body bgcolor=\"#CCFFCC\"><h4>404 NOT FOUND</h4><hr> sorry </body></html>"
		w.Write([]byte(errHtml))
	}
	defer f.Close()

	buf := make([]byte, 4096)
	// 文件存在， 循环读取，写给浏览器
	for {
		n, _ := f.Read(buf)
		if n == 0 {
			break
		}
		w.Write(buf[:n])
	}
}

//处理函数
func serverHandle(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.String()
	fmt.Println("urlfileName=", fileName)

	// 封装函数，去到服务器指定目录中找寻文件，存在打开写会给浏览器， 不存在报错
	openSendFile(fileName, w)
}

//取出文件
func Server3() {
	http.HandleFunc("/", serverHandle)
	http.ListenAndServe("127.0.0.1:8808", nil)
}
