package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDescription(t *testing.T) {
	data := "苏A12345|南京南站|0小时35分钟|15元"
	d := Description(data)
	log.Println(d)
}

//		   	1|04|3201|32023411111|1101|20200503 121311|03|3201|32023411111|1101|20200503 111311
//交易详细信息 1|04|3201|3201000002|1101|20200508 134217|03|3201|3201000002|1002|20200506 161547
//      收费车型[ok]|出口类型[04]|出口路网号[3201]|出口站/广场号 |出口车道号|出口时间【ok】 |入口类型【03】｜入口路网号 ｜入口站/广场号 ｜入口车道号 ｜入口时间【ok】
//cx:车型 ckz：出口站、入口站ckcd：出口车道，入口车道cksj：出口时间 rksj：入口时间
func TestDetail(t *testing.T) {
	s := Detail("1", "32023411111", "1101", "2020-05-03 12:13:11", "2020-05-03 11:13:11")
	log.Println(s)
}

func TestJytime(t *testing.T) {
	s := Jytime("2020-05-03 12:13:11")
	log.Println(s)

}
