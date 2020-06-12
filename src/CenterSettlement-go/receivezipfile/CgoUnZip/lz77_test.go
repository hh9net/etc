package CgoUnZip

import "testing"

func TestLz77UnZip(t *testing.T) {
	//测试加压缩
	Lz77UnZip("00000000000000100025.xml.lz77")

}
