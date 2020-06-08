package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDateformat(t *testing.T) {

	//2006-01-02T15:04:05
	DateFormat()
	//2006-01-02
	TodayFormat()
}
func TestDataTimeFormatHandle(t *testing.T) {
	d := "2017-07-07T11:33:53"
	s := DataTimeFormatHandle(d)

	log.Println(s)

}

func TestStrTimeTotime(t *testing.T) {
	T := StrTimeTotime("2017-07-07 11:33:53")
	log.Println(T)
}
