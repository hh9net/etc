package lib7z

/*
#cgo CFLAGS: -I./
#cgo LDFLAGS: -L./ -llz77 -ldl
#include <stdio.h>
#include <stdlib.h>
#include "Lz77include.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Lz77so struct {
}

func (p *Lz77so) Comresslz77(source string, dest string) {
	csourcestr := C.CString(source)
	cdeststr := C.CString(dest)

	fmt.Println("libinlz77so \n")
	C.Compressfilefunc(csourcestr, cdeststr)
	C.free(unsafe.Pointer(csourcestr))
	C.free(unsafe.Pointer(cdeststr))
	fmt.Println("out liblz77so \n")

}

func (p *Lz77so) Depresslz77(source string, dest string) {

	csourcestr := C.CString(source)
	cdeststr := C.CString(dest)

	C.Decompressfilefunc(csourcestr, cdeststr)
	C.free(unsafe.Pointer(csourcestr))
	C.free(unsafe.Pointer(cdeststr))
}

func (p *Lz77so) Test() {
	C.Testcode()
}
