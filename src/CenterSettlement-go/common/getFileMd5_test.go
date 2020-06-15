package common

import (
	"log"
	"testing"
)

func TestGetFileMd5(t *testing.T) {
	// 当前目录的csv配置文件为例00000000000000109139.xml
	fileName := "00000000000000010796.xml"
	md5Val := GetFileMd5(fileName)
	log.Println("配置文件的md5值：", md5Val)
}
