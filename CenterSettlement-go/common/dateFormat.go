package common

import (
	"fmt"
	"log"
	"time"
)

//时间格式化
//2006-01-02 15:04:05
func DateTimeformat() string {
	var dateformat string
	datetime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无输出
	log.Println("now is :", datetime)
	return dateformat
}

//2006-01-02T15:04:05
func Dateformat() string {
	var dateformat string

	today := time.Now().Format("2006-01-02")
	datetime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出

	fmt.Println("now is :", today)
	fmt.Println("now is :", datetime)
	return dateformat
}

// 2006-01-02
func Todayformat() string {
	var dateformat string
	t := time.Now()
	y, m, d := t.Date()
	today := time.Now().Format("2006-01-02")
	datetime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出

	fmt.Println("time is : ", t)
	fmt.Println("y m d is : ", y, m, d)
	fmt.Println("now is :", today)
	fmt.Println("now is :", datetime)
	return dateformat
}
