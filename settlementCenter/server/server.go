package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"settlementCenter/client"
	"time"
)

//读取联网中心的数据
func ReadFile(fileName string, connect net.Conn) {
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
	defer connect.Close()
	buff := make([]byte, 1024*4)
	size, rerr := connect.Read(buff)
	if rerr != nil {
		fmt.Println("Read", rerr)
		return
	}
	fileName := string(buff[:size])
	connect.Write([]byte("ok"))
	//把读到的文件
	ReadFile(fileName, connect)
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
		// go Response(connect);
		go HandleTask(connect)
	}
	defer listen.Close()
}

//处理任务
func HandleTask(conn net.Conn) {
	//起goroutine handle函数
	//handle函数为总的逻辑处理入口
	//在这个handle函数中起四个goroutine
	//ch:=make(chan )
	//goroutine 1:
	go HandleTable()
	//goroutine 2:
	go HandleSendXml()

	//goroutine 3:
	go HandleChannelMessage(conn)
	//goroutine 4:
	go AnalyzeDataPakage()

}

//数据打包
func HandleTable() {

	//查询原始交易数据（在数据层）
	//准备数据（在数据层）
	//Xml数据生成Xml文件、压缩，存文件
	fname := Generatexml()

	//压缩
	Lz77Compress(fname)

	//Md5计算

	//7z压缩（Cgo解决）
	//插入原始交易消息包记录
	//更新原始交易数据的状态

	//调用打包函数

}

//线程2 发送数据包
func HandleSendXml() {
	//从文件夹sendxml中扫描打包文件（判断这个文件夹下面有没有文件）
	tiker := time.NewTicker(time.Second * 2)
	for {
		fmt.Println(<-tiker.C)
	}
	//if true
	//read 后  这里应该是要调 数据服务的一个接口
	client.Sendxml()
	//调接口成功后 mv文件夹到另一个文件中
	//
	//定时器定期扫描sendzipxml文件
	//	读取文件
	//	准备报文
	//	发送报文
	//	发送成功则mv消息包至sendsucceed
	//	发送失败 触发重发机制
}

//线程3  接收数据包
func HandleChannelMessage(conn net.Conn) {

	//监听channel数据，将接收到的数据存储到指定文件夹
	Response(conn)
	//监听联网中心端口
	//	接收数据包
	//	存储数据包至receive文件夹
}

//线程4 处理数据包  定期扫描 接收联网的接收数据的文件夹 receive，如果有文件就解压， 解压后分析包。
func AnalyzeDataPakage() {

	//定期检查文件夹receive    解压后
	tiker := time.NewTicker(time.Second * 5)
	for {
		<-tiker.C
		//扫描receive 文件夹

	}
	//定时器定期扫描receive文件夹
	//		读取文件
	//		解析文件
	//		应答数据包
	//		记账数据包
	//		争议数据包
	//		清分数据包
	//		退费数据包
}
