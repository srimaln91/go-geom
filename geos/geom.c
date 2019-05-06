#include "geos.h"

GEOSGeometry*
simplified_buffer(GEOSContextHandle_t handle, GEOSGeometry *g, double width, double tolerance)
{
    GEOSGeometry *simple_geom;
    GEOSGeometry *buffered_geom;

    // Simplify geometry with a given tolerance
    simple_geom = GEOSSimplify_r(handle, g, tolerance);

    //Create buffer
    buffered_geom = GEOSBuffer_r(handle, simple_geom, width, 8);

    //Memory cleanup
    GEOSGeom_destroy_r(handle, simple_geom);

    return buffered_geom;
}
