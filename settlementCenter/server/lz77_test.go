package server

import "testing"

func TestLz77Compress(t *testing.T) {
	//测试压缩

	fname := "CZ_3201_00000000000000999999.xml"
	Lz77Compress(fname)
}
