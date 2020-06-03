package lz77zip

/*
#include <stdlib.h>
#include "load_so.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

// 静态编译 g++ Lz77.cpp -fPIC -shared -o lz77.so
func Lz77zip(fname string) {
	//把CZ_origin.xml  压缩成 "2.xml.lz77"
	orilz77file := "../sendzipxml/" + fname + ".lz77"
	fn := "../generatexml/" + fname
	src2 := C.CString(fn)
	dest2 := C.CString(orilz77file)
	// 压缩
	C.Compressfile(src2, dest2)
	defer C.free(unsafe.Pointer(src2))
	defer C.free(unsafe.Pointer(dest2))
}

func Lz77Unzipxml(fname string) {
	//2.xml.lz77 解压为 1.xml "00000000000000100025.xml.lz77"
	fstr := strings.Split(fname, ".lz77")

	originfile := "../generatexml/" + fstr[0]
	orilz77file := "../generatexml/" + fname

	src1 := C.CString(orilz77file)
	dest1 := C.CString(originfile)
	// 解压
	C.Decompressfile(src1, dest1)

	defer C.free(unsafe.Pointer(src1))
	defer C.free(unsafe.Pointer(dest1))

}
