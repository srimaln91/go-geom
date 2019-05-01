package geos

/*
#include "geos.h"
*/
import "C"
import "errors"

// CoordinateSeq wraps C coordiante sequence
type CoordinateSeq struct {
	CSeq *C.GEOSCoordSequence
}

func initCoordSeq(size int, dims int) (*CoordinateSeq, error) {

	seq := C.GEOSCoordSeq_create_r(ctxHandler, C.uint(size), C.uint(dims))

	if seq == nil {
		return nil, errors.New("Could not create coordinate sequence")
	}

	return &CoordinateSeq{
		CSeq: seq,
	}, nil
}

// SetX sets the longitude of the coordinate sequence element at given index
func (cs *CoordinateSeq) SetX(idx uint, x float64) {
	C.GEOSCoordSeq_setX_r(ctxHandler, cs.CSeq, C.uint(idx), C.double(x))
}

// SetY sets the latitude of the coordinate sequence element at given index
func (cs *CoordinateSeq) SetY(idx uint, y float64) {
	C.GEOSCoordSeq_setY_r(ctxHandler, cs.CSeq, C.uint(idx), C.double(y))
}

// SetZ sets the altitude of the coordinate sequence element at given index
func (cs *CoordinateSeq) SetZ(idx uint, z float64) {
	C.GEOSCoordSeq_setZ_r(ctxHandler, cs.CSeq, C.uint(idx), C.double(z))
}

// GetX returns the value of x of a given index
func (cs *CoordinateSeq) GetX(idx uint) float64 {
	var val C.double

	res := C.GEOSCoordSeq_getX_r(ctxHandler, cs.CSeq, C.uint(idx), &val)

	if res == 0 {
		return 0.0
	}

	return float64(val)
}

// GetY returns the value of y of a given index
func (cs *CoordinateSeq) GetY(idx uint) float64 {
	var val C.double

	res := C.GEOSCoordSeq_getY_r(ctxHandler, cs.CSeq, C.uint(idx), &val)

	if res == 0 {
		return 0.0
	}

	return float64(val)
}

// GetZ returns the value of z of a given index
func (cs *CoordinateSeq) GetZ(idx uint) float64 {
	var val C.double

	res := C.GEOSCoordSeq_getZ_r(ctxHandler, cs.CSeq, C.uint(idx), &val)

	if res == 0 {
		return 0.0
	}

	return float64(val)
}

// GetSize retuns the size of coordinate sequence
func (cs *CoordinateSeq) GetSize() uint {
	var size C.uint

	res := C.GEOSCoordSeq_getSize_r(ctxHandler, cs.CSeq, &size)

	if res == 0 {
		return 0
	}

	return uint(size)
}

// Destroy clears the coordinate sequence allocations from the memory
func (cs *CoordinateSeq) Destroy() {
	C.GEOSCoordSeq_destroy_r(ctxHandler, cs.CSeq)
}
