package geom

/*
#include <liblwgeom.h>
#include <geos_c.h>
#include "lwgeom_geos.h"
#include "geos.h"
#include "lwgeom.h"
char *cnull = NULL;
*/
import (
	"C"
)
import "unsafe"

// LwGeom Go type to wrap lwgeom
type LwGeom struct {
	LwGeom *C.LWGEOM
}

// LwGeomFromGeoJSON creates lwgeom from GeoJson
func LwGeomFromGeoJSON(json string) *LwGeom {

	lwgeom := C.lwgeom_from_geojson(C.CString(json), &C.cnull)

	return &LwGeom{
		LwGeom: lwgeom,
	}
}

// LwGeomFromGEOS convert GEOS geometry to lwgeom
func LwGeomFromGEOS(geosGeom *C.GEOSGeometry) *LwGeom {
	lwGeom := C.GEOS2LWGEOM(geosGeom, C.uchar(0))

	return &LwGeom{
		LwGeom: lwGeom,
	}
}

// Free clears the memory allocated to lwgeom
func (lwg *LwGeom) Free() {
	C.lwgeom_free(lwg.LwGeom)
}

// LwGeomToGeoJSON generates geojson from lwgeom
func (lwg *LwGeom) LwGeomToGeoJSON(precisoin int, hasBbox int) string {

	geojson := C.lwgeom_to_geojson(lwg.LwGeom, C.cnull, C.int(precisoin), C.int(hasBbox))

	return C.GoString(geojson)
}

// LineSubstring returns a part of the linestring
func (lwg *LwGeom) LineSubstring(from float64, to float64) {

	defer C.lwfree(unsafe.Pointer(lwg.LwGeom))

	newGeom := C.lwgeom_line_substring(lwg.LwGeom, C.double(from), C.double(to))

	lwg.LwGeom = newGeom
}

// SetSRID sets the SRID of the geometry
func (lwg *LwGeom) SetSRID(srid int) {
	C.lwgeom_set_srid(lwg.LwGeom, C.int(srid))
}

// GetSRID returns the SRID of the geometry
func (lwg *LwGeom) GetSRID() int {
	return int(C.lwgeom_get_srid(lwg.LwGeom))
}

// ToGEOS converts lwgeom to GEOS geometry
func (lwg *LwGeom) ToGEOS() *C.GEOSGeometry {
	C.init_geos()
	defer C.finishGEOS()
	return C.LWGEOM2GEOS(lwg.LwGeom, C.uchar(0))
}
