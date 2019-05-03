package geos

/*
#include "geos.h"
*/
import "C"
import "errors"

// Geom base type for geometric operations
type Geom struct {
	cGeom *C.GEOSGeometry
}

// Buffer Styles
var (
	GeosbufCapRound  capstyle = C.GEOSBUF_CAP_ROUND
	GeosbufCapFlat   capstyle = C.GEOSBUF_CAP_FLAT
	GeosbufCapSquare capstyle = C.GEOSBUF_CAP_SQUARE

	GeosbufJoinRound joinstyle = C.GEOSBUF_JOIN_ROUND
	GeosbufJoinMitre joinstyle = C.GEOSBUF_JOIN_MITRE
	GeosbufJoinBevel joinstyle = C.GEOSBUF_JOIN_BEVEL
)

// Create dedicated type for cap styles and join styles
type capstyle int
type joinstyle int

// GenerateGEOM generates a new go geom that wraps libgeos GEOSGeometry type
func GenerateGEOM(cGeom *C.GEOSGeometry) *Geom {
	return &Geom{
		cGeom: cGeom,
	}
}

// CreatePoint create geometry point
func CreatePoint(x float64, y float64) *Geom {

	coordSeq, _ := initCoordSeq(1, 2)
	// defer coordSeq.Destroy()

	coordSeq.SetX(0, x)
	coordSeq.SetY(0, y)

	return GenerateGEOM(C.GEOSGeom_createPoint_r(ctxHandler, coordSeq.CSeq))
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
func (g *Geom) Buffer(width float64) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSBuffer_r(ctxHandler, g.cGeom, C.double(width), C.int(8))
}

// BufferWithParams creates a buffer around the geometry with specified bufffer options
func (g *Geom) BufferWithParams(width float64, params *BufferParams) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSBufferWithParams_r(ctxHandler, g.cGeom, params.CBufP, C.double(width))
}

// BufferWithStyle buffers a geometry using given style params
func (g *Geom) BufferWithStyle(width float64, quadSegs int, endCapStyle capstyle, joinStyle joinstyle, mitreLimit float64) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)
	g.cGeom = C.GEOSBufferWithStyle_r(ctxHandler, g.cGeom, C.double(width), C.int(quadSegs), C.int(endCapStyle), C.int(joinStyle), C.double(mitreLimit))

}

// Destroy releases the memory allocated to GEOM
func (g *Geom) Destroy() {
	C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)
}

// Simplify simplifies the geometry with a given tolerance
func (g *Geom) Simplify(tolerance float32) {

	//Destroy old geom
	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

	g.cGeom = C.GEOSSimplify_r(ctxHandler, g.cGeom, C.double(tolerance))
}

/*
SimplifyPreserveTopology simplifies the geometry and will avoid creating derived
geometries (polygons in particular) that are invalid.
*/
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

/*
Disjoints Overlaps, Touches, Within all imply geometries are not spatially disjoint.
If any of the aforementioned returns true, then the geometries are not spatially disjoint.
Disjoint implies false for spatial intersection.
*/
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

// Contains returns TRUE if geometry B is completely inside geometry A.
func (g *Geom) Contains(g1 *Geom) (bool, error) {
	contains := C.GEOSContains_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(contains)
}

// Overlaps returns TRUE if geometry B is completely inside geometry A.
func (g *Geom) Overlaps(g1 *Geom) (bool, error) {
	overlaps := C.GEOSOverlaps_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(overlaps)
}

/*
Equals returns true if the DE-9IM intersection matrix for the two Geometrys is T*F**FFF*.

a and b are topologically equal. "Two geometries are topologically equal if their interiors
intersect and no part of the interior or boundary of one geometry intersects the exterior of the other".[9]

equals to Within & Contains
*/
func (g *Geom) Equals(g1 *Geom) (bool, error) {
	eq := C.GEOSEquals_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(eq)
}

// EqualsExact returns true if the two Geometrys are of the same type and their
// vertices corresponding by index are equal up to a specified tolerance.
func (g *Geom) EqualsExact(g1 *Geom, tolerance float64) (bool, error) {
	eqExact := C.GEOSEqualsExact_r(ctxHandler, g.cGeom, g1.cGeom, C.double(tolerance))
	return geosBoolResult(eqExact)
}

/*
Covers returns true if this geometry covers the specified geometry.

The covers predicate has the following equivalent definitions:

    - Every point of the other geometry is a point of this geometry.
    - The DE-9IM Intersection Matrix for the two geometries is T*****FF* or *T****FF* or ***T**FF* or ****T*FF*
    - g.coveredBy(this) (covers is the inverse of coveredBy)

If either geometry is empty, the value of this predicate is false.

This predicate is similar to contains, but is more inclusive (i.e. returns true for more cases).
In particular, unlike contains it does not distinguish between points in the boundary and in the
interior of geometries. For most situations, covers should be used in preference to contains.
As an added benefit, covers is more amenable to optimization, and hence should be more performant.
*/
func (g *Geom) Covers(g1 *Geom) (bool, error) {
	covers := C.GEOSCovers_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(covers)
}

/*
CoveredBy tests whether this geometry is covered by the specified geometry.

The coveredBy predicate has the following equivalent definitions:

    Every point of this geometry is a point of the other geometry.
    The DE-9IM Intersection Matrix for the two geometries matches [T*F**F***] or [*TF**F***] or [**FT*F***] or [**F*TF***]
    g.covers(this) (coveredBy is the converse of covers)

If either geometry is empty, the value of this predicate is false.

This predicate is similar to within, but is more inclusive (i.e. returns true for more cases).
*/
func (g *Geom) CoveredBy(g1 *Geom) (bool, error) {
	coverdBy := C.GEOSCoveredBy_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(coverdBy)
}

// Crosses returns true if this geometry crosses the specified geometry.
func (g *Geom) Crosses(g1 *Geom) (bool, error) {
	crosses := C.GEOSCrosses_r(ctxHandler, g.cGeom, g1.cGeom)
	return geosBoolResult(crosses)
}

// GetNumCoordinates returns the count of this Geometrys vertices.
func (g *Geom) GetNumCoordinates() (int, error) {
	numcoords := C.GEOSGetNumCoordinates_r(ctxHandler, g.cGeom)

	//GEOS return -1 on exception
	if numcoords == C.int(-1) {
		return 0, errors.New(GeosError)
	}

	return int(numcoords), nil

}

// Area returns the area of the geometry
func (g *Geom) Area() (float64, error) {
	var area C.double

	ret := C.GEOSArea_r(ctxHandler, g.cGeom, &area)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(area), nil
}

// Length returns the length of the geometry
func (g *Geom) Length() (float64, error) {
	var len C.double

	ret := C.GEOSLength_r(ctxHandler, g.cGeom, &len)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(len), nil
}

// Distance returns the distance between two geometries
func (g *Geom) Distance(g1 *Geom) (float64, error) {
	var dist C.double

	ret := C.GEOSDistance_r(ctxHandler, g.cGeom, g1.cGeom, &dist)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(dist), nil
}
