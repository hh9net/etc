package common

import (
	"encoding/binary"
	"log"
	"testing"
)

func TestToHex(t *testing.T) {
	a := int64(49)
	s := ToHex(a)
	log.Println(s)
	//00000000000000000031 20位

}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func TestHexToTen(t *testing.T) {
	//HexToTen("31")
	HexToTen("0fa8")
}

func TestToHexFormat(t *testing.T) {
	s := ToHexFormat(int64(49), 8)
	log.Println(s)
}
