package server

import (
	"log"
	"testing"
)

func TestGeneratexml(t *testing.T) {
	Generatexml()
}

func TestGetMD5Encode(t *testing.T) {
	data := []byte("hello xin")
	str := GetMD5Encode(data)
	log.Println(str, len(str))

}
