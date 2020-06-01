package service

import (
	"CenterSettlement-go/types"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

//保存联网中心的数据
func SaveFile(connect net.Conn) string {
	var fileName string

	buf := make([]byte, 1024*4)
	n, rerr := connect.Read(buf)
	if rerr != nil {
		if rerr == io.EOF {
			fmt.Println("EOF error", rerr)
		} else {
			fmt.Println("Read error", rerr)
		}
		return ""
	}

	data := string(buf[:n])
	log.Println("接收数据")
	log.Println(data)

	msgid := string(buf[:20])
	log.Println("消息包Massageid", msgid)

	msglength := string(buf[20:26])
	log.Println("消息包msglength", msglength)

	msgmd5 := string(buf[26:58])
	log.Println("消息包msgmd5：", msgmd5)

	msg := string(buf[58:])
	log.Println("消息包msg：")
	log.Println(msg)
	//包号作为文件的名字 ，后面再修改名字在解析文件之后，修改文件名
	fileName = msgid
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Create", ferr)
		return ""
	}

	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, buf[58:]...)
	//这里可以不写，直接使用channel发送给线程2
	//写入文件
	_, fwerr := file.Write((xmlOutPutData))
	if fwerr != nil {
		log.Printf("Write xml file error: %v\n", ferr)
	}
	//更新消息包信息
	defer file.Close()
	return fileName
}

//响应联网中心的应答
func Response(Filename string, connect net.Conn) {

	//即时应答
	var replyStru types.ReplyStru
	replyStru.Massageid = Filename
	replyStru.Result = "1"
	m := []byte(replyStru.Massageid)
	r := []byte(replyStru.Result)
	d := append(m, r...)
	log.Println(string(d))
	// 返回ok
	connect.Write(d)
}

//接收联网中心发来数据包
func Receive() {
	log.Println("执行线程4")
	//监听联网中心数据端口
	listen, lerr := net.Listen("tcp", "127.0.0.1:8808")
	if lerr != nil {
		fmt.Println("Listen", lerr)
		return
	}
	log.Println("等待联网中心发送文件")
	for {
		connect, cerr := listen.Accept()
		if cerr != nil {
			fmt.Println("Accept", cerr)
			return
		}
		go HandleTask(connect)
	}
	defer listen.Close()
}

//处理任务
func HandleTask(conn net.Conn) {
	go HandleMessage(conn)
}

//线程3  接收数据包
func HandleMessage(conn net.Conn) {

	//fileName := ""
	//把读到的数据 以文件记录
	Filename := SaveFile(conn)
	//接收数据即时应答
	Response(Filename, conn)
	defer conn.Close()

}
