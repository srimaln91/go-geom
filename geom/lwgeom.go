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

type LwGeom struct {
	LwGeom *C.LWGEOM
}

func LwGeomFromGeoJson(json string) *LwGeom {

	lwgeom := C.lwgeom_from_geojson(C.CString(json), &C.cnull)

	return &LwGeom{
		LwGeom: lwgeom,
	}
}

func LwGeomFromGEOS(geosGeom *C.GEOSGeometry) *LwGeom {
	lwGeom := C.GEOS2LWGEOM(geosGeom, C.uchar(0))

	return &LwGeom{
		LwGeom: lwGeom,
	}
}

func (lwg *LwGeom) Free() {
	C.lwgeom_free(lwg.LwGeom)
}

func (lwg *LwGeom) LwGeomToGeoJson(precisoin int, hasBbox int) string {

	geojson := C.lwgeom_to_geojson(lwg.LwGeom, C.cnull, C.int(precisoin), C.int(hasBbox))

	return C.GoString(geojson)
}

func (lwg *LwGeom) LineSubstring(from float64, to float64) {

	defer C.lwfree(unsafe.Pointer(lwg.LwGeom))

	newGeom := C.lwgeom_line_substring(lwg.LwGeom, C.double(from), C.double(to))

	lwg.LwGeom = newGeom
}

func (lwg *LwGeom) SetSRID(srid int) {
	C.lwgeom_set_srid(lwg.LwGeom, C.int(srid))
}

func (lwg *LwGeom) GetSRID() int {
	return int(C.lwgeom_get_srid(lwg.LwGeom))
}

func (lwg *LwGeom) ToGEOS() *C.GEOSGeometry {
	C.init_geos()
	defer C.finishGEOS()
	return C.LWGEOM2GEOS(lwg.LwGeom, C.uchar(0))
}
