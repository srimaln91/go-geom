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
	"errors"
	"unsafe"
)

// Geom Go type to wrap lwgeom
type Geom struct {
	LwGeom *C.LWGEOM
}

// SRS Stores required spatial reference systems
var SRS = map[string]string{
	"EPSG:4326": "+proj=longlat +datum=WGS84 +no_defs",
	"EPSG:4978": "+proj=geocent +datum=WGS84 +units=m +no_defs",
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
func FromGeoJSON(geojson string) (*Geom, error) {

	geojsonCstring := C.CString(geojson)
	defer C.lwfree(unsafe.Pointer(geojsonCstring))

	lwgeom := C.lwgeom_from_geojson(geojsonCstring, &C.cnull)
	if lwgeom == nil {
		return nil, errors.New("Lwgeom exception on lwgeom_from_geojson")
	}

	return &Geom{
		LwGeom: lwgeom,
	}, nil
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
func (lwg *Geom) ToGeoJSON(precisoin int, hasBbox int) (string, error) {

	geojson := C.lwgeom_to_geojson(lwg.LwGeom, C.cnull, C.int(precisoin), C.int(hasBbox))
	if geojson == nil {
		return "", errors.New("Lwgeom exception on lwgeom_to_geojson")
	}
	defer C.lwfree(unsafe.Pointer(geojson))

	return C.GoString(geojson), nil
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

// ClosestPoint returns the 2-dimensional point on g1 that is closest to g2. This is the first point of the shortest line.
func (lwg *Geom) ClosestPoint(lwg2 *Geom) (*Geom, error) {

	lwgeomPoint := C.closest_point(lwg.LwGeom, lwg2.LwGeom)

	if lwgeomPoint == nil {
		return nil, errors.New("Lwgeom Error")
	}

	return &Geom{
		LwGeom: lwgeomPoint,
	}, nil
}

/*
Split supports splitting a line by (multi)point, (multi)line or (multi)polygon boundary, a (multi)polygon by line.
The returned geometry is always a collection.
*/
func (lwg *Geom) Split(blade *Geom) (*Geom, error) {
	lwgeomOut := C.split(lwg.LwGeom, blade.LwGeom)

	if lwgeomOut == nil {
		return nil, errors.New("LWgeom Error")
	}

	return &Geom{
		LwGeom: lwgeomOut,
	}, nil
}

// GetSubGeom returns the subgeom from given index from a geometry collection
func (lwg *Geom) GetSubGeom(index int) (*Geom, error) {
	lwgeomOut := C.get_subgeom(lwg.LwGeom, C.int(index))

	if lwgeomOut == nil {
		return nil, errors.New("Lwgeom Error")
	}

	return &Geom{
		LwGeom: lwgeomOut,
	}, nil
}

// Equals returns true if a given geometry equals with another geometry
func (lwg *Geom) Equals(g2 *Geom) bool {
	isSame := C.lwgeom_same(lwg.LwGeom, g2.LwGeom)

	return C.int(isSame) == 1
}

/*
LineLocatePoint returns a float between 0 and 1 representing the location of the closest point on LineString
to the given Point, as a fraction of total 2d line length.
https://postgis.net/docs/manual-1.5/ST_Line_Locate_Point.html
*/
func (lwg *Geom) LineLocatePoint(point *Geom) (float64, error) {
	ret := C.line_locate_point(lwg.LwGeom, point.LwGeom)

	if ret == 0 {
		return 0, errors.New("Invalid geometries")
	}

	return float64(ret), nil
}

// CreateFromWKT parse WKT and create a Geom object
func CreateFromWKT(wkt string) (*Geom, error) {
	cWkt := C.CString(wkt)
	defer C.lwfree(unsafe.Pointer(cWkt))

	lwgeom := C.lwgeom_from_wkt(cWkt, C.LW_PARSER_CHECK_NONE)

	if lwgeom == nil {
		return nil, errors.New("Error parsing WKT")
	}

	return &Geom{
		LwGeom: lwgeom,
	}, nil
}

// ToWKT returns a WKT string
func (lwg *Geom) ToWKT(precision int) ([]byte, error) {

	var size C.size_t
	wkt := C.lwgeom_to_wkt(lwg.LwGeom, C.WKT_ISO, C.int(precision), &size)

	if wkt == nil {
		return nil, errors.New("Error creating WKT")
	}

	return []byte(C.GoString(wkt)), nil
}

// Area returns the area of the geometry
func (lwg *Geom) Area() (float64, error) {
	area := C.lwgeom_area(lwg.LwGeom)
	return float64(area), nil
}

// Length returns the length of the geometry
func (lwg *Geom) Length() (float64, error) {
	length := C.lwgeom_length(lwg.LwGeom)
	return float64(length), nil
}
