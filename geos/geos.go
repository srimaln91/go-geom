package geos

/*
#cgo CFLAGS: -g -Wall -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -L/usr/local/lib -L/usr/lib -lgeos_c
#include "geos.h"
*/
import "C"

var (
	ctxHandler C.GEOSContextHandle_t
)

func init() {
	ctxHandler = C.init_geos()
}

// GoFinishGEOS removed libgeos allocations from the memory
func GoFinishGEOS() {
	C.finishGEOS_r(ctxHandler)
}
