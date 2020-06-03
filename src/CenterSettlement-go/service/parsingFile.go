package service

import (
	"log"
	"os"
)

//文件处理
func GetFileSize(fname string) int64 {
	path := "../sendzipxml/" + fname + ".lz77"
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("获取文件大小 error ", err)
	}
	//文件大小
	log.Println("文件大小", fileInfo.Size()) //返回的是字节
	return fileInfo.Size()
}
