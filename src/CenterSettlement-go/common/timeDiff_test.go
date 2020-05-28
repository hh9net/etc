package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestTimeDifference(t *testing.T) {
	//进入时间   出去时间
	var date string
	date = TimeDifference("2020-05-19 01:00:00", "2020-05-20 20:00:00")
	log.Println("进入时间   出去时间 时间差", date)

}
