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
func Detail(cx string, ckz string, ckcd string, cksj string, rksj string) string {
	//2020-05-03 12:13:11
	//log.Println("cksj", cksj)
	//log.Println("rksj", rksj)

	//log.Printf("(cx %s, ckz %s, ckcd %s, cksj %s, rksj %s", cx, ckz, ckcd, cksj, rksj)
	csj := timeDetail(cksj)
	rsj := timeDetail(rksj)
	detail := cx + "|04|" + "3201|" + ckz + "|" + ckcd + "|" + csj + "|03|" + "3201|" + ckz + "|" + ckcd + "|" + rsj
	log.Printf(detail)

	return detail
}

//20200508 134217
func timeDetail(sj string) string {
	s := strings.Split(sj, "-")
	s1 := s[0] + s[1] + s[2]
	s2 := strings.Split(s1, ":")
	s3 := s2[0] + s2[1] + s2[2]
	return s3
}

//20200508134217
func Jytime(sj string) string {
	s := strings.Split(sj, "-")
	s1 := s[0] + s[1] + s[2]
	s2 := strings.Split(s1, ":")
	s3 := s2[0] + s2[1] + s2[2]
	s4 := strings.Split(s3, " ")
	s5 := s4[0] + s4[1]
	return s5
}

//出口类型 写固定值【04】    出口网络号 写固定值【3201】 出口站 写F_VC_TINGCCBH， 出口车道号，写F_VC_CHEDID, 入口类型，写固定【03】，
//入口站 写F_VC_TINGCCBH ,入口车道号，写F_VC_CHEDID

func Name(data string) string {
	s := strings.Split(data, "|")
	d := s[1]
	log.Println(data)
	return d
}
func GetLiush(FVcChedjyxh string) string {
	Liush := []byte(FVcChedjyxh)
	//log.Println(len(FVcChedjyxh))
	if len(FVcChedjyxh) == 0 {
		return "00"
	}

	if len(FVcChedjyxh) == 1 {
		return "0" + string(Liush[:])

	}
	return string(Liush[len(FVcChedjyxh)-2:])
}

//停车场消费交易编号(停车场编号+交易发生的时间+流水号 )
func GetId(tingccbh string, jysj string, liush string) string {
	sj := Jytime(jysj)
	hb := tingccbh + sj + liush
	log.Println("停车场消费交易编号(停车场编号+交易发生的时间+流水号 )", hb)
	return hb
}

//TAC校验码	TAC	hexbinary	8【TAC码 FVcTacm（32位）】
//交易金额	TxnAmt	hexbinary	8【ok】
//交易类型	TxnType	hexbinary	2【ok固定 09】
//终端机编号	PosId	hexbinary	12 【加密卡好】
//交易的终端流水号	PosSeq	hexbinary	8【加密序列号】
//交易时间	TxTime(TxnDate+TxnTime)	hexbinary	14【ok】
//TAC保留串	Reserve	hexbinary	4【0000】
//交易后余额	Balance	hexbinary	8【ok】
//交易前余额	PreBalance	hexbinary	8【ok】
//卡计数器	Iccounter	hexbinary	4【卡交易序号】

// <CustomizedData>7888180f 000000C8 09 01320002d9f7 000017d5 20200513143434 0000 7FC4BD32 7FC4BDFA 0da2</CustomizedData>
func CustomizedData(tac string, jyje string, jylx string, zdjbh string, zdlsh string, jysj string, jyhye string, jyqye string, kjsq string) string {
	//2020-05-03 12:13:11
	sj := Jytime(jysj)
	s := tac + jyje + jylx + zdjbh + zdlsh + sj + "0000" + jyhye + jyqye + kjsq
	return s
}

/*
时间常量
*/
const (
	//定义每分钟的秒数
	SecondsPerMinute = 60
	//定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	//定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

/*
时间转换函数
*/
func ResolveTime(seconds int) (day int, hour int, minute int) {
	//每分钟秒数
	minute = seconds / SecondsPerMinute
	//每小时秒数
	hour = seconds / SecondsPerHour
	//每天秒数
	day = seconds / SecondsPerDay
	return
}
