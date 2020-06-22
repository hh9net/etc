package centerServer

import (
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

//模拟联网中心发送 记账数据、争议数据、清分数据
func CenterClient() {
	tiker := time.NewTicker(time.Second * 5)
	for {
		log.Println("执行 模拟联网中心发送 记账数据、争议数据、清分数据", <-tiker.C)

		//发送记账包数据
		err1 := SendKeepAccount()
		if err1 != nil {
			log.Println("发送记账包数据 error", err1)
		}

		//发送争议包数据
		err2 := SendDispute()
		if err2 != nil {
			log.Println("发送争议包数据 error", err2)
		}

		//发送清分包数据
		//err3 := SendClearing()
		//if err3 != nil {
		//	log.Println("发送清分包数据 error ", err3)
		//}

		//发送应答包
	}
}

//发送记账包
func SendKeepAccount() error {
	//1、扫描文件夹获取数据
	log.Println("执行发送记账包")
	for {
		pwd := "../centerkeepaccount/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		if len(fileInfoList) == 0 {
			return nil
		}

		for i := range fileInfoList {
			//判断文件的	前缀名
			if strings.HasPrefix(fileInfoList[i].Name(), "JZB") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())

				//压缩xml文件
				zjzerr := CenterZipxml(fileInfoList[i].Name(), "jz")
				if zjzerr != nil {
					return errors.New("记账包压缩失败")
				}
				log.Println("记账包压缩成功，可以移到已压缩文件夹 ：centerCompressed")
				//移动xml文件
				s1 := "../centerkeepaccount/" + fileInfoList[i].Name()
				s2 := "../centerCompressed/" + fileInfoList[i].Name()
				mverr := MoveFile(s1, s2)
				if mverr != nil {
					return errors.New("记账包xml 移动失败")
				}
				log.Println("记账包xml文件移动成功，可以进行解析xml，准备好发送报文")

				//		解析xml文件  获取数据
				sendStru, perr := ParsingXMLFiles(fileInfoList[i].Name())
				if perr != nil {
					return perr
				}

				//发送
				jzserr := Sendxml(sendStru)
				if jzserr != nil {
					return jzserr
				}
			}
		}
	}
}

//发送争议包
func SendDispute() error {
	//1、扫描文件夹获取数据
	log.Println("执行发送争议包")
	//tiker := time.NewTicker(time.Second * 5)
	for {
		//log.Println("执行线程2", <-tiker.C)
		pwd := "../centerdispute/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		if len(fileInfoList) == 0 {
			return nil
		}
		for i := range fileInfoList {
			//判断文件的	前缀名

			if strings.HasPrefix(fileInfoList[i].Name(), "ZYB") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())

				//压缩xml文件
				zipzyerr := CenterZipxml(fileInfoList[i].Name(), "zy")
				if zipzyerr != nil {
					return errors.New("争议包xml 压缩失败")
				}
				log.Println("争议包压缩成功")
				//MoveFile
				s1 := "../centerdispute/" + fileInfoList[i].Name()
				s2 := "../centerCompressed/" + fileInfoList[i].Name()

				mverr := MoveFile(s1, s2)
				if mverr != nil {
					return errors.New("争议包xml 移动失败")
				}
				log.Println("争议包xml文件移动成功")

				//		解析xml文件  获取数据
				sendStru, perr := ParsingXMLFiles(fileInfoList[i].Name())
				if perr != nil {
					return perr
				}

				//发送
				zyserr := Sendxml(sendStru)
				if zyserr != nil {
					return zyserr
				}
			}
		}
	}
}

