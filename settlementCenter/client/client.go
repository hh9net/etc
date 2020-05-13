package client

import (
	"compress/flate"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func NewConn() net.Conn {
	//1、客户端主动连接服务器 http://127.0.0.1:8808
	conn, err := net.Dial("tcp", "127.0.0.1:8808")
	if err != nil {
		log.Fatal("Dial 联网中心 err:", err)
		return nil
	}
	log.Fatal("tcp连接成功")
	defer conn.Close() //条件反射出来 延迟关闭
	return conn
}

func Sendxml() {
	conn := NewConn()
	//2、模拟浏览器，组织一个最简单的请求报文。包含请求行，请求头，空行即可。
	//requestHttpHeader := "GET /ysjyxx0000011.xml HTTP/1.1\r\nHost:127.0.0.1:8808\r\n\r\n"
	requestHttpHeader := "GET /ysjyxx0000011.xml HTTP/1.1\r\nHost:127.0.0.1:8808\r\n\r\n"
	// 准备原始记录的xml数据包的请求包
	//3、给服务器发送请求报文
	conn.Write([]byte(requestHttpHeader))

	//4、读取 服务器回复 响应报文
	// 读取响应缓冲区
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read err:", err)
		return
	}
	//5、打印观察
	fmt.Printf("#\n%s#", string(buf[:n]))
}

func ReadAndCompress(filename string) {
	filepath := "../generatexml/" + filename

	fbuf := make([]byte, 1024*4)

	file, fErr := os.OpenFile(filepath, os.O_RDWR, 0777)
	if fErr != nil {
		log.Fatal("Open generatexml err:", fErr)
		return
	}
	defer file.Close()
	// 向缓存区写入数据
	Size, readErr := file.Read(fbuf)
	if readErr != nil {
		if readErr == io.EOF {
			log.Fatal("EOF", readErr)
		} else {
			log.Fatal("Read:", readErr)
		}
		return
	}
	// 一个缓冲区压缩的内容
	//buf := bytes.NewBuffer(nil)
	fn := []byte(filename)
	fname := string(fn[:len(fn)-3])
	//创建文件
	//fw, f_werr := os.Create("../sendzipxml/"+fname + "lz77")
	fw, f_werr := os.Create("../sendzipxml/" + fname + "xml")
	if f_werr != nil {
		log.Fatal("Read:", f_werr)
	}
	defer fw.Close()
	// 创建一个flate.Writer，压缩级别最好
	flateWrite, err := flate.NewWriter(fw, flate.BestCompression)
	if err != nil {
		log.Fatalln(err)
	}
	defer flateWrite.Close()
	// 写入待压缩内容
	flateWrite.Write(fbuf[:Size])
	flateWrite.Flush()
	//flateWrite.
	fmt.Println(fbuf)
}

//XML文件 按照 LZ77 方式打包
func Compress() {
	myfolder := `../generatexml`
	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println(file.Name())
			ReadAndCompress(file.Name())
		}
	}
}
