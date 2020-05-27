package common

import (
	"log"
	"strings"
)

func Description(data string) string {
	s := strings.Split(data, "|")
	d := s[1] + "|" + s[2]
	log.Println(data)
	return d
}

//交易详细信息 1|04|3201|3201000002|1101|20200508 134217|03|3201|3201000002|1002|20200506 161547
//      收费车型[ok]|出口类型[04]|出口路网号[3201]|出口站/广场号 |出口车道号|出口时间【ok】 |入口类型【03】｜入口路网号 ｜入口站/广场号 ｜入口车道号 ｜入口时间【ok】
//cx:车型 ckz：出口站、入口站ckcd：出口车道，入口车道cksj：出口时间 rksj：入口时间
func Detail(cx string,ckz string,ckcd string,cksj string,rksj string) string {
	//2020-05-03 12:13:11
	csj:=timeDetail(cksj)
	rsj:=timeDetail(rksj)
	detail := cx+"|04|"+"3201|"+ckz+"|"+ckcd+"|"+csj+"|03|"+"3201|"+ckz+"|"+ckcd+"|"+rsj
	return detail
}
func timeDetail(sj string)string{
	s := strings.Split(sj, "-")
	s1:=s[0]+s[1]+s[2]
	s2 := strings.Split(s1, ":")
	s3:=s2[0]+s2[1]+s2[2]
	return s3
}

//出口类型 写固定值【04】    出口网络号 写固定值【3201】 出口站 写F_VC_TINGCCBH， 出口车道号，写F_VC_CHEDID, 入口类型，写固定【03】，
//入口站 写F_VC_TINGCCBH ,入口车道号，写F_VC_CHEDID

func Name(data string) string {
	s := strings.Split(data, "|")
	d := s[1] 
	log.Println(data)
	return d
}

//停车场消费交易编号(停车场编号+交易发生的时间+流水号 )

func GetId()  {
	
}