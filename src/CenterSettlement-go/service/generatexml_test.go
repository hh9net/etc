package service

import (
	"CenterSettlement-go/generatexml"
	"log"
	"testing"
)

func TestGeneratexml(t *testing.T) {
	//Generatexml(22)
	s := Generatexml(22)
	generatexml.Lz77zip(s)
}

func TestGetMD5Encode(t *testing.T) {
	data := []byte("hello xin")
	str := GetMD5Encode(data)
	log.Println(str, len(str))

}
