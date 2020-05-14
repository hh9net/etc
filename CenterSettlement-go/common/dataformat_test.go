package common

import "testing"

func TestDateformat(t *testing.T) {
	//2006-01-02 15:04:05
	DateTimeFormat()
	//2006-01-02T15:04:05
	DateFormat()
	//2006-01-02
	TodayFormat()
}
