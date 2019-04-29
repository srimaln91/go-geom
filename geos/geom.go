package geos

/*
#include "geos.h"
*/
import "C"

// Geom base type for geometric operations
type Geom struct {
	cGeom *C.GEOSGeometry
}

// GenerateGEOM generates a new go geom that wraps libgeos GEOSGeometry type
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

// Intersection returns the intersection of 2 geometries
func (g *Geom) Intersection(g1 *Geom) *Geom {
	intersection := C.GEOSIntersection_r(ctxHandler, g.cGeom, g1.cGeom)
	return GenerateGEOM(intersection)
}

// Intersects checks whether the geom intersects with an another geom
func (g *Geom) Intersects(g1 *Geom) (bool, error) {
	intersects := C.GEOSIntersects_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(intersects)
}

// Disjoints Overlaps, Touches, Within all imply geometries are not spatially disjoint.
// If any of the aforementioned returns true, then the geometries are not spatially disjoint. Disjoint implies false for spatial intersection.
func (g *Geom) Disjoints(g1 *Geom) (bool, error) {
	disjoint := C.GEOSDisjoint_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(disjoint)
}

// Touches returns TRUE if the only points in common between g1 and g2 lie in the union of the boundaries of g1 and g2.
func (g *Geom) Touches(g1 *Geom) (bool, error) {
	touches := C.GEOSTouches_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(touches)
}

// Within returns TRUE if geometry A is completely inside geometry B
func (g *Geom) Within(g1 *Geom) (bool, error) {
	within := C.GEOSWithin_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(within)
}
