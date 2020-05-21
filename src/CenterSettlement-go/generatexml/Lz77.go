package generatexml

/*
#include <stdlib.h>
#include "load_so.h"
*/
import "C"

import (
	log "github.com/sirupsen/logrus"
	"unsafe"
)

// 动态库编译 g++ Lz77.cpp -fPIC -shared -o lz77.so
func Lz77zip(fname string) {

	//把CZ_origin.xml  压缩成 "2.xml.lz77"
	orilz77file := "../sendzipxml/" + fname + ".lz77"

	log.Println("ch name:", fname)
	log.Println("orilz77file :=", orilz77file)

	src2 := C.CString(fname)
	dest2 := C.CString(orilz77file)
	// 压缩
	C.Compressfile(src2, dest2)
	defer C.free(unsafe.Pointer(src2))
	defer C.free(unsafe.Pointer(dest2))
}

func Lz77Unzip() {
	//2.xml.lz77 解压为 1.xml
	originfile := "1.xml"

	orilz77file := "2.xml.lz77"

	// src1 := C.CString(lz77file)
	src1 := C.CString(orilz77file)
	dest1 := C.CString(originfile)
	// 解压
	C.Decompressfile(src1, dest1)

	defer C.free(unsafe.Pointer(src1))
	defer C.free(unsafe.Pointer(dest1))

}
