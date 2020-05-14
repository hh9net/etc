package common

import (
	"log"
	"time"
)

//时间格式化
//2006-01-02 15:04:05
func DateTimeFormat() string {
	var dateformat string
	datetime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无输出
	log.Println("now is :", datetime)
	return dateformat
}

//2006-01-02T15:04:05
func DateFormat() string {
	datetime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出
	b := []byte(datetime)
	b[10] = 'T'
	log.Println(string(b))
	return string(b)
}

// 2006-01-02
func TodayFormat() string {
	today := time.Now().Format("2006-01-02")
	log.Println("now is :", today)
	return today
}
