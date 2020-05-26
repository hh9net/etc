package common

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func TimeDifference(intime time.Time, outTime time.Time) {
	t := outTime.Sub(intime)
	log.Println("时间差", t)

}
