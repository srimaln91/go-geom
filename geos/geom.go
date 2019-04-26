package geos

/*
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

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSBuffer_r(ctxHandler, g.cGeom, C.double(width), C.int(8))
}

// Destroy releases the memory allocated to GEOM
func (g *Geom) Destroy() {
	C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)
}

// Simplify simplifies the geometry with given tolerance
func (g *Geom) Simplify(tolerance float32) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSSimplify_r(ctxHandler, g.cGeom, C.double(tolerance))
}

// SimplifyPreserveTopology simplifies the geometry and will avoid creating derived geometries (polygons in particular) that are invalid.
func (g *Geom) SimplifyPreserveTopology(tolerance float32) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSTopologyPreserveSimplify_r(ctxHandler, g.cGeom, C.double(tolerance))
}

// Reverse reverses the geom
func (g *Geom) Reverse() {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSReverse_r(ctxHandler, g.cGeom)
}

// Union returns the union of two geometries
func (g *Geom) Union(g1 *Geom) *Geom {

	union := C.GEOSUnion_r(ctxHandler, g.cGeom, g1.cGeom)
	return GenerateGEOM(union)
}

// Intersection retuerns the intersection of 2 geometries
func (g *Geom) Intersection(g1 *Geom) *Geom {
	intersection := C.GEOSIntersection_r(ctxHandler, g.cGeom, g1.cGeom)
	return GenerateGEOM(intersection)
}
