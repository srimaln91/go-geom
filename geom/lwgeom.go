package geom

/*
#include <liblwgeom.h>
#include <geos_c.h>
#include "lwgeom_geos.h"
#include "geos.h"
#include "geom.h"
char *cnull = NULL;
*/
import (
	"C"
)

import (
	"unsafe"
)

// Geom Go type to wrap lwgeom
type Geom struct {
	LwGeom *C.LWGEOM
}

// SRS Stores required spatial reference systems
var SRS = map[string]string{
	"EPSG:4326": "+proj=longlat +datum=WGS84 +no_defs",
	"EPSG:3857": "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs",
}

// LwGeomVersion the version number of liblwgeom
func LwGeomVersion() string {
	version := C.lwgeom_version()
	
	return C.GoString(version)
}

// GEOSVersion returns the version number of libgeos
func GEOSVersion() string {
	geosVersion := C.lwgeom_geos_version()

	return C.GoString(geosVersion)
}

// FromGeoJSON creates lwgeom from GeoJson
func FromGeoJSON(geojson string) *Geom {

	geojsonCstring := C.CString(geojson)
	lwgeom := C.lwgeom_from_geojson(geojsonCstring, &C.cnull)
	defer C.lwfree(unsafe.Pointer(geojsonCstring))

	return &Geom{
		LwGeom: lwgeom,
	}
}

// LwGeomFromGEOS convert GEOS geometry to lwgeom
func LwGeomFromGEOS(geosGeom *C.GEOSGeometry) *Geom {

	lwGeom := C.GEOS2LWGEOM(geosGeom, C.uchar(0))

	return &Geom{
		LwGeom: lwGeom,
	}
}

// Free clears the memory allocated to lwgeom
func (lwg *Geom) Free() {
	C.lwgeom_free(lwg.LwGeom)
}

// ToGeoJSON generates geojson from lwgeom
func (lwg *Geom) ToGeoJSON(precisoin int, hasBbox int) string {

	geojson := C.lwgeom_to_geojson(lwg.LwGeom, C.cnull, C.int(precisoin), C.int(hasBbox))
	defer C.lwfree(unsafe.Pointer(geojson))
	return C.GoString(geojson)
}

// LineSubstring returns a part of the linestring
func (lwg *Geom) LineSubstring(from float64, to float64) {

	defer C.lwgeom_free(lwg.LwGeom)

	lwg.LwGeom = C.lwgeom_line_substring(lwg.LwGeom, C.double(from), C.double(to))

}

// SetSRID sets the SRID of the geometry
func (lwg *Geom) SetSRID(srid int) {
	C.lwgeom_set_srid(lwg.LwGeom, C.int(srid))
}

// GetSRID returns the SRID of the geometry
func (lwg *Geom) GetSRID() int {
	return int(C.lwgeom_get_srid(lwg.LwGeom))
}

// ToGEOS converts lwgeom to GEOS geometry
func (lwg *Geom) ToGEOS() *GEOSGeom {

	return GenerateGeosGeom(C.LWGEOM2GEOS(lwg.LwGeom, C.uchar(0)))
}

/*
Project transforms (reproject) a geometry from one SRS to another.
You will have to use the WKT versions of SRS definition.
references: https://epsg.io
*/
func (lwg *Geom) Project(fromSRS string, toSRS string) {
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
