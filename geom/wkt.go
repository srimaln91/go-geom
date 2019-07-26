package geom

/*
#include "geos.h"
*/
import "C"

import (
	"unsafe"
)

type wktReader struct {
	c *C.GEOSWKTReader
}

type wktWriter struct {
	c *C.GEOSWKTWriter
}

func (r *wktReader) read(wkt string) *GEOSGeom {
	cs := C.CString(wkt)
	defer C.free(unsafe.Pointer(cs))

	cGeom := C.GEOSWKTReader_read(r.c, cs)
	return GenerateGeosGeom(cGeom)
}

func (w *wktWriter) write(g *GEOSGeom) string {

	//Set output dimension (2 or 3)
	outputDimention := C.GEOSGeom_getCoordinateDimension(g.cGeom)
	C.GEOSWKTWriter_setOutputDimension(w.c, outputDimention)

	wkt := C.GEOSWKTWriter_write(w.c, g.cGeom)

	//Clear buffers
	defer C.GEOSFree(unsafe.Pointer(wkt))

	return C.GoString(wkt)
}

func createWktReader() *wktReader {

	c := C.GEOSWKTReader_create()
	if c == nil {
		return nil
	}

	r := &wktReader{c: c}

	// Instruct Go garbage collector to clean C memory blocks
	// runtime.SetFinalizer(r, func(r *wktReader) {
	// 	C.GEOSWKTReader_destroy_r(ctxHandler, r.c)
	// })

	return r
}

// Destroy releases the memory allocated to WKT reader
func (r *wktReader) destroy() {
	C.GEOSWKTReader_destroy(r.c)
}

func createWktWriter() *wktWriter {
	c := C.GEOSWKTWriter_create()
	if c == nil {
		return nil
	}

	w := &wktWriter{c: c}

	// Instruct Go garbage collector to clean C memory blocks
	// runtime.SetFinalizer(w, func(w *wktWriter) {
	// 	C.GEOSWKTWriter_destroy_r(ctxHandler, w.c)
	// })

	return w
}

// Destroy releases the memory allocated to WKT writer
func (w *wktWriter) destroy() {
	C.GEOSWKTWriter_destroy(w.c)
}
