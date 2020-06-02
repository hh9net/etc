package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestTimeDifference(t *testing.T) {
	//进入时间   出去时间
	var date string
	var t1 time.Time
	var t2 time.Time
	t1 = time.Now()
	log.Println(t1)
	time.Sleep(time.Second * 10)
	t2 = time.Now()
	log.Println(t2)
	date = TimeDifference(t1, t2)
	log.Println("进入时间   出去时间 时间差", date)

}
