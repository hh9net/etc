package lz77zip

import (
	"CenterSettlement-go/lz77zip/Cgo"
	"testing"
)

//go zip
func TestZipLz77(t *testing.T) {
	//f:="CZ_3101_00000000000000100136.xml"
	//ZipLz77(f)
}

//go unzip
func TestUnZipLz77(t *testing.T) {
	//f:="CZ_3101_00000000000000100136.xml.lz77"
	//UnZipLz77(f)

}

//Cgo  交叉测试
func TestUnZipLz772(t *testing.T) {
	f := "CZ_3101_00000000000000100136.xml.lz77"
	Cgo.Lz77Unzipxml(f)
}