//发送清分包
func SendClearing() error {
	//1、扫描文件夹获取数据
	log.Println("执行发送争议包")
	//tiker := time.NewTicker(time.Second * 5)
	for {
		pwd := "../centerClearing/"
		fileInfoList, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("该文件夹下有文件的数量 ：", len(fileInfoList))
		if len(fileInfoList) == 0 {
			return nil
		}
		for i := range fileInfoList {
			//判断文件的	前缀名
			if strings.HasPrefix(fileInfoList[i].Name(), "QFB") {
				log.Println("打印当前文件或目录下的文件名", fileInfoList[i].Name())

				//压缩xml文件
				zipqferr := CenterZipxml(fileInfoList[i].Name(), "qf")
				if zipqferr != nil {
					return errors.New("清分包xml 压缩失败")
				}
				log.Println("清分包压缩成功")
				//MoveFile
				s1 := "../centerClearing/" + fileInfoList[i].Name()
				s2 := "../centerCompressed/" + fileInfoList[i].Name()

				mverr := MoveFile(s1, s2)
				if mverr != nil {
					return errors.New("清分包xml 移动失败")
				}
				log.Println("清分包xml文件移动 成功")

				//		解析xml文件  获取数据
				sendStru, perr := ParsingXMLFiles(fileInfoList[i].Name())
				if perr != nil {
					return perr
				}

				//发送
				qfserr := Sendxml(sendStru)
				if qfserr != nil {
					return qfserr
				}
			}
		}
	}
}

//解析xml文件
func ParsingXMLFiles(fname string) (*types.SendStru, error) {
	var sendStru types.SendStru
	//1、获取消息包序号Massageid
	fnstr := strings.Split(fname, "_")
	idstr := strings.Split(fnstr[1], ".")
	sendStru.Massageid = idstr[0] //20位

	//2、获取压缩文件大小
	lengthstr := fmt.Sprintf("%06d", GetFileSize(fname))
	sendStru.Xml_length = lengthstr
	log.Println("发送文件大小", lengthstr)

	//3、获取文件md5  "JZB_3301_00000000000000100094.xml"
	sendStru.Md5_str = GetXmlFileMd5(fname)
	if sendStru.Md5_str != "" {
		log.Println("文件md5：", sendStru.Md5_str)
	} else {
		log.Println("获取文件md5 error ")
		return nil, errors.New("获取文件md5 error")
	}

	//4、获得文件名
	sendStru.Xml_msgName = fname + ".lz77" //压缩文件名
	//
	log.Println("报文信息：", sendStru)

	////5、更新数据    根据 包号 更新记账消息包的【发送状态   发送中】
	//Mid, _ := strconv.Atoi(sendStru.Massageid)
	//err := storage.UpdateYuansjyxx(int64(Mid))
	//if err != nil {
	//	log.Println("根据 包号 更新原始交易消息包的发送状态  error: ", err)
	//
	//}

	return &sendStru, nil
}

func connsendFile(data []byte, fname string, conn *net.Conn, sendStru *types.SendStru) error {
	connect := *conn
	n, err := connect.Write(data)
	if err != nil {
		log.Printf("发送字节数%d，错误为：%s", n, err)
		return err
	}

	path := "../centerSendxmlzip/" + fname
	log.Println("发送文件 path:=", path)
	file, oserr := os.Open(path)
	if oserr != nil {
		log.Println("os.Open error", oserr)
		return oserr
	}
	log.Println("文件conn发送开始")
	defer file.Close()
	//获取文件大小
	DataLen, _ := strconv.Atoi(sendStru.Xml_length)
	var buff []byte
	if DataLen > 0 {
		buff = make([]byte, DataLen)
	}
	log.Println("联网中心要发送的压缩文件的大小", DataLen)
	total := 0
	for {
		size, rerr := file.Read(buff)
		log.Println("读取文件的大小为：", size, rerr)
		if rerr != nil && rerr != io.EOF {
			log.Println("file.Read err:", rerr)
			break
		}
		if rerr == io.EOF {
			log.Println("文件读取完毕")
			log.Println("文件长度", total)
			break
		}

		_, werr := connect.Write(buff[:size])
		if werr != nil {
			log.Println("err", werr)
			return werr
		}
		total += size
		log.Println("已发送文件的大小为：", total)
		if total == DataLen {
			log.Println("发送报文到 通行宝  成功")
			return nil
		}
	}
	return nil
}

