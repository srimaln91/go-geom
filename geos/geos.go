package geos

/*
#include "geos.h"
*/
import "C"
import "errors"

// GeosError generalizes GEOS error
// Should look into std_err in order to get descriptive output
const GeosError = "GEOS ERROR"

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

// GEOS Binary predicates - return 2 on exception, 1 on true, 0 on false
func geosBoolResult(char _Ctype_char) (bool, error) {

	switch C.int(char) {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, errors.New(GeosError)
	}
}
