package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestTimeDifference(t *testing.T) {
	//进入时间   出去时间
	TimeDifference("2020-03-14 13:13:12", "2020-05-20 15:10:05")
	h := getHourDiffer("2020-03-14 13:13:12", "2020-05-20 15:10:05")
	m := getMDiffer("2020-03-14 13:13:12", "2020-05-20 15:10:05")
	s := getSDiffer("2020-03-14 13:13:12", "2020-05-20 15:10:05")
	log.Println(h/24, m/24/60, s/24/60/60)
	log.Println(h, m, s)
}
