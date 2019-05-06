#include "geos.h"

/*
Parse WKT into GEOS geometry
*/
GEOSGeometry* 
from_wkt(GEOSContextHandle_t handle, char *wkt)
{
    GEOSWKTReader *wkt_reader;
    GEOSGeometry *geom;

    wkt_reader = GEOSWKTReader_create_r(handle);
    geom = GEOSWKTReader_read_r(handle, wkt_reader, wkt);

    GEOSWKTReader_destroy_r(handle, wkt_reader);  

    return geom;  
}

/*
Get WKT from GEOS geometry
*/
char*
to_wkt(GEOSContextHandle_t handle, GEOSGeometry *g)
{
    GEOSWKTWriter *wkt_writer;
    int output_dimention;

    wkt_writer = GEOSWKTWriter_create_r(handle);

    // Set output dimention
	output_dimention = GEOSGeom_getCoordinateDimension_r(handle, g);
	GEOSWKTWriter_setOutputDimension_r(handle, wkt_writer, output_dimention);

    char *wkt;
    wkt = GEOSWKTWriter_write_r(handle, wkt_writer, g);

    GEOSWKTWriter_destroy_r(handle, wkt_writer);

    return wkt;
}

/*
Simplify a given geometry and generate a buffered geometry
*/
GEOSGeometry*
simplified_buffer(GEOSContextHandle_t handle, GEOSGeometry *g, double width, double tolerance)
{
    GEOSGeometry *simple_geom;
    GEOSGeometry *buffered_geom;

    // Simplify geometry with a given tolerance
    simple_geom = GEOSTopologyPreserveSimplify_r(handle, g, tolerance);

    //Create buffer
    buffered_geom = GEOSBuffer_r(handle, simple_geom, width, 8);

    //Memory cleanup
    GEOSGeom_destroy_r(handle, simple_geom);

    return buffered_geom;
}

/*
Simplify and generate a buffer for a given WKT text.
This function returns WKT.
*/
char* 
simplified_buffer_from_wkt(GEOSContextHandle_t handle, char *inwkt, double width, double tolerance)
{
    GEOSGeometry *in_geom;
    GEOSGeometry *buffered_geom;
    char *out_wkt;

    in_geom = from_wkt(handle, inwkt);
    buffered_geom = simplified_buffer(handle, in_geom, width, tolerance);

    out_wkt = to_wkt(handle, buffered_geom);

    GEOSGeom_destroy_r(handle, in_geom);
    GEOSGeom_destroy_r(handle, buffered_geom);

    return out_wkt;
}
