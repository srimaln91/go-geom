package geos

/*
#cgo CFLAGS: -g -Wall -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -L/usr/local/lib -L/usr/lib -lgeos_c
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <geos_c.h>
*/
import "C"
import "sync"

var (
	mu sync.Mutex
)

// Geom base type for the geometric operations
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
	// defer wktReader.Destroy()
	geom := wktReader.read(wkt)

	return geom
}

// ToWKT returns a WKT string
func (g *Geom) ToWKT() string {
	wktWriter := createWktWriter()
	// defer wktWriter.Destroy()
	return wktWriter.write(g)
}

//Buffer creates a buffer around the geometry
func (g *Geom) Buffer(width float32) {
	mu.Lock()
	defer mu.Unlock()
	g.cGeom = C.GEOSBuffer_r(ctxHandler, g.cGeom, C.double(width), C.int(8))
}

// Destroy releases the memory allocated to GEOM
func (g *Geom) Destroy() {
	C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)
}
