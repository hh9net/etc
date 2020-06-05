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
	orilz77file := "CenterSettlement-go/sendzipxml/" + fname + ".lz77"
	fn := "CenterSettlement-go/generatexml/" + fname
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
	originfile := "CenterSettlement-go/center_server/" + fstr[0] //.xml
	orilz77file := "CenterSettlement-go/center_server/" + fname  //.xml.lz77
	//originfile := "../sendxmlsucceed/" + fstr[0]//.xml
	//orilz77file := "../sendxmlsucceed/" + fname//.xml.lz77

	src1 := C.CString(orilz77file)
	dest1 := C.CString(originfile)
	// 解压
	C.Decompressfile(src1, dest1)

	defer C.free(unsafe.Pointer(src1))
	defer C.free(unsafe.Pointer(dest1))

}

func Zip(fname string) {
	//把CZ_origin.xml  压缩成 "2.xml.lz77"  CenterSettlement-go
	orilz77file := "../center_server/" + fname + ".lz77"
	fn := "../center_server/" + fname
	src2 := C.CString(fn)
	dest2 := C.CString(orilz77file)
	// 压缩
	C.Compressfile(src2, dest2)
	defer C.free(unsafe.Pointer(src2))
	defer C.free(unsafe.Pointer(dest2))
}
