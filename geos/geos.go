package geos

/*
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

// Version returns the GEOS version
func Version() string {
	version := C.GEOSversion()
	return C.GoString(version)
}
