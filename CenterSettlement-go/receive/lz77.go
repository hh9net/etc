package generatexml

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
func Lz77UnZip() {
	fmt.Println("data :", C.testa())
	lz77file := "CZ_3201_00000000000000079752.xml.lz77"
	//lz77file := "CZ_3201_00000000000000000001.xml.lz77"

	//xml源文件
	originfile := "00000000000000079752.xml"

	// 解压文件
	src1 := C.CString(lz77file)
	dest1 := C.CString(originfile)
	// 解压
	C.decompresslz77(src1, dest1)
	defer C.free(unsafe.Pointer(src1))
	defer C.free(unsafe.Pointer(dest1))
}
