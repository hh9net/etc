package common

import (
	"testing"
	"time"
)

func TestTimeDifference(t *testing.T) {
	t1 := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second * 10)
	t2 := time.Now().Format("2006-01-02 15:04:05")
	TimeDifference(t2, t1)
}
