package server

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type DataPacket struct {
	Type string
	Body string
}

//模拟联网中心，处理结算数据   业务 模拟
func Server() {
	//绑定端口
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8808")
	if err != nil {
		log.Println(err.Error())
	}
	//监听
	log.Println("监听客户端 127.0.0.1:8808")
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
	_, err := conn.Read(buf)
	if err != nil {
		return
	}
	result, Body := check(buf)
	if result {
		log.Printf("接收到报文内容:{ %s }\n", hex.EncodeToString(Body))
	}
}
func check(buf []byte) (bool, []byte) {
	Length := DataLength(buf)
	if Length < 3 || Length > 4096 {
		return false, nil
	}
	Body := buf[:Length]
	return uint16(len(Body))-2 != Length, Body
}
func DataLength(buf []byte) uint16 {
	return binary.BigEndian.Uint16(inversion(buf[:2])) + 2
}

//反转字节
func inversion(buf []byte) []byte {
	for i := 0; i < len(buf)/2; i++ {
		temp := buf[i]
		buf[i] = buf[len(buf)-1-i]
		buf[len(buf)-1-i] = temp
	}
	return buf
}

func Server1() {
	//监听客户端
	log.Println("监听客户端 127.0.0.1:8808")
	http.HandleFunc("/", ServerHandle)
	http.ListenAndServe("127.0.0.1:8808", nil)
}

//处理函数
func ServerHandle(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.String()
	log.Println("urlfileName=", fileName)

	con, _ := ioutil.ReadAll(r.Body) //获取post的数据
	log.Println(string(con))
	w.Write([]byte("已经接收"))
	r.Body.Close()

	//如果文件比较大要循环读取

	//把读取的内容保存成文件 存放在center中

	//判断是否接收完， 如果接收完毕，则回复应答
	//准备应答报文

	// 封装函数，去到服务器指定目录中找寻文件，存在打开写会给浏览器， 不存在报错
	//openSendFile(fileName, w)

}
