package centerServer

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/common"
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//
//type DataPacket struct {
//	Type string
//	Body string
//}

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
	buffer := make([]byte, 58)
	d, e := conn.Read(buffer)
	log.Println("描述文件的数据data：", string(buffer[:d]), e)
	fileNameid := string(buffer[:20])

	msglength := string(buffer[20:26])
	log.Println("消息包msglength", msglength)

	msgmd5 := string(buffer[26:58])
	log.Println("消息包msgmd5：", msgmd5)

	//msgid := string(buffer[:20])
	//log.Println("消息包Massageid", msgid)

	//接收文件，保存文件，即使应答
	rferr := RevFile(fileNameid, conn, msglength)
	if rferr != nil {
		log.Println("文件保存失败", rferr)
		return
	}
	//校验文件
	cerr := CheckFile(msgmd5, fileNameid)
	if cerr != nil {

		// 1、应答确认
		log.Println("触发应答确认消息发送")

		//2、移动文件
		//移动xml
		x1 := "../centerYuanshi/" + fileNameid + ".xml"
		x2 := "../centeryuanshixmlparsed/errorxml/" + fileNameid + ".xml"
		mxerr := common.MoveFile(x1, x2)
		if mxerr != nil {
			log.Println("移动errorxml 失败")
			return
		}

		//移动xmlzip
		xz1 := "../centeryuanshixmlzip/" + fileNameid + ".xml.lz77"
		xz2 := "../centeryuanshixmlzip/errorxml/" + fileNameid + ".xml.lz77"
		mxzerr := common.MoveFile(xz1, xz2)
		if mxzerr != nil {
			log.Println("移动errorxml 失败")
			return
		}
	}

	//处理文件
	HandleFile()

	//获取数据
	//GetData(buffer[:n])
}

func CheckFile(msgmd5, fileNameid string) error {
	pwd := "./"
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))

	for i := range fileInfoList {
		//判断文件的结尾名
		if strings.HasSuffix(fileInfoList[i].Name(), ".lz77") {
			log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())

			//解压缩
			zerr := UnZipLz77(fileInfoList[i].Name())
			if zerr != nil {
				log.Println(" 解压缩文件失败", zerr)
				return zerr
			}
			if zerr == nil {
				//移动xml文件
				s := "./" + fileInfoList[i].Name()
				des := "../centeryuanshixmlzip/" + fileInfoList[i].Name()
				merr := client.MoveFile(s, des)
				if merr != nil {
					log.Println("client.MoveFile err : ", merr)
					return merr
				}

			}
			log.Println("校验文件")
			fstr := strings.Split(fileInfoList[i].Name(), ".xml.lz77")

			log.Println(fstr[0], fileNameid)
			if fileNameid == fstr[0] {
				//校验文件 校验其md5

				fmd5 := GetFileMd5(fileNameid + ".xml")
				if fmd5 == msgmd5 {
					log.Println("文件md5一致")
					log.Println("解析原始交易消息文件，  获取数据")
					//解析文件 导入数据库
					//ParsingYSXMLFile()
					return nil
				} else {
					log.Println("文件md5 不一致 ,通行宝 文件发送格式不正确")
					return errors.New("文件md5 不一致 ,通行宝 文件发送格式不正确")
				}
			}
		}
	}
	return nil
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

//接收文件 保存为文件 即时应答
func RevFile(fileNameid string, conn net.Conn, msglength string) error {
	//创建文件
	fs, err := os.Create("../centerServer/" + fileNameid + ".xml.lz77")
	defer fs.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return err
	}

	j, _ := strconv.Atoi(msglength) //压缩文件长度
	total := 0
	i := 0 //循环次数
	//每次读取数据长度
	buf := make([]byte, 100)
	for {
		//读取内容
		n, rerr := conn.Read(buf)
		if rerr != nil {
			log.Println("conn.Read err:", rerr)
			return rerr
		}
		total += n
		i += 1

		if n == 0 {
			log.Println("文件读取完毕")
			break
		}
		//写入文件
		fs.Write(buf[:n])

		//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
		if total == j {
			log.Println("文件接受成功,共", total, "字节")
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

	}
	log.Println("循环写文件次数", i)
	return nil
}

