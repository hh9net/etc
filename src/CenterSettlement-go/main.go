package main

import (
	"CenterSettlement-go/conf"
	"CenterSettlement-go/database"
	"CenterSettlement-go/service"
	commonUtils "CenterSettlement-go/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	// 日志初始化
	conf := conf.LogConfigInit() //日志配置
	commonUtils.InitLogrus(conf.LogPath, conf.LogFileName, time.Duration(24*conf.LogmaxAge)*time.Hour, time.Duration(conf.LogrotationTime)*time.Hour)
	database.DBInit()
	log.Println(database.XormClient)
	//goroutine1
	go service.Generatexml()
	//goroutine2
	go service.HandleSendXml()
	//goroutine4
	go service.AnalyzeDataPakage()
	//goroutine3
	service.Receive()
}
