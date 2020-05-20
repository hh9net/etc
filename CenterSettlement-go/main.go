package main

import (
	"CenterSettlement-go/database"
	"CenterSettlement-go/service"
	commonUtils "CenterSettlement-go/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

//日志配置
var (
	LogmaxAge       = 365 //日志最大保存时间（天）
	LogrotationTime = 24  //日志切割时间间隔（小时）
	LogPath         = "./log"
	LogFileName     = "CenterSettlement.log"
)

func main() {
	// 日志初始化
	commonUtils.InitLogrus(LogPath, LogFileName, time.Duration(24*LogmaxAge)*time.Hour, time.Duration(LogrotationTime)*time.Hour)
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
