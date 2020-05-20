package client

import (
	"log"
	"os"
)

//如果发送成功，把文件移动到指定的文件夹
func MoveFile(src string, des string) {
	//err := os.Rename("./a", "/tmp/a")
	err := os.Rename(src, des)
	if err != nil {
		log.Fatalln("移动文件错误", err)
		return
	}

}

//删除文件
func DelFile(src string) {
	//"./1.txt"
	del := os.Remove(src)
	if del != nil {
		log.Println(del)
	}
}

//删除指定path下的所有文件
func DelAllFile(src string) {
	//"./dir"
	delDir := os.RemoveAll(src)
	if delDir != nil {
		log.Println(delDir)
	}
}
