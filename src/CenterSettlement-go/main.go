package main

import (
	"CenterSettlement-go/conf"
	"CenterSettlement-go/database"
	"CenterSettlement-go/service"
	commonUtils "CenterSettlement-go/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

//
//func test(lib *lib7z.Lz77so) {
//	file := "generatexml/CZ_3201_00000000000000100048.xml"
//
//	for i := 0; i <= 5; i++ {
//		filenew := fmt.Sprintf("%d.xml.lz77", i)
//		fileextra := fmt.Sprintf("%d.xml", i)
//		lib.Comresslz77(file, filenew)
//		lib.Depresslz77(filenew, fileextra)
//
//	}
//}

func main() {
	//var libtest lib7z.Lz77so
	// 日志初始化
	conf := conf.LogConfigInit() //日志配置
	commonUtils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogmaxAge)*time.Hour,
		time.Duration(conf.LogrotationTime)*time.Hour)
	database.DBInit() //连接数据库 初始化为全局变量
	//goroutine1
	go service.HandleGeneratexml()
	//goroutine2
	go service.HandleSendXml()
	//goroutine4
	go service.AnalyzeDataPakage()
	//goroutine3
	service.Receive()
	//主线程处理压缩与解压缩

	//for true {
	//	test(&libtest)
	//	time.Sleep(time.Second * 3)
	//}
	for {
		tiker := time.NewTicker(time.Second * 2)
		for {
			log.Println(<-tiker.C)
		}
	}
}
