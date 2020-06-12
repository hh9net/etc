package service

import (
	"CenterSettlement-go/client"
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	"CenterSettlement-go/lz77zip"
	"CenterSettlement-go/types"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

//接收联网中心发来数据包
func Receive() {
	log.Println("执行线程4 接收联网中心发送文件")
	//监听联网中心数据端口 "127.0.0.1:8809"
	address := conf.ListeningAddressConfigInit()
	Address := address.ListeningAddressIp + ":" + address.ListeningAddressPort
	listen, lerr := net.Listen("tcp", Address)
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
		go HandleTask(&connect)
	}
	defer listen.Close()
}

//处理任务
func HandleTask(conn *net.Conn) {
	go HandleMessage(conn)
}

//线程3  接收数据包
func HandleMessage(conn *net.Conn) {

	//把读到的数据 以文件记录  获取xml name
	Filename, err := Save(conn)
	if err != nil {
		log.Println("记录文件时 Save(conn) 错误", err)
		return
	}

	//解析文件
	ParsingFile(Filename)

	defer (*conn).Close()

}

//保存联网中心发来的数据
func Save(conn *net.Conn) (string, error) {
	//var fileName string
	//接受客户端发来的要传文件的文件信息data
	buffer := make([]byte, 58)
	d, e := (*conn).Read(buffer)
	log.Println("描述文件的数据data：", string(buffer[:d]), e)
	fileNameid := string(buffer[:20])

	msglength := string(buffer[20:26])
	log.Println("消息包msglength", msglength)

	msgmd5 := string(buffer[26:58])
	log.Println("消息包msgmd5：", msgmd5)

	//msgid := string(buffer[:20])
	//log.Println("消息包Massageid", msgid)

	//接收文件，保存文件，即使应答
	rferr, fname := RevFile(fileNameid, conn, msglength)
	if rferr != nil {
		log.Println(" 处理接收文件，保存文件，即使应答  错误：", rferr)
		return fname, rferr
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
			return "", mxerr
		}

		//移动xmlzip
		xz1 := "../centeryuanshixmlzip/" + fileNameid + ".xml.lz77"
		xz2 := "../centeryuanshixmlzip/errorxml/" + fileNameid + ".xml.lz77"
		mxzerr := common.MoveFile(xz1, xz2)
		if mxzerr != nil {
			log.Println("移动errorxml 失败")
			return "", mxzerr
		}
	}

	return fname, nil
}

//接收文件 保存为文件 即时应答
func RevFile(fileNameid string, conn *net.Conn, msglength string) (error, string) {
	//创建文件
	fs, err := os.Create("CenterSettlement-go/receivezipfile/" + fileNameid + ".xml.lz77")
	defer fs.Close()
	if err != nil {
		log.Println("os.Create err =", err)
		return err, ""
	}
	//获取文件大小
	DataLen, _ := strconv.Atoi(msglength) //压缩文件长度
	//动态分配切片大小
	//一次性读取文件
	//一次性写完
	//if DataLen > 0 {
	//	data = make([]byte, DataLen)
	//	if _, err := io.ReadFull(*conn, data); err != nil {
	//		log.Println("read msg data error ", err)
	//		return err
	//	}
	//}
	var buff []byte
	if DataLen > 0 {
		buff = make([]byte, DataLen)
	}
	log.Println("要接收的文件的长度：", DataLen)
	//读取内容
	n, rerr := (*conn).Read(buff)
	if rerr != nil {
		log.Println("conn.Read err:", rerr)
		return rerr, ""
	}

	size, werr := fs.Write(buff[:n])
	log.Println("写入文件的大小", size, rerr)
	if werr != nil {
		log.Println("写入文件时的错误为：", rerr)
		return werr, ""
	}
	//如果实际总接受字节数与客户端给的要传输字节数相等，说明传输完毕
	if size == DataLen {
		log.Println("文件消息包 接受成功,共：", size, "个字节")

		//即时应答
		var replyStru types.ReplyStru
		//replyStru.Massageid = string(buf[:20])
		replyStru.Massageid = fileNameid

		replyStru.Result = "1"
		m := []byte(replyStru.Massageid)
		r := []byte(replyStru.Result)
		d := append(m, r...)
		InstantResponse(d, conn)
	}

	return nil, fileNameid + ".xml.lz77"
}

//应答
func InstantResponse(d []byte, conn *net.Conn) {
	// 返回接收成功
	_, err := (*conn).Write(d)
	if err != nil {
		log.Println("联网中心 conn.Write 错误")
	}

	(*conn).Close()
}

func CheckFile(msgmd5, fileNameid string) error {
	pwd := "CenterSettlement-go/receivezipfile/"
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
			zerr := lz77zip.UnZipLz77(fileInfoList[i].Name())
			if zerr != nil {
				log.Println(" 解压缩文件失败", zerr)
				return zerr //此处解压失败，回头处理
			}
			if zerr == nil {
				//解压缩成功 移动zipxml文件
				s := "CenterSettlement-go/receivezipfile/" + fileInfoList[i].Name()
				des := "CenterSettlement-go/receiveUnzipsucceed/" + fileInfoList[i].Name()
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

// 获取xml文件msg的md5码
func GetFileMd5(filename string) string {
	// 文件全路径名
	path := "CenterSettlement-go/receivexml/" + filename
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

//解析文件
func ParsingFile(Filename string) {
	//1、扫描文件夹，改名字
	log.Println("去执行线程4吧，本线程3的工作完成了！")
}
