package lz77zip

import "testing"

//不用动态库的
func TestLz77zip(t *testing.T) {
	//f:=make(chan string,0)
	f := "JZ_3301_00000000000000100094.xml"

	Lz77zip(f)
}
func TestLz77Unzip(t *testing.T) {
	f := "JZ_3301_00000000000000100094.xml.lz77"
	Lz77Unzipxml(f)
}
