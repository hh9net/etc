package common

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func TimeDifference(intime time.Time, outtime time.Time) string {

	var day string
	in := intime.Format("2006-01-02 15:04:05")
	out := outtime.Format("2006-01-02 15:04:05")
	inT, _ := time.Parse("2006-01-02 15:04:05", in)
	outT, _ := time.Parse("2006-01-02 15:04:05", out)
	t := outT.Sub(inT).String()
	//t := outTime.Sub(inTime).String()
	log.Println("时间差", t)
	if strings.Contains(t, "h") {
		//h
		hourstr := strings.Split(t, "h")
		hour, _ := strconv.Atoi(hourstr[0])
		//m eg:59m59s
		m := strings.Split(hourstr[1], "m")
		//s eg:59s
		s := strings.Split(m[1], "s")
		//h
		if hour >= 24 {
			//小时转天
			d := hour / 24
			h := hour % 24
			dstr := strconv.Itoa(d)
			hstr := strconv.Itoa(h)
			day = dstr + "天" + hstr + "时" + m[0] + "分" + s[0] + "秒"
			return day
		}

		if 24 > hour && hour > 0 {
			day = hourstr[0] + "时" + m[0] + "分" + s[0] + "秒"
			return day
		}
	} else if strings.Contains(t, "m") {
		//m eg:59m59s
		m := strings.Split(t, "m")
		//s eg:59s
		s := strings.Split(m[1], "s")

		day = m[0] + "分" + s[0] + "秒"
		return day
	} else if strings.Contains(t, "s") {
		//s eg:59s
		s := strings.Split(t, "s")
		day = s[0] + "秒"
		return day
	}
	return ""
}
