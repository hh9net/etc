package service

import (
	"log"
	"os"
)

//文件处理
func GetFileSize(fname string) int64 {
	path := "../sendzipxml/" + fname
	fileInfo, _ := os.Stat(path)
	//文件大小
	log.Println("文件大小", fileInfo.Size()) //返回的是字节
	return fileInfo.Size()
}
