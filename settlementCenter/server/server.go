package server

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"settlementCenter/client"
)

func ReadFile(fileName string, connect net.Conn) {
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Create", ferr)
		return
	}
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
	go HandleChannelMessage()
	//goroutine 4:
	go AnalyzeDataPakage()

}
func HandleTable() {
	//连接数据库 操作  链接数据库这个逻辑可以放外面统一的句柄
	//也可以使用链接池：
	//select * FROM B_TXF_CHEDXFYSSJ WHERE （按照五分种和未打包筛选）
	//查出来五分种的条数 判断不足一百条则取五分种内所有  超过一百条取前一百条
	//将结构体数据解析成XML格式并写入文件，
	v := ""
	xmlOutPut, outPutErr := xml.Marshal(v)
	if outPutErr != nil {
		fmt.Println("error:", outPutErr)
	}
	fmt.Println(string(xmlOutPut))
	//保存文件
	if outPutErr == nil {
		headerBytes := []byte(xml.Header)                  //加入XML头
		xmlOutPutData := append(headerBytes, xmlOutPut...) //拼接XML头和实际XML内容
		//写入文件
		ioutil.WriteFile("文件路径", xmlOutPutData, os.ModeAppend)
	} else {
		fmt.Println(outPutErr)
	}

	//打包
	client.Compress()

}

//线程2 发送数据包
func HandleSendXml() {
	//从文件夹sendxml中扫描打包文件（判断这个文件夹下面有没有文件）
	//tiker := time.NewTicker(time.Second * 2)
	//for  {
	//	fmt.Println(<-tiker.C)
	//}
	//if true
	//read 后  这里应该是要调 数据服务的一个接口
	client.Sendxml()
	//调接口成功后 mv文件夹到另一个文件中
	//
}

//线程3  接收数据包
func HandleChannelMessage() {

	//监听channel数据，将接收到的数据存储到指定文件夹
	//
	Receive()
}

//线程4 处理数据包
func AnalyzeDataPakage() {

	//定期检查上面的文件夹    解压后
}
