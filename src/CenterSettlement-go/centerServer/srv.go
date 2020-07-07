package centerServer

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/common"
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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
	log.Println("监听客户端 ", tcpAddr)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println(err.Error())
	}
	//defer listener.Close()

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
	//校验文件md5
	log.Println("执行 CheckFile校验文件")
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
			log.Println("移动CheckFile失败 的xml error")
			return
		}

		//移动xmlzip
		xz1 := "../centeryuanshixmlzip/" + fileNameid + ".xml.lz77"
		xz2 := "../centeryuanshixmlzip/errorxmlzip/" + fileNameid + ".xml.lz77"
		mxzerr := common.MoveFile(xz1, xz2)
		if mxzerr != nil {
			log.Println("移动CheckFile失败 的zipxml error\"")
			return
		}
	}
}

func CheckFile(msgmd5, fileNameid string) error {
	pwd := "./"
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CheckFile 时 该联网文件夹下有文件的数量 ：", len(fileInfoList))

	for i := range fileInfoList {
		//判断文件的结尾名
		if strings.HasSuffix(fileInfoList[i].Name(), ".lz77") {
			log.Println("执行 CheckFile 打印当前目录下的压缩文件的名字", fileInfoList[i].Name())

			//解压缩
			zerr := UnZipLz77(fileInfoList[i].Name())
			if zerr != nil {
				log.Println(" 解压缩文件失败", zerr)
				return zerr
			}
			if zerr == nil {
				//移动zipxml文件
				log.Println("原始交易xml消息解压成功,移动zipxml文件")
				s := "./" + fileInfoList[i].Name()
				des := "../centeryuanshixmlzip/" + fileInfoList[i].Name()
				merr := client.MoveFile(s, des)
				if merr != nil {
					log.Println("client.MoveFile err : ", merr)
					return merr
				}
			}
			log.Println("校验文件Md5")
			fstr := strings.Split(fileInfoList[i].Name(), ".xml.lz77")

			log.Println(fstr[0], fileNameid)
			if fileNameid == fstr[0] {
				//校验文件 校验其md5

				fmd5 := GetFileMd5(fileNameid + ".xml")

				if fmd5 == msgmd5 {
					log.Println("文件md5一致")
					log.Println("此文件可以进行 解析原始交易消息文件")
					//解析文件 导入数据库
					//ParsingYSXMLFile()
					return nil
				} else {
					log.Println("new Md5:", fmd5)
					log.Println("发送过来的md5：", msgmd5)
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

	fs, err := os.Create("../centerServer/" + fileNameid + ".xml.lz77")
	defer fs.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return err
	}
	//获取文件大小
	DataLen, _ := strconv.Atoi(msglength) //压缩文件长度
	var buff []byte
	if DataLen > 0 {
		buff = make([]byte, DataLen)
	}
	log.Println("要接收的文件的长度：", DataLen, len(buff))
	i := DataLen

	for {
		//读取内容
		n, rerr := (conn).Read(buff)
		i = i - n
		if rerr != io.EOF && rerr != nil {
			log.Println("conn.Read 读文件出错:", rerr)
			return rerr
		} //读文件出错

		log.Printf("本次读出文件的大小 %d 还要读出%d 错误为：%s", n, i, rerr)

		if rerr == io.EOF {
			log.Println("文件读结束了", rerr)
			return nil
		} //文件读结束了

		if n == 0 {
			log.Println("n=0,文件读结束了")
			return nil
		}

		//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
		if i < 0 {
			//减少最后一个不需要的字节
			_, werr := fs.Write(buff[:n-1])
			if werr != nil {
				log.Println("写入文件时的错误为：", werr)
				return werr
			}
			log.Println("文件消息包 接受成功,共：", DataLen, "个字节")
			//即时应答
			var replyStru types.ReplyStru
			replyStru.Massageid = fileNameid
			replyStru.Result = "1"
			m := []byte(replyStru.Massageid)
			r := []byte(replyStru.Result)
			d := append(m, r...)
			InstantResponse(d, conn)
			return nil
		}

		//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
		if i == 0 {
			//减少最后一个不需要的字节
			_, werr := fs.Write(buff[:n])
			if werr != nil {
				log.Println("写入文件时的错误为：", werr)
				return werr
			}
			log.Println("文件消息包 接受成功,共：", DataLen, "个字节")
			//即时应答
			var replyStru types.ReplyStru
			replyStru.Massageid = fileNameid
			replyStru.Result = "1"
			m := []byte(replyStru.Massageid)
			r := []byte(replyStru.Result)
			d := append(m, r...)
			InstantResponse(d, conn)
			return nil
		}
		size, werr := fs.Write(buff[:n])
		if werr != nil {
			log.Println("写入文件时的错误为：", werr, size)
			return werr
		}

	}
	//return nil
}
func Handle() {
	//处理文件
	//tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("++++++++++处理文件+++++++++++")
		//log.Println("扫描centerYuanshi文件夹,解析文件、数据入库", <-tiker.C)
		herr := HandleFile()
		if herr == errors.New("该centerYuanshi 文件夹下没有需要解析的文件") {
			log.Println("该centerYuanshi 文件夹下没有需要解析的文件")
			return
		}
		if herr != nil {
			log.Println("HandleFile 执行失败 ")
			return
		}

		if herr == nil {
			log.Println("HandleFile 执行  ok ")
		}
	}
}
func HandleFile() error {
	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("扫描centerYuanshi文件夹,解析文件、数据入库", <-tiker.C)

		//扫描原始数据包
		pwd := "../centerYuanshi/"
		fileList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Println("该centerYuanshi文件夹下有xml文件的数量 ：", len(fileList))
		if len(fileList) == 0 {
			log.Println("该centerYuanshi 文件夹下没有需要解析的文件")
			return errors.New("该centerYuanshi 文件夹下没有需要解析的文件")
		}
		for i := range fileList {
			log.Println("该centerYuanshi文件夹需要数据入库的 xml 名字为:", fileList[i].Name()) //打印当前文件或目录下的文件或目录名

			//判断文件的结尾名
			if strings.HasSuffix(fileList[i].Name(), ".xml") {
				//解析xml文件
				content, rderr := ioutil.ReadFile("../centerYuanshi/" + fileList[i].Name())
				if rderr != nil {
					log.Println("读文件位置错误信息：", rderr)
					return rderr
				}

				//将xml文件转换为对象
				var msg Message
				unzerr := xml.Unmarshal(content, &msg)
				if unzerr != nil {
					log.Println("解析 receive文件夹中xml文件内容的错误信息：", unzerr)
					return unzerr
				}

				log.Println("msg:", msg.Header.MessageClass, msg.Header.MessageType, msg.Body.ContentType, msg.Header.MessageId, msg.Body.Transaction[0].ICCard.CardType)
				//原始交易数据

				// 数据导入数据库
				perr := ParsingYSXMLFile(msg, fileList[i].Name())
				if perr != nil {
					log.Println("数据导入数据库失败", perr)
					return perr
				}

			}

		}
		return nil
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
		rnerr := common.FileRename(s1, des)
		if rnerr != nil {
			log.Println("储值卡改文件名 error")
			return rnerr
		}
		log.Println("储值卡改文件名成功，数据插入数据库")
		//2、解析xml数据 把数据导入数据库
		perr, filename, errlx := Parsexml(des, "CZ_"+Diqu+"_"+fname, msg)
		switch errlx {
		case 0:
			log.Println(filename)
			log.Println("解析xml数据 把数据导入数据库 成功")
		case 1:
			log.Println(filename)
			log.Println("联网中心 新增结算数据 时 错误")
		case 2:
			log.Println(filename)
			log.Println("联网中心 新增结算数据明细 时 错误")
		case 3:
			log.Println(filename)
			log.Println("联网中心 新增结算处理数据 时 错误")
		case 4:
			log.Println(filename)
			log.Println("xml文件移到解析过的文件里 时 错误")
		}
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
		rnzerr := common.FileRename(s1, des)
		if rnzerr != nil {
			log.Println("记账卡改文件名 error")
			return rnzerr
		}

		//2、解析xml数据 把数据导入数据库
		perr1, filename, errlx := Parsexml(des, "JZ_"+Diqu+"_"+fname, msg)
		switch errlx {
		case 0:
			log.Println(filename)
			log.Println("解析xml数据 把数据导入数据库 成功")
		case 1:
			log.Println(filename)
			log.Println("联网中心 新增结算数据 时 错误")
		case 2:
			log.Println(filename)
			log.Println("联网中心 新增结算数据明细 时 错误")
		case 3:
			log.Println(filename)
			log.Println("联网中心 新增结算处理数据 时 错误")
		case 4:
			log.Println(filename)
			log.Println("xml文件移到解析过的文件里 时 错误")
		}
		if perr1 != nil {
			return perr1
		}
		return nil
	}

	return nil
}

func Parsexml(pwdfname string, fname string, msg Message) (error, string, int) {
	log.Println("把原始数据xml插入数据库，文件的路径以及名称为：", pwdfname)
	//新增结算数据消息包
	err := JieSuanMessageInset(msg)
	if err != nil {
		log.Println("新增结算数据消息包 error", err)
		return err, fname, 1
	}
	//
	////新增结算数据明细
	//mxerr := JieSuanMessageMxInset(msg)
	//if mxerr != nil {
	//	log.Println("新增结算数据明细 error", mxerr)
	//	return mxerr, fname, 2
	//}
	////新增结算处理
	//clerr := JieSuanMessageChuliInset(msg)
	//if clerr != nil {
	//	log.Println("新增结算处理记录 error ", clerr)
	//	return clerr, fname, 3
	//}

	//把解析过的原始交易消息xml移动文件夹
	x1 := pwdfname
	x2 := "../centeryuanshixmlparsed/" + fname
	mverr := common.MoveFile(x1, x2)
	if mverr != nil {
		log.Println("移动解析过的原始交易消息xml  error ", mverr)
		return errors.New("移动解析过的原始交易消息xml 失败"), fname, 4
	}
	log.Printf("此%s消息文件数据入库成功,xml文件移到解析过的文件里成功", fname)
	return nil, fname, 0
}

// 获取xml文件msg的md5码
func GetFileMd5(filename string) string {
	// 文件全路径名
	path := "../centerYuanshi/" + filename
	log.Println("GetFileMd5 path:", path)
	pFile, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)
	log.Println("成功获取此文件md5", filename)
	//str :=
	return strings.ToUpper(hex.EncodeToString(md5h.Sum(nil)))
}
