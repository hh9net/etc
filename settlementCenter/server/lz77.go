package server

/*
#include <stdio.h>
#include <stdlib.h>
#include "load_so.h"
#cgo LDFLAGS: -ldl

// 文件写在此处,才会触发c部分代码的重新编译，直接写在.h文件离的不会触发重新编译
int testa()
{
	return testcode_c();
}

// 压缩
int compresslz77(char *srcfile, char *destfile)
{
	return compress(srcfile, destfile);
}

// 解压缩
int decompresslz77(char *srcfile , char *destfile)
{
	return decompress(srcfile, destfile);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// 动态库编译 g++ Lz77.cpp -fPIC -shared -o lz77.so
func Lz77Compress(fname string) {
	fmt.Println("data :", C.testa())
	//xml文件
	xmlfile := fname
	//压缩后的文件名
	orilz77file := fname + ".lz77"
	//压缩文件
	src2 := C.CString(xmlfile)
	dest2 := C.CString(orilz77file)
	// 压缩
	C.compresslz77(src2, dest2)

	defer C.free(unsafe.Pointer(src2))
	defer C.free(unsafe.Pointer(dest2))
}
