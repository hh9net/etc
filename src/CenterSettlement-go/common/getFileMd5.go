package common

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// 获取xml文件msg的md5码
func GetFileMd5(filename string) string {
	// 文件全路径名
	path := "CenterSettlement-go/compressed_xml/" + filename
	//path := "../compressed_xml/" + filename  //test

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
