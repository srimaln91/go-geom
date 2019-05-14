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

import (
	"unsafe"
)

// LwGeom Go type to wrap lwgeom
type LwGeom struct {
	LwGeom *C.LWGEOM
}

// LwGeomFromGeoJSON creates lwgeom from GeoJson
func LwGeomFromGeoJSON(geojson string) *LwGeom {

	geojsonCstring := C.CString(geojson)
	lwgeom := C.lwgeom_from_geojson(geojsonCstring, &C.cnull)
	defer C.lwfree(unsafe.Pointer(geojsonCstring))

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

// ToGeoJSON generates geojson from lwgeom
func (lwg *LwGeom) ToGeoJSON(precisoin int, hasBbox int) string {

	geojson := C.lwgeom_to_geojson(lwg.LwGeom, C.cnull, C.int(precisoin), C.int(hasBbox))
	defer C.lwfree(unsafe.Pointer(geojson))
	return C.GoString(geojson)
}

// LineSubstring returns a part of the linestring
func (lwg *LwGeom) LineSubstring(from float64, to float64) {

	defer C.lwgeom_free(lwg.LwGeom)

	lwg.LwGeom = C.lwgeom_line_substring(lwg.LwGeom, C.double(from), C.double(to))

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
func (lwg *LwGeom) ToGEOS() *GEOSGeom {

	return GenerateGeosGeom(C.LWGEOM2GEOS(lwg.LwGeom, C.uchar(0)))
}

/*
Project transforms (reproject) a geometry from one SRS to another.
You will have to use the WKT versions of SRS definition.
references: https://epsg.io
*/
func (lwg *LwGeom) Project(fromSRS string, toSRS string) {
	fromSrsCtype := C.CString(fromSRS)
	defer C.free(unsafe.Pointer(fromSrsCtype))

	toSrsCtype := C.CString(toSRS)
	defer C.free(unsafe.Pointer(toSrsCtype))

	from := C.lwproj_from_string(fromSrsCtype)
	defer C.pj_free(from)

	to := C.lwproj_from_string(toSrsCtype)
	defer C.pj_free(to)

	C.lwgeom_transform(lwg.LwGeom, from, to)
}