//发送
func Sendxml(sendStru *types.SendStru /*, conn *net.Conn*/) error {
	//连接通行包
	address := ListeningAddressConfigInit()
	Address := address.ListeningAddressIp + ":" + address.ListeningAddressPort
	//Dial

	conn, derr := net.Dial("tcp", Address)
	if derr != nil {
		log.Println("Dial 通行宝 error ", derr)
		return derr
	}
	if conn != nil {
		log.Println("Dial 通行宝 成功")
	}
	//把报文写给服务端
	data := []byte(sendStru.Massageid)
	length := []byte(sendStru.Xml_length)
	Md5 := []byte(sendStru.Md5_str)
	data = append(data, length...)
	data = append(data, Md5...)
	//发送
	err := connsendFile(data, sendStru.Xml_msgName, &conn, sendStru)
	if err != nil {
		log.Println("connsendFile error:", err)
		return err
	}

	//发送成功  更新 发送状态 发送时间 发送成功后消息包的文件路径
	Mid, _ := strconv.Atoi(sendStru.Massageid)
	//消息包发送成功更新 发送状态 发送时间 发送成功后消息包的文件路径

	log.Printf("消息包：%d 发送成功 ", Mid)
	//log.Printf("发送后收到的接收应答 ")

	//读取通行宝的即时应答
	buf := make([]byte, 100)
	n, err2 := (conn).Read(buf)
	if err2 != nil && err2 != io.EOF {
		log.Println("conn.Read err = ", err2)
		return err2
	}
	if err2 == io.EOF {
		log.Println("conn.Read 读取联网中心的应答结束", err2)
		//对联网中心的接收应答处理

		(conn).Close()
		return nil
	}
	log.Println("收到的即时应答 :", string(buf[:n]))

	return nil

	//err1, DBsj := storage.SendedUpdateYuansjyxx(int64(Mid), sendStru.Xml_msgName)
	//if err1 != nil {
	//	log.Println("storage.SendedUpdateYuansjyxx  error:", err1)
	//}
	////TCP发送记录
	//has, serr, count := storage.GetTcpSendRecord(strconv.Itoa(Mid))
	//if serr != nil {
	//	log.Println("查询TCP发送记录错误")
	//}
	//var record storage.BJsTcpqqjl
	//
	//if has == false {
	//	log.Println("此msgid属于第一次发送")
	//	record.FVcXiaoxxh = strconv.Itoa(Mid)
	//	record.FDtZuixsj = DBsj
	//	record.FNbChongfcs = 0
	//	record.FNbFasz = 1
	//	record.FNbMd5 = sendStru.Md5_str
	//	err2 := storage.TcpSendRecordInsert(record)
	//	if err2 != nil {
	//		log.Println("storage.TcpSendRecordInsert error:", err2)
	//	}
	//}
	//if has == true {
	//	log.Println("此msgid 已经发送")
	//	record.FVcXiaoxxh = strconv.Itoa(Mid)
	//	record.FNbChongfcs = count + 1
	//	err3 := storage.TcpSendRecordUpdate(record)
	//	if err3 != nil {
	//		log.Println("storage.TcpSendRecordUpdate error:", err3)
	//	}
	//}
	//log.Println("TCP发送记录 插入完成")

}

//获取文件大小
func GetFileSize(fname string) int64 {

	path := "../centerSendxmlzip/" + fname + ".lz77"
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("获取文件大小 error ", err)
	}
	//文件大小
	//log.Println("文件大小", fileInfo.Size()) //返回的是字节
	return fileInfo.Size()
}

//移动文件
func MoveFile(src string, des string) error {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Fatalln("移动文件错误", err)
		return err
	}
	log.Printf("移动文件%s to %s 成功", src, des)
	return nil
}

// 获取xml文件msg的md5码
func GetXmlFileMd5(filename string) string {
	// 文件全路径名
	path := "../centerCompressed/" + filename
	pFile, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return strings.ToUpper(hex.EncodeToString(md5h.Sum(nil)))
}
