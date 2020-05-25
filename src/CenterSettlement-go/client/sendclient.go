package client

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func sendFile(path string, connect net.Conn) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return
	}
	defer file.Close()
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
func Sendxml1() string {
	//如果xml发送成功  会返回一个成功 触发 文件的移动 "ok"
	//path := "../sendzipxml/"
	//info, serr := os.Stat(path)
	//if serr != nil {
	//	fmt.Println("Stat error", serr)
	//	return ""
	//}

	//连接服务器
	conn, derr := net.Dial("tcp", "127.0.0.1:8808")
	if derr != nil {
		log.Println("Dial", derr)
		return ""
	}
	log.Println("Dial 成功", conn)

	//把报文写给服务端

	//_, w_err := conn.Write([]byte(info.Name()))
	_, w_err := conn.Write([]byte("原始交易数据"))
	if w_err != nil {
		log.Println("Write error", w_err)
		return ""
	}

	buff := make([]byte, 4096)
	//读取返回内容
	size, r_err := conn.Read(buff)
	if r_err != nil {
		log.Println("Read error", r_err)
		return ""
	}
	//判断返回内容（联网中心的响应）
	if "ok" == string(buff[:size]) {
		//如果收到ok应答，发送文件
		//sendFile(path, conn)
	}
	conn.Close()
	return "ok"
}

//客户端对象
type TcpClient struct {
	connection *net.TCPConn
	server     *net.TCPAddr
	stopChan   chan struct{}
}

func Sendxml() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8181")
	if err != nil {
		log.Println(err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	send(conn)
}

func send(conn *net.TCPConn) {
	defer conn.Close()
	decodeString, _ := hex.DecodeString("2300a78c070000000000000a00000000000000030000000000010618120250ba2700004eb0")
	//发送
	_, err := conn.Write(decodeString)
	if err != nil {
		log.Println(err.Error())
	}
}

// 接收数据包
func (client *TcpClient) receivePackets() {
	reader := bufio.NewReader(client.connection)
	for {
		//承接上面说的服务器端的偷懒，我这里读也只是以\n为界限来读区分包
		msg, err := reader.ReadString('\n')
		if err != nil {
			//在这里也请处理如果服务器关闭时的异常
			close(client.stopChan)
			break
		}
		log.Println(msg)
	}
}
func (client *TcpClient) send() {
}
