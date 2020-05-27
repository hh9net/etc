package common

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func TimeDifference(intime string, outtime string) {
	inTime, _ := time.Parse("2006-01-02 15:04:05", intime)
	outTime, _ := time.Parse("2006-01-02 15:04:05", outtime)
	t := outTime.Sub(inTime)

	log.Println("时间差", t)
	log.Println("时间差 小时", t.Hours(), "=========", t.Hours())
	log.Println("时间差  分", t.Minutes())
	log.Println("时间差  秒", t.Seconds())
}
