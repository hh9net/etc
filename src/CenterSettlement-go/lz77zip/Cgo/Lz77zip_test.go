package Cgo

import "testing"

//不用动态库的 静态库
func TestLz77zip(t *testing.T) {
	f := "CZ_3101_00000000000000100136.xml"
	Lz77zip(f)
}
func TestLz77Unzip(t *testing.T) {
	f := "CZ_3101_00000000000000100136.xml.lz77"
	Lz77Unzipxml(f)
}