func HandleFile( /*msgmd5, fileNameid string*/ ) {
	//扫描文件夹  解压缩文件
	//从文件夹center_server中扫描打包文件（判断这个文件夹下面有没有文件）

	tiker := time.NewTicker(time.Second * 5)

	for {
		log.Println("扫描联网中心文件夹", <-tiker.C)

		//扫描原始数据包
		pwd := "../centerYuanshi/"
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

				log.Println("msg:", msg.Header.MessageClass, msg.Header.MessageType, msg.Body.ContentType, msg.Header.MessageId, msg.Body.Transaction[0].ICCard.CardType)
				//原始交易数据

				//解析文件 导入数据库
				ParsingYSXMLFile(msg, fileList[i].Name())

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
			}
			//else {
			//	log.Println("该联网中心没有需要解析原始交易消息xml")
			//	log.Println("fileInfoList[i].Name():",fileInfoList[i].Name())
			//}
		}
	}
}

//解析原始交易数据
func ParsingYSXMLFile(msg Message, fname string) error {

	if msg.Header.MessageClass == 5 && msg.Header.MessageType == 7 && msg.Body.ContentType == 1 && msg.Body.Transaction[0].ICCard.CardType == 22 {
		//储值卡 原始交易数据
		//1、改文件名
		Diqu := msg.Body.Transaction[0].ICCard.NetNo
		s1 := "../centerYuanshi/" + fname
		des := "../centerYuanshi/" + "CZ_" + Diqu + "_" + fname
		mverr := common.MoveFile(s1, des)
		if mverr != nil {
			return mverr
		}
		//2、解析xml数据 把数据导入数据库
		perr := Parsexml(des, "CZ_"+Diqu+"_"+fname, msg)
		if perr != nil {
			return perr
		}
		return nil
	}

	if msg.Header.MessageClass == 5 && msg.Header.MessageType == 7 && msg.Body.ContentType == 1 && msg.Body.Transaction[0].ICCard.CardType == 23 {
		//记账卡 原始交易数据
		//1、改文件名
		Diqu := msg.Body.Transaction[0].ICCard.NetNo
		s1 := "../centerYuanshi/" + fname
		des := "../centerYuanshi/" + "JZ_" + Diqu + "_" + fname
		common.MoveFile(s1, des)

		//2、解析xml数据 把数据导入数据库
		perr1 := Parsexml(des, "JZ_"+Diqu+"_"+fname, msg)
		if perr1 != nil {
			return perr1
		}

	}

	return nil
}

func Parsexml(pwdfname string, fname string, msg Message) error {
	log.Println("把原始数据xml插入数据库", pwdfname)

	//新增结算数据消息包
	err := JieSuanMessageInset(msg)
	if err != nil {
		log.Println("新增结算数据消息包 error")
		return err
	}

	//新增结算数据明细
	mxerr := JieSuanMessageMxInset(msg)
	if mxerr != nil {
		log.Println("新增结算数据明细 error")
		return mxerr
	}
	//新增结算处理
	clerr := JieSuanMessageChuliInset(msg)
	if clerr != nil {
		log.Println("新增结算处理记录 error ")
		return clerr
	}

	//把解析过的原始交易消息xml移动文件夹
	x1 := pwdfname
	x2 := "../centeryuanshixmlparsed/" + fname
	mverr := common.MoveFile(x1, x2)
	if mverr != nil {
		return errors.New("移动解析过的原始交易消息xml 失败")
	}
	return nil
}

// 获取xml文件msg的md5码
func GetFileMd5(filename string) string {
	// 文件全路径名
	path := "../centerYuanshi/" + filename
	pFile, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	log.Println("成功获取md5")
	return hex.EncodeToString(md5h.Sum(nil))
}
