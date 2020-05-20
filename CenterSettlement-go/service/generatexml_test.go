package service

import (
	"CenterSettlement-go/generatexml"
	"log"
	"testing"
)

func TestGeneratexml(t *testing.T) {

	ch := make(chan string, 0)
	go generatexml.Lz77zip(ch)
	Generatexml(ch)

}

func TestGetMD5Encode(t *testing.T) {
	data := []byte("hello xin")
	str := GetMD5Encode(data)
	log.Println(str, len(str))

}
