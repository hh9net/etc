package common

import (
	"log"
	"testing"
)

func TestGetFileMd5(t *testing.T) {
	// 当前目录的csv配置文件为例
	fileName := "JZ_3301_00000000000000100094.xml.lz77"
	md5Val := GetFileMd5(fileName)
	log.Println("配置文件的md5值", md5Val)
}
