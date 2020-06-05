package common

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// 获取xml文件msg的md5码
func GetFileMd5(filename string) string {
	// 文件全路径名
	path := "CenterSettlement-go/generatexml/" + filename
	pFile, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败，filename=%v, err=%v", filename, err)
		return ""
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil))
}
