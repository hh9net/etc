package server

import (
	"CenterSettlement-go/conf"
	"CenterSettlement-go/types"
	"io/ioutil"
	"log"
	"net"
)

type DataPacket struct {
	Type string
	Body string
}

//模拟联网中心，处理结算数据   业务 模拟
func Server() {
	//绑定端口8806
	address := conf.AddressConfigInit()
	Address := address.AddressIp + ":" + address.AddressPort
	tcpAddr, err := net.ResolveTCPAddr("tcp", Address)
	if err != nil {
		log.Println(err.Error())
	}
	//监听
	log.Println("监听客户端 127.0.0.1:8806")
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err.Error())
	}
	defer listener.Close()

	//联网中心开始接收数据
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		go ReceiveHandler(conn)
	}
}

func ReceiveHandler(conn net.Conn) {
	defer conn.Close()
	//每次读取数据长度
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	//获取数据
	GetData(buf[:n])

	//保存为文件
	rerr := RevFile(string(buf[:20]), buf[58:])
	if rerr != nil {
		log.Println("文件保存失败")
	}

	//即时应答
	var replyStru types.ReplyStru
	replyStru.Massageid = string(buf[:20])
	replyStru.Result = "1"
	m := []byte(replyStru.Massageid)
	r := []byte(replyStru.Result)
	d := append(m, r...)
	InstantResponse(d, conn)

	//如果文件比较大要循环读取
	//把读取的内容保存成文件 存放在center中
	//判断是否接收完， 如果接收完毕，则回复应答
	//准备应答报文
	// 封装函数，去到服务器指定目录中找寻文件，存在打开写会给浏览器， 不存在报错
	//openSendFile(fileName, w)
}

//应答
func InstantResponse(d []byte, conn net.Conn) {
	// 返回接收成功
	_, err := conn.Write(d)
	if err != nil {
		log.Println("联网中心 conn.Write 错误")
	}

	conn.Close()
}

//获取数据
func GetData(buf []byte) {
	data := string(buf)
	log.Println("接收数据:")
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
}

//保存为文件
func RevFile(fileName string, data []byte) error {
	//err := ioutil.WriteFile("test.txt",data, 0644)
	err := ioutil.WriteFile("../center_server/"+fileName+".xml.lz77", data, 0644)
	if err != nil {
		log.Println("联网中心保存文件失败  ioutil.WriteFile  err =", err)
		return err
	}
	log.Println("联网中心保存文件成功")
	return nil
	//fs,err := os.Create("../receice/"+fileName)
	//defer fs.Close()
	//if err != nil {
	//	log.Println("os.Create err =",err)
	//	return
	//}
	//// 拿到数据
	////buf := make([]byte ,1024*10)
	//for {
	//	n,err := fs.Read(buf)
	//	if err != nil {
	//		log.Println("fs.Read error =",err)
	//		if err == io.EOF {
	//			log.Println("文件结束了",err)
	//		}
	//		return
	//	}
	//	if n == 0 {
	//		log.Println("文件结束了",err)
	//		return
	//	}
	//	n,werr :=fs.Write(buf[:n])
	//	if werr != nil {
	//		log.Println("联网中心保存文件失败",err)
	//	}
	//}
}
