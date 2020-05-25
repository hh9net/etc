package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDateformat(t *testing.T) {
	//2006-01-02 15:04:05
	DateTimeFormat()
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
