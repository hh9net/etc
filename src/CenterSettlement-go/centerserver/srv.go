package centerserver

import (
	"CenterSettlement-go/centerYuanshi"
	"CenterSettlement-go/client"
	"CenterSettlement-go/types"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataPacket struct {
	Type string
	Body string
}

//模拟联网中心，处理结算数据   业务 模拟
func Server() {
	//绑定端口8806
	address := AddressConfigInit()
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
	//接受客户端发来的要传文件的文件信息data
	buffer := make([]byte, 100)
	d, e := conn.Read(buffer)
	log.Println("文件数据data：", string(buffer[:d]), e)
	fileNameid := string(buffer[:20])

	msglength := string(buffer[20:26])
	log.Println("消息包msglength", msglength)

	msgmd5 := string(buffer[26:58])
	log.Println("消息包msgmd5：", msgmd5)

	//msgid := string(buffer[:20])
	//log.Println("消息包Massageid", msgid)
	//创建文件
	fs, err := os.Create("../centerserver/" + fileNameid + ".xml.lz77")
	defer fs.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return
	}

	j, _ := strconv.Atoi(msglength)
	total := 0
	i := 0
	//每次读取数据长度
	buf := make([]byte, 100)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		total += n
		i += 1

		if n == 0 {
			log.Println("文件读取完毕")
			break
		}
		if err != nil {
			log.Println("conn.Read err:", err)
			return
		}
		fs.Write(buf[:n])

		//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
		if total == j {
			fmt.Println("文件接受成功,共", total, "字节")
			//回复客户端已收到文件
			//conn.Write([]byte("文件接受成功"))

			//即时应答
			var replyStru types.ReplyStru
			//replyStru.Massageid = string(buf[:20])
			replyStru.Massageid = fileNameid

			replyStru.Result = "1"
			m := []byte(replyStru.Massageid)
			r := []byte(replyStru.Result)
			d := append(m, r...)
			InstantResponse(d, conn)
			break
		}
		rerr := RevFile(string(buf[:20]), buf[58:])
		if rerr != nil {
			log.Println("文件保存失败")
		}
	}
	log.Println("循环写文件次数", i)
	//获取数据
	//GetData(buffer[:n])
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
	//保存为文件
	//iowerr := ioutil.WriteFile("../centerserver/"+fileNameid+".xml.lz77", buf[:n], 0644)
	//if iowerr != nil {
	//	log.Println("联网中心保存文件失败  ioutil.WriteFile  err =", iowerr)
	//	return
	//}
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

func HandleFile() {
	//扫描文件夹  解压缩文件
	//从文件夹center_server中扫描打包文件（判断这个文件夹下面有没有文件）

	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("扫描联网中心文件夹", <-tiker.C)
		//pwd := "centerserver/"
		pwd := "./"

		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		for i := range fileInfoList {
			//判断文件的结尾名
			if strings.HasSuffix(fileInfoList[i].Name(), ".xml.lz77") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())
				//解压缩
				//zerr:= UnZipLz77(fileInfoList[i].Name())//新版bug
				zerr := centerYuanshi.UnZipLz77(fileInfoList[i].Name())
				if zerr != nil {
					log.Println("发送文件时 压缩xml文件失败")
				}
				if zerr == nil {
					//移动xml文件
					s := "./" + fileInfoList[i].Name()
					des := "../centerYuanshi/zip/" + fileInfoList[i].Name()
					merr := client.MoveFile(s, des)
					if merr != nil {
						log.Println("client.MoveFile err : ", merr)
						return
					}
				}

				//解析文件
				//		解析文件  获取数据
				ParsingYSXMLFile()

				//log.Println(sendStru)

				////连接联网中心服务器
				//address := conf.AddressConfigInit()
				//Address := address.AddressIp + ":" + address.AddressPort
				//conn, derr := net.Dial("tcp", Address)
				//if derr != nil {
				//	log.Println("Dial", derr)
				//	//return ""
				//}
				//if conn != nil {
				//	log.Println("Dial 成功")
				//}
				////发送
				//client.Sendxml(&sendStru, conn)
				//
				//buf := make([]byte, 1024)
				//n, err2 := conn.Read(buf)
				//if err2 != nil {
				//	log.Println("conn.Read err = ", err2)
				//	return
				//}
				////str := string(buf[:n])
				////对联网中心的接收应答处理
				//ImmediateResponseProcessing(string(buf[:n]), fileInfoList[i].Name())
				//
				//conn.Close()
			} else {
				log.Println("该联网中心没有需要解析原始交易消息xml")
				break
			}
		}
	}
}

//解析原始交易数据
func ParsingYSXMLFile() {
	//扫描原始数据包
	pwd := "../centerYuanshi/"

	//pwd := "CenterSettlement-go/receivexml/"
	fileList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileList))
	if len(fileList) == 0 {
		log.Println("该centerYuanshi 文件夹下没有需要解析的文件")
		return
	}
	for i := range fileList {
		log.Println("xml name:", fileList[i].Name()) //打印当前文件或目录下的文件或目录名
		//判断文件的结尾名
		if strings.HasSuffix(fileList[i].Name(), ".xml") {

			//解析xml文件

			//content, err := ioutil.ReadFile("../receivexml/" + fileInfoList[i].Name())
			//
			content, err := ioutil.ReadFile("../centerYuanshi/" + fileList[i].Name())
			if err != nil {
				log.Println("读文件位置错误信息：", err)
				return
			}

			//将xml文件转换为对象
			var msg Message
			err = xml.Unmarshal(content, &msg)
			if err != nil {
				log.Println("解析 receive文件夹中xml文件内容的错误信息：", err)
			}

			log.Println("msg:", msg.Header.MessageClass, msg.Header.MessageType, msg.Body.ContentType, msg.Header.MessageId)
			//原始交易数据
			if msg.Header.MessageClass == 5 && msg.Header.MessageType == 7 && msg.Body.ContentType == 1 {
				//原始交易数据
				//解析xml数据 把数据导入数据库
				Parsexml("../centerClient/", fmt.Sprintf("%020d", msg.Header.MessageId)+".xml")

				return
			}

			//if msg.Header.MessageClass == 5 && msg.Header.MessageType == 5 && msg.Body.ContentType == 1 {
			//	//记账数据包
			//	//1、修改文件名字  2、移动文件
			//	src := "/Users/nicker/Desktop/Xmlfilebak(3)/" + fileList[i].Name()
			//	des := "../keepAccountFile/" + "JZB_" + fmt.Sprintf("%020d", msg.Header.MessageId) + ".xml"
			//	frerr := common.FileRename(src, des)
			//	if frerr != nil {
			//		log.Println("记账数据包 修改文件名字错误：", frerr)
			//		return
			//	}

			//解析xml数据 把数据导入数据库
			//	Parsexml("../keepAccountFile/", "JZB"+fmt.Sprintf("%020d", msg.Header.MessageId)+".xml")
			//}
		}
	}

}

func Parsexml(pwd string, fname string) {

	log.Println("把原始数据xml插入数据库", pwd, fname)

	//把解析过的原始交易消息xml移动文件夹
}
