package geom

/*
#include "geos.h"
*/
import "C"

// BufferParams wraps Ctype BufferParams
type BufferParams struct {
	CBufP *C.GEOSBufferParams
}

// CreateBufferParams creates a bufferparams struct
func CreateBufferParams() (*BufferParams, error) {

	bufP := C.GEOSBufferParams_create()

	return &BufferParams{
		CBufP: bufP,
	}, nil
}

// Destroy clears the memory allcated to BufferParams
func (bp *BufferParams) Destroy() {
	C.GEOSBufferParams_destroy(bp.CBufP)
}

/*
SetEndCapStyle sets the end cap style.

Should be one of below values
GEOSBUF_CAP_ROUND
GEOSBUF_CAP_FLAT
GEOSBUF_CAP_SQUARE
*/
func (bp *BufferParams) SetEndCapStyle(style capstyle) {
	C.GEOSBufferParams_setEndCapStyle(bp.CBufP, C.int(style))
}

/*
SetJoinStyle sets the join cap style
defaults to GEOSBUF_CAP_ROUND
*/
func (bp *BufferParams) SetJoinStyle(style joinstyle) {
	C.GEOSBufferParams_setJoinStyle(bp.CBufP, C.int(style))
}

// SetMitreLimit sets the metre limit which can be used with GEOSBUF_JOIN_MITRE
func (bp *BufferParams) SetMitreLimit(limit float64) {
	C.GEOSBufferParams_setMitreLimit(bp.CBufP, C.double(limit))
}

// SetQuadrantSegments sets the number of segments used to approximate a quarter circle (defaults to 8).
func (bp *BufferParams) SetQuadrantSegments(quadSegs int) {
	C.GEOSBufferParams_setQuadrantSegments(bp.CBufP, C.int(quadSegs))
}

//SetSingleSided sets the params to perform a single-sided buffer on the geometry
func (bp *BufferParams) SetSingleSided() {
	C.GEOSBufferParams_setSingleSided(bp.CBufP, C.int(1))
}
