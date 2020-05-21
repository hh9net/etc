package generatexml

import "testing"

//不用动态库的
func TestLz77zip(t *testing.T) {
	//f:=make(chan string,0)
	f := "CZ_3201_00000000000000100022xml"

	Lz77zip(f)
}
func TestLz77Unzip(t *testing.T) {
	Lz77Unzip()
}
