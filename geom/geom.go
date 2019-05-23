package geom

/*
#include <liblwgeom.h>
#include <geos_c.h>
#include "lwgeom_geos.h"
#include "geos.h"
#include "geom.h"
*/
import (
	"C"
)
import "errors"

// Buffer creates a buffer around a geometry object
func (lwg *Geom) Buffer(width float64) error {
	bufferedGeom := C.buffer(lwg.LwGeom, C.double(width), C.int(8))
	defer C.lwgeom_free(lwg.LwGeom)

	if bufferedGeom == nil { //Cleanup
		// cleanup(geom1, geom2, geom3, geom4)
		return errors.New("Error creating Buffer")
	}

	lwg.LwGeom = bufferedGeom

	return nil
}

// BufferWithParams creates a buffer around a geometry using BufferParams object
func (lwg *Geom) BufferWithParams(params *BufferParams, width float64) error {
	bufferedGeom := C.buffer_with_params(lwg.LwGeom, C.double(width), params.CBufP)
	defer C.lwgeom_free(lwg.LwGeom)

	if bufferedGeom == nil {
		return errors.New("Error creating Buffer")
	}

	lwg.LwGeom = bufferedGeom

	return nil
}

// Union returns the union of two geometries
func (lwg *Geom) Union(g1 *Geom) (*Geom, error) {

	union := C.geos_union(lwg.LwGeom, g1.LwGeom)

	if union == nil {
		return nil, errors.New("Error in GEOS union operation")
	}

	return &Geom{
		LwGeom: union,
	}, nil
}

// Intersection finds the intersection of two geometries
func (lwg *Geom) Intersection(g1 *Geom) (*Geom, error) {

	intersection := C.geos_intersection(lwg.LwGeom, g1.LwGeom)

	if intersection == nil {
		return nil, errors.New("Error in GEOS intersection operation")
	}

	return &Geom{
		LwGeom: intersection,
	}, nil
}

// Intersects checks whether the geom intersects with an another geom
func (lwg *Geom) Intersects(g1 *Geom) (bool, error) {
	intersects := C.geos_intersects(lwg.LwGeom, g1.LwGeom)

	if intersects == C.char(2) {
		return false, errors.New("Error in GEOS intersects operation")
	}

	return geosBoolResult(intersects)
}

/*
Disjoints Overlaps, Touches, Within all imply geometries are not spatially disjoint.
If any of the aforementioned returns true, then the geometries are not spatially disjoint.
Disjoint implies false for spatial intersection.
*/
func (lwg *Geom) Disjoints(g1 *Geom) (bool, error) {
	disjoints := C.geos_disjoints(lwg.LwGeom, g1.LwGeom)

	if disjoints == C.char(2) {
		return false, errors.New("Error in GEOS disjoint operation")
	}

	return geosBoolResult(disjoints)
}

// Touches returns TRUE if the only points in common between g1 and g2 lie in the union of the boundaries of g1 and g2.
func (lwg *Geom) Touches(g1 *Geom) (bool, error) {
	touches := C.geos_disjoints(lwg.LwGeom, g1.LwGeom)

	if touches == C.char(2) {
		return false, errors.New("Error in GEOS touches operation")
	}

	return geosBoolResult(touches)
}

// Within returns TRUE if geometry A is completely inside geometry B
func (lwg *Geom) Within(g1 *Geom) (bool, error) {
	within := C.geos_within(lwg.LwGeom, g1.LwGeom)

	if within == C.char(2) {
		return false, errors.New("Error in GEOS within operation")
	}

	return geosBoolResult(within)
}

// Contains returns TRUE if geometry B is completely inside geometry A.
func (lwg *Geom) Contains(g1 *Geom) (bool, error) {
	contains := C.geos_contains(lwg.LwGeom, g1.LwGeom)

	if contains == C.char(2) {
		return false, errors.New("Error in GEOS contains operation")
	}

	return geosBoolResult(contains)
}

// Overlaps returns TRUE if geometry B is completely inside geometry A.
func (lwg *Geom) Overlaps(g1 *Geom) (bool, error) {
	overlaps := C.geos_overlaps(lwg.LwGeom, g1.LwGeom)

	if overlaps == C.char(2) {
		return false, errors.New("Error in GEOS overlaps operation")
	}

	return geosBoolResult(overlaps)
}

/*
GEOSEquals returns true if the DE-9IM intersection matrix for the two Geometrys is T*F**FFF*.

a and b are topologically equal. "Two geometries are topologically equal if their interiors
intersect and no part of the interior or boundary of one geometry intersects the exterior of the other".[9]

equals to Within & Contains
*/
func (lwg *Geom) GEOSEquals(g1 *Geom) (bool, error) {
	equals := C.geos_equals(lwg.LwGeom, g1.LwGeom)

	if equals == C.GEOS_EXCEPTION {
		return false, errors.New("Error in GEOS equals operation")
	}

	return geosBoolResult(equals)
}

// GEOSEqualsExact returns true if the two Geometrys are of the same type and their
// vertices corresponding by index are equal up to a specified tolerance.
func (lwg *Geom) GEOSEqualsExact(g1 *Geom, tolerance float64) (bool, error) {
	eqExact := C.geos_equals_exact(lwg.LwGeom, g1.LwGeom, C.double(tolerance))

	if C.GEOS_EXCEPTION == eqExact {
		return false, errors.New("Error in GEOS equals operation")
	}

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
func (lwg *Geom) Covers(g1 *Geom) (bool, error) {
	covers := C.geos_covers(lwg.LwGeom, g1.LwGeom)

	if covers == C.GEOS_EXCEPTION {
		return false, errors.New("Error in GEOS equals operation")
	}

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
func (lwg *Geom) CoveredBy(g1 *Geom) (bool, error) {
	coveredBy := C.geos_covered_by(lwg.LwGeom, g1.LwGeom)

	if coveredBy == C.GEOS_EXCEPTION {
		return false, errors.New("Error in GEOS equals operation")
	}

	return geosBoolResult(coveredBy)
}

// Crosses returns true if this geometry crosses the specified geometry.
func (lwg *Geom) Crosses(g1 *Geom) (bool, error) {
	crosses := C.geos_crosses(lwg.LwGeom, g1.LwGeom)

	if crosses == C.GEOS_EXCEPTION {
		return false, errors.New("Error in GEOS equals operation")
	}

	return geosBoolResult(crosses)
}
