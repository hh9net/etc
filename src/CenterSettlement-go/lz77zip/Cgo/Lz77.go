package Cgo

/*
#include <stdlib.h>
#include "load_so.h"
*/
import "C"

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"unsafe"
)

// 静态编译 g++ Lz77.cpp -fPIC -shared -o lz77.so
func Lz77zip(fname string) {
	//把CZ_origin.xml  压缩成 "2.xml.lz77"
	orilz77file := "../../sendzipxml/" + "Cgo_" + fname + ".lz77"
	fn := "../../generatexml/" + fname
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
	originfile := "../../receivexml/" + "CgoUnZip_" + fstr[0] //.xml
	orilz77file := "../../sendzipxml/" + fname                //.xml.lz77
	//originfile := "../sendxmlsucceed/" + fstr[0]//.xml
	//orilz77file := "../sendxmlsucceed/" + fname//.xml.lz77
	log.Println("7zfile:", orilz77file)
	log.Println("xmlfile:", originfile)

	src1 := C.CString(orilz77file)
	dest1 := C.CString(originfile)
	// 解压
	C.Decompressfile(src1, dest1)

	defer C.free(unsafe.Pointer(src1))
	defer C.free(unsafe.Pointer(dest1))

}

func Zip(fname string) {
	//把CZ_origin.xml  压缩成 "2.xml.lz77"  CenterSettlement-go
	orilz77file := "../centerserver/" + fname + ".lz77"
	fn := "../centerserver/" + fname
	src2 := C.CString(fn)
	dest2 := C.CString(orilz77file)
	// 压缩
	C.Compressfile(src2, dest2)
	defer C.free(unsafe.Pointer(src2))
	defer C.free(unsafe.Pointer(dest2))
}
