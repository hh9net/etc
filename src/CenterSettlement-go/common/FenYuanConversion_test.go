package common

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestFen2Yuan(t *testing.T) {
	Fen2Yuan(1850)
}

func TestYuan2Fen(t *testing.T) {
	s := Yuan2Fen(18.5675)
	log.Println(s)
}
