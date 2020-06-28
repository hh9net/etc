package common

import (
	"log"
	"time"
)

//时间格式化
//转换成时间格式： 2006-01-02 15:04:05
func DateTimeFormat(t time.Time) string {
	datetime := t.Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无输出
	//log.Println("DateTimeFormat :", datetime)
	return datetime
}

//2006-01-02 15:04:05
func DateTimeNowFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
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
	return time.Now().Format("2006-01-02")
}

//处理时间2017-07-07T11:33:53 转为 2017-07-07 11:33:53
func DataTimeFormatHandle(datetime string) string {
	b := []byte(datetime)
	b[10] = ' '
	log.Println(datetime, "to", string(b))
	return string(b)
}

//处理时间字符串转时间
func StrTimeTotime(strTime string) time.Time {

	const Layout = "2006-01-02 15:04:05" //时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	//log.Println(err)
	t, _ := time.ParseInLocation(Layout, strTime /*需要转换的时间类型字符串*/, loc)
	//log.Println(toerr)

	return t
}

//处理时间2017-07-07T11:33:53 转为 2017-07-07 11:33:53
func DataTimeFormat(datetime string) string {

	b := []byte(datetime)
	b[10] = ' '
	log.Println(datetime, "to", string(b))
	return string(b)
}
