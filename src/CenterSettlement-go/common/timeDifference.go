package common

import (
	"fmt"
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

func FDate() {
	a, _ := time.Parse("2006-01-02", "2018-09-01")
	b, _ := time.Parse("2006-01-02", "2018-09-02")
	d := a.Sub(b)
	fmt.Println(d.Hours() / 24)
}

func getHourDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}

func getMDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 60
		return hour
	} else {
		return hour
	}
}

func getSDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff
		return hour
	} else {
		return hour
	}
}
