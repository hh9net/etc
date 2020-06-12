package common

import (
	"log"
	"os"
)

//文件处理
func GetFileSize(fname string) int64 {
	//path := "CenterSettlement-go/sendzipxml/" + fname + ".lz77"

	path := "CenterSettlement-go/sendzipxml/" + fname + ".lz77"
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("获取文件大小 error ", err)
	}
	//文件大小
	//log.Println("文件大小", fileInfo.Size()) //返回的是字节
	return fileInfo.Size()
}

//改文件名字
func FileRename(src string, des string) error {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Fatalln("移动文件错误", err)
		return err
	}
	log.Printf("移动文件%s to %s 成功", src, des)
	return nil
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
