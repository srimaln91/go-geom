package geos

/*
#cgo CFLAGS: -g -Wall -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -L/usr/local/lib -L/usr/lib -lgeos_c
#include "geos.h"
*/
import "C"

// Geom base type for geometric operations
type Geom struct {
	cGeom *C.GEOSGeometry
}

// GenerateGEOM generates a new go geom
func GenerateGEOM(cGeom *C.GEOSGeometry) *Geom {
	return &Geom{
		cGeom: cGeom,
	}
}

// FromWKT creates a Geom from WKT string
func FromWKT(wkt string) *Geom {

	wktReader := createWktReader()
	defer wktReader.destroy()

	geom := wktReader.read(wkt)

	return geom
}

// ToWKT returns a WKT string
func (g *Geom) ToWKT() string {

	wktWriter := createWktWriter()
	defer wktWriter.destroy()

	wkt := wktWriter.write(g)
	return wkt
}

// SetSRID sets SRID of the geometry
func (g *Geom) SetSRID(srid int) {
	C.GEOSSetSRID_r(ctxHandler, g.cGeom, C.int(srid))
}

// GetSRID returns the SRID of the geometry
func (g *Geom) GetSRID() int {
	srid := C.GEOSGetSRID_r(ctxHandler, g.cGeom)
	return int(srid)
}

// Buffer creates a buffer around the geometry
func (g *Geom) Buffer(width float32) {
	g.cGeom = C.GEOSBuffer_r(ctxHandler, g.cGeom, C.double(width), C.int(8))
}

// Destroy releases the memory allocated to GEOM
func (g *Geom) Destroy() {
	C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)
}
