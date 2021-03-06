package geom

/*
#include "geos.h"
#include "geom.h"
*/
import "C"
import "errors"

/*
GeosError generalizes GEOS error
Should look into std_err in order to get descriptive output
*/
const GeosError = "GEOS ERROR"

var (
	ctxHandler C.GEOSContextHandle_t
)

// init initializes the package
func init() {
	GoInitGEOS()
}

// GoInitGEOS initializes libgeos
func GoInitGEOS() {
	C.init_geos()
}

// GoFinishGEOS remove libgeos allocations from the memory
func GoFinishGEOS() {
	C.finish_geos()
}

// Version returns the GEOS version
func Version() string {
	version := C.GEOSversion()
	return C.GoString(version)
}

// geosBoolResult evaluates C.char into boolean
func geosBoolResult(result C.char) (bool, error) {

	// GEOS Binary predicates - return 2 on exception, 1 on true, 0 on false
	switch result {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, errors.New(GeosError)
	}
}

// GEOSGeom base type for geometric operations
type GEOSGeom struct {
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

// GenerateGeosGeom generates a new go geom that wraps libgeos GEOSGeometry type
func GenerateGeosGeom(cGeom *C.GEOSGeometry) *GEOSGeom {
	return &GEOSGeom{
		cGeom: cGeom,
	}
}

// CreatePoint create geometry point
func CreatePoint(x float64, y float64) *GEOSGeom {

	coordSeq, _ := initCoordSeq(1, 2)
	// defer coordSeq.Destroy()

	coordSeq.SetX(0, x)
	coordSeq.SetY(0, y)

	return GenerateGeosGeom(C.GEOSGeom_createPoint(coordSeq.CSeq))
}

// FromWKT creates a Geom from WKT string
func FromWKT(wkt string) *GEOSGeom {

	wktReader := createWktReader()
	defer wktReader.destroy()

	geom := wktReader.read(wkt)

	return geom
}

// ToWKT returns a WKT string
func (g *GEOSGeom) ToWKT() string {

	wktWriter := createWktWriter()
	defer wktWriter.destroy()

	wkt := wktWriter.write(g)
	return wkt
}

// SimplifiedBufferFromWkt simplifies and buffers inut wkt and
func SimplifiedBufferFromWkt(wkt string, width float64, tolerance float64) string {
	return C.GoString(C.simplified_buffer_from_wkt(C.CString(wkt), C.double(width), C.double(tolerance)))
}

// SetSRID sets SRID of the geometry
func (g *GEOSGeom) SetSRID(srid int) {
	C.GEOSSetSRID(g.cGeom, C.int(srid))
}

// GetSRID returns the SRID of the geometry
func (g *GEOSGeom) GetSRID() int {

	srid := C.GEOSGetSRID(g.cGeom)

	return int(srid)
}

// Buffer creates a buffer around the geometry
func (g *GEOSGeom) Buffer(width float64) {

	//Destroy old geom
	defer C.GEOSGeom_destroy(g.cGeom)

	g.cGeom = C.GEOSBuffer(g.cGeom, C.double(width), C.int(8))
}

// BufferWithParams creates a buffer around the geometry with specified bufffer options
func (g *GEOSGeom) BufferWithParams(width float64, params *BufferParams) {

	//Destroy old geom
	defer C.GEOSGeom_destroy(g.cGeom)

	g.cGeom = C.GEOSBufferWithParams(g.cGeom, params.CBufP, C.double(width))
}

// BufferWithStyle buffers a geometry using given style params
func (g *GEOSGeom) BufferWithStyle(width float64, quadSegs int, endCapStyle capstyle, joinStyle joinstyle, mitreLimit float64) {

	//Destroy old geom
	defer C.GEOSGeom_destroy(g.cGeom)
	g.cGeom = C.GEOSBufferWithStyle(g.cGeom, C.double(width), C.int(quadSegs), C.int(endCapStyle), C.int(joinStyle), C.double(mitreLimit))

}

// SimplifiedBuffer simplifies a geometry with a given tolerance and creates a buffer around that
func (g *GEOSGeom) SimplifiedBuffer(tolerance float64, width float64) {
	defer C.GEOSGeom_destroy(g.cGeom)
	g.cGeom = C.simplified_buffer(g.cGeom, C.double(width), C.double(tolerance))
}

// Destroy releases the memory allocated to GEOM
func (g *GEOSGeom) Destroy() {
	C.GEOSGeom_destroy(g.cGeom)
}

// Simplify simplifies the geometry with a given tolerance
func (g *GEOSGeom) Simplify(tolerance float32) {

	//Destroy old geom
	defer C.GEOSGeom_destroy(g.cGeom)

	g.cGeom = C.GEOSSimplify(g.cGeom, C.double(tolerance))
}

/*
SimplifyPreserveTopology simplifies the geometry and will avoid creating derived
geometries (polygons in particular) that are invalid.
*/
func (g *GEOSGeom) SimplifyPreserveTopology(tolerance float32) {

	//Destroy old geom
	defer C.GEOSGeom_destroy(g.cGeom)

	g.cGeom = C.GEOSTopologyPreserveSimplify(g.cGeom, C.double(tolerance))
}

// Reverse reverses the geom
// func (g *GEOSGeom) Reverse() {

// 	//Destroy old geom
// 	defer C.GEOSGeom_destroy_r(ctxHandler, g.cGeom)

// 	g.cGeom = C.GEOSReverse_r(ctxHandler, g.cGeom)
// }

// Union returns the union of two geometries
func (g *GEOSGeom) Union(g1 *GEOSGeom) *GEOSGeom {

	union := C.GEOSUnion(g.cGeom, g1.cGeom)
	return GenerateGeosGeom(union)
}

// Intersection returns the intersection of 2 geometries
func (g *GEOSGeom) Intersection(g1 *GEOSGeom) *GEOSGeom {
	intersection := C.GEOSIntersection(g.cGeom, g1.cGeom)
	return GenerateGeosGeom(intersection)
}

// Intersects checks whether the geom intersects with an another geom
func (g *GEOSGeom) Intersects(g1 *GEOSGeom) (bool, error) {
	intersects := C.GEOSIntersects(g.cGeom, g1.cGeom)
	return geosBoolResult(intersects)
}

/*
Disjoints Overlaps, Touches, Within all imply geometries are not spatially disjoint.
If any of the aforementioned returns true, then the geometries are not spatially disjoint.
Disjoint implies false for spatial intersection.
*/
func (g *GEOSGeom) Disjoints(g1 *GEOSGeom) (bool, error) {
	disjoint := C.GEOSDisjoint(g.cGeom, g1.cGeom)
	return geosBoolResult(disjoint)
}

// Touches returns TRUE if the only points in common between g1 and g2 lie in the union of the boundaries of g1 and g2.
func (g *GEOSGeom) Touches(g1 *GEOSGeom) (bool, error) {
	touches := C.GEOSTouches(g.cGeom, g1.cGeom)
	return geosBoolResult(touches)
}

// Within returns TRUE if geometry A is completely inside geometry B
func (g *GEOSGeom) Within(g1 *GEOSGeom) (bool, error) {
	within := C.GEOSWithin(g.cGeom, g1.cGeom)
	return geosBoolResult(within)
}

// Contains returns TRUE if geometry B is completely inside geometry A.
func (g *GEOSGeom) Contains(g1 *GEOSGeom) (bool, error) {
	contains := C.GEOSContains(g.cGeom, g1.cGeom)
	return geosBoolResult(contains)
}

// Overlaps returns TRUE if geometry B is completely inside geometry A.
func (g *GEOSGeom) Overlaps(g1 *GEOSGeom) (bool, error) {
	overlaps := C.GEOSOverlaps(g.cGeom, g1.cGeom)
	return geosBoolResult(overlaps)
}

/*
Equals returns true if the DE-9IM intersection matrix for the two Geometrys is T*F**FFF*.

a and b are topologically equal. "Two geometries are topologically equal if their interiors
intersect and no part of the interior or boundary of one geometry intersects the exterior of the other".[9]

equals to Within & Contains
*/
func (g *GEOSGeom) Equals(g1 *GEOSGeom) (bool, error) {
	eq := C.GEOSEquals(g.cGeom, g1.cGeom)
	return geosBoolResult(eq)
}

// EqualsExact returns true if the two Geometrys are of the same type and their
// vertices corresponding by index are equal up to a specified tolerance.
func (g *GEOSGeom) EqualsExact(g1 *GEOSGeom, tolerance float64) (bool, error) {
	eqExact := C.GEOSEqualsExact(g.cGeom, g1.cGeom, C.double(tolerance))
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
func (g *GEOSGeom) Covers(g1 *GEOSGeom) (bool, error) {
	covers := C.GEOSCovers(g.cGeom, g1.cGeom)
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
func (g *GEOSGeom) CoveredBy(g1 *GEOSGeom) (bool, error) {
	coverdBy := C.GEOSCoveredBy(g.cGeom, g1.cGeom)
	return geosBoolResult(coverdBy)
}

// Crosses returns true if this geometry crosses the specified geometry.
func (g *GEOSGeom) Crosses(g1 *GEOSGeom) (bool, error) {
	crosses := C.GEOSCrosses(g.cGeom, g1.cGeom)
	return geosBoolResult(crosses)
}

// GetNumCoordinates returns the count of this Geometrys vertices.
func (g *GEOSGeom) GetNumCoordinates() (int, error) {
	numcoords := C.GEOSGetNumCoordinates(g.cGeom)

	//GEOS return -1 on exception
	if numcoords == C.int(-1) {
		return 0, errors.New(GeosError)
	}

	return int(numcoords), nil

}

// Area returns the area of the geometry
func (g *GEOSGeom) Area() (float64, error) {
	var area C.double

	ret := C.GEOSArea(g.cGeom, &area)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(area), nil
}

// Length returns the length of the geometry
func (g *GEOSGeom) Length() (float64, error) {
	var len C.double

	ret := C.GEOSLength(g.cGeom, &len)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(len), nil
}

// Distance returns the distance between two geometries
func (g *GEOSGeom) Distance(g1 *GEOSGeom) (float64, error) {
	var dist C.double

	ret := C.GEOSDistance(g.cGeom, g1.cGeom, &dist)

	if ret == C.int(-1) {
		return 0.0, errors.New(GeosError)
	}

	return float64(dist), nil
}

/*
NumPoints returns the count of points of the geometry
Geometry type must be a LineString.
*/
func (g *GEOSGeom) NumPoints() (int, error) {

	points := C.GEOSGeomGetNumPoints(g.cGeom)

	if points < 1 {
		return 0, errors.New(GeosError)
	}

	return int(points), nil
}

// GEOSGeomTypeID returns the type ID of the geometry
func (g *GEOSGeom) GEOSGeomTypeID() (int, error) {

	geomTypeID := C.GEOSGeomTypeId(g.cGeom)

	if geomTypeID < 0 {
		return 0, errors.New(GeosError)
	}

	return int(geomTypeID), nil
}

// GetRawGeom returns a pointer to the base C type
func (g *GEOSGeom) GetRawGeom() *C.GEOSGeometry {
	return g.cGeom
}
