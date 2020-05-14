package client

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

//如果发送成功，把文件移动到指定的文件夹
func MoveFile() {
	//打开文件，返回文件指针
	file, err := os.Open("../genetatexml/CZ_3201_00000000000000999999.xml")
	if err != nil {
		log.Println(err)
	}
	log.Println(file)
	defer file.Close()

	//以读写方式打开文件，如果不存在，则创建
	file2, err := os.OpenFile("./2.txt", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		log.Println(err)
	}
	log.Println(file2)
	defer file2.Close()

	//创建文件
	//Create函数也是调用的OpenFile
	file3, err := os.Create("./3.txt")
	if err != nil {
		log.Println(err)
	}
	log.Println(file3)
	defer file3.Close()

	//读取文件内容
	file4, err := os.Open("./1.txt")
	if err != nil {
		log.Println(err)
	}
	//创建byte的slice用于接收文件读取数据
	buf := make([]byte, 1024)
	//循环读取buf
	for {
		//Read函数会改变文件当前偏移量
		n, _ := file4.Read(buf)
		//读取字节数为0时跳出循环
		if n == 0 {
			break
		}
		fmt.Println(string(buf))
	}
	_ = file4.Close()
	//读取文件内容
	file5, err := os.Open("./1.txt")
	if err != nil {
		fmt.Println(err)
	}
	buf2 := make([]byte, 1024)
	ix := 0
	for {
		//ReadAt从指定的偏移量开始读取，不会改变文件偏移量
		n, _ := file5.ReadAt(buf2, int64(ix))
		ix = ix + n
		if n == 0 {
			break
		}
		fmt.Println(string(buf2))
	}
	_ = file5.Close()
	//写入文件
	file6, err := os.Create("./4.txt")
	if err != nil {
		fmt.Println(err)
	}
	data := "我是数据\r\n"
	for i := 0; i < 10; i++ {
		//写入byte的slice数据
		_, _ = file6.Write([]byte(data))
		//写入字符串
		_, _ = file6.WriteString(data)
	}
	_ = file6.Close()
	//写入文件
	file7, err := os.Create("./5.txt")
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 10; i++ {
		//按指定偏移量写入数据
		ix := i * 64
		_, err := file7.WriteAt([]byte("我是数据"+strconv.Itoa(i)+"\r\n"), int64(ix))
		if err != nil {
			fmt.Println(err)
		}
	}
	ferr := file7.Close()
	if ferr != nil {
		fmt.Println(ferr)
	}
	//删除文件
	del := os.Remove("./1.txt")
	if del != nil {
		fmt.Println(del)
	}
	//删除指定path下的所有文件
	delDir := os.RemoveAll("./dir")
	if delDir != nil {
		fmt.Println(delDir)
	}
}
