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

	if bufferedGeom == nil {
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
		return nil, errors.New("Error in union operation")
	}

	return &Geom{
		LwGeom: union,
	}, nil
}
