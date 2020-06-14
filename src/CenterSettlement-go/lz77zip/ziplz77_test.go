package lz77zip

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"testing"
)

//go centeryuanshixmlzip
func TestZipLz77(t *testing.T) {
	f := "CZ_3101_00000000000000100136.xml"
	err := ZipLz77(f)
	if err != nil {
		log.Fatalln(err)
	}
}

//go unzip
func TestUnZipLz77(t *testing.T) {
	f := "00000000000000000001.xml.lz77"
	err := UnZipLz77(f)
	if err != nil {
		log.Fatalln(err)
	}
}

//交叉测试 压缩
func TestGetFileMd5Test(t *testing.T) {
	//压缩：1
	f1 := "Cgo_CZ_3101_00000000000000100136.xml.lz77"
	s1 := GetFileMd5Test(f1, 1)

	log.Println("s1md5", s1)
	f2 := "CZ_3101_00000000000000100136.xml.lz77"
	s2 := GetFileMd5Test(f2, 1)
	log.Println("s2md5", s2)
	if s1 == s2 {
		log.Println("重写成功")
	} else {
		log.Println("失败成功")
	}
}

//交叉测试 解压
func TestGetFilexmlMd5Test(t *testing.T) {
	//解压 2
	f1 := "Cgo_CZ_3101_00000000000000100136.xml"
	s1 := GetFileMd5Test(f1, 2)
	log.Println("s1md5", s1)
	f2 := "CgoUnZip_CZ_3101_00000000000000100136.xml"
	s2 := GetFileMd5Test(f2, 2)
	log.Println("s2md5", s2)
	if s1 == s2 {
		log.Println("重写成功")
	} else {
		log.Println("重写失败")
	}
}

// 获取xml文件msg的md5码
func GetFileMd5Test(filename string, fs int) string {
	// 文件全路径名
	//path := "CenterSettlement-go/compressed_xml/" + filename
	var path string
	if fs == 1 {
		path = "../sendzipxml/" + filename //压缩
	}
	if fs == 2 {
		path = "../receivexml/" + filename //解压
	}
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

func TestGetFileTest(t *testing.T) {
	//解压 2
	f1 := "../centerServer/00000000000000000001.xml.lz77"
	s1 := GetFileMd5Test(f1, 2)
	log.Println("s1md5", s1)
	f2 := "CgoUnZip_CZ_3101_00000000000000100136.xml"
	s2 := GetFileMd5Test(f2, 2)
	log.Println("s2md5", s2)
	if s1 == s2 {
		log.Println("重写成功")
	} else {
		log.Println("重写失败")
	}
}
