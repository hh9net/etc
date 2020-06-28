package main

import (
	"CenterSettlement-go/common"
	"CenterSettlement-go/conf"
	"CenterSettlement-go/database"
	"CenterSettlement-go/service"
	commonUtils "CenterSettlement-go/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	// 打开文件锁 处理进程互斥
	lock, e := common.Create("./LockFile.txt")
	if e != nil {
		log.Println("打开文件锁 错误 ", e) // handle error
	}
	defer lock.Release()

	// 尝试独占文件锁
	e = lock.Lock()
	if e != nil {
		log.Println("独占文件锁 错误 ", e) // handle error
		os.Exit(1)
	}
	defer lock.Unlock()

	// 日志初始化
	conf := conf.LogConfigInit() //日志配置
	commonUtils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogmaxAge)*time.Hour,
		time.Duration(conf.LogrotationTime)*time.Hour)
	database.DBInit() //连接数据库 初始化为全局变量
	//goroutine1
	//go service.HandleGeneratexml()
	//goroutine2
	//go service.HandleSendXml()
	//goroutine4
	go service.AnalyzeDataPakage()
	//goroutine3
	//service.Receive()

	for {
		tiker := time.NewTicker(time.Second * 30)
		for {
			log.Println("执行主go程 ", <-tiker.C)
		}
	}
}

//
