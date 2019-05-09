#include <stdio.h>
#include <geos_c.h>
#include <liblwgeom.h>
#include "lwgeom_geos.h"
#include "lwgeom_log.h"

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



/*
Return a linestring being a substring of the input one starting and ending at the given fractions of total 2d length.
Second and third arguments are float8 values between 0 and 1. This only works with LINESTRINGs

Copied from PostGIS source and modified accordingly.
*/
LWGEOM*
lwgeom_line_substring(LWGEOM *ingeom, double from, double to)
{
	// GSERIALIZED *geom = PG_GETARG_GSERIALIZED_P(0);
	// double from = PG_GETARG_FLOAT8(1);
	// double to = PG_GETARG_FLOAT8(2);
	LWGEOM *olwgeom = NULL;
	POINTARRAY *ipa, *opa;
	// GSERIALIZED *ret;
	int type = lwgeom_get_type(ingeom);

	// if ( from < 0 || from > 1 )
	// {
	// 	// elog(ERROR,"line_interpolate_point: 2nd arg isn't within [0,1]");
	// 	return (int*)-1;
	// }

	// if ( to < 0 || to > 1 )
	// {
	// 	// elog(ERROR,"line_interpolate_point: 3rd arg isn't within [0,1]");
	// 	return -1;
	// }

	// if ( from > to )
	// {
	// 	// elog(ERROR, "2nd arg must be smaller then 3rd arg");
	// 	return -1;
	// }

	if ( type == LINETYPE )
	{
		LWLINE *iline = lwgeom_as_lwline(ingeom);

		if ( lwgeom_is_empty((LWGEOM*)iline) )
		{
			/* TODO return empty line */
			lwline_release(iline);
			printf("test1");
			// return -1;
		}

		ipa = iline->points;

		opa = ptarray_substring(ipa, from, to, 0);

		if ( opa->npoints == 1 ) /* Point returned */
			olwgeom = (LWGEOM *)lwpoint_construct(iline->srid, NULL, opa);
		else
			olwgeom = (LWGEOM *)lwline_construct(iline->srid, NULL, opa);

	}
	else if ( type == MULTILINETYPE )
	{
		LWMLINE *iline;
		uint32_t i = 0, g = 0;
		int homogeneous = LW_TRUE;
		LWGEOM **geoms = NULL;
		double length = 0.0, sublength = 0.0, minprop = 0.0, maxprop = 0.0;

		iline = lwgeom_as_lwmline(ingeom);

		if ( lwgeom_is_empty((LWGEOM*)iline) )
		{
			/* TODO return empty collection */
			lwmline_release(iline);
			printf("test2");
			
		}

		/* Calculate the total length of the mline */
		for ( i = 0; i < iline->ngeoms; i++ )
		{
			LWLINE *subline = (LWLINE*)iline->geoms[i];
			if ( subline->points && subline->points->npoints > 1 )
				length += ptarray_length_2d(subline->points);
		}

		geoms = lwalloc(sizeof(LWGEOM*) * iline->ngeoms);

		/* Slice each sub-geometry of the multiline */
		for ( i = 0; i < iline->ngeoms; i++ )
		{
			LWLINE *subline = (LWLINE*)iline->geoms[i];
			double subfrom = 0.0, subto = 0.0;

			if ( subline->points && subline->points->npoints > 1 )
				sublength += ptarray_length_2d(subline->points);

			/* Calculate proportions for this subline */
			minprop = maxprop;
			maxprop = sublength / length;

			/* This subline doesn't reach the lowest proportion requested
			   or is beyond the highest proporton */
			if ( from > maxprop || to < minprop )
				continue;

			if ( from <= minprop )
				subfrom = 0.0;
			if ( to >= maxprop )
				subto = 1.0;

			if ( from > minprop && from <= maxprop )
				subfrom = (from - minprop) / (maxprop - minprop);

			if ( to < maxprop && to >= minprop )
				subto = (to - minprop) / (maxprop - minprop);


			opa = ptarray_substring(subline->points, subfrom, subto, 0);
			if ( opa && opa->npoints > 0 )
			{
				if ( opa->npoints == 1 ) /* Point returned */
				{
					geoms[g] = (LWGEOM *)lwpoint_construct(SRID_UNKNOWN, NULL, opa);
					homogeneous = LW_FALSE;
				}
				else
				{
					geoms[g] = (LWGEOM *)lwline_construct(SRID_UNKNOWN, NULL, opa);
				}
				g++;
			}



		}
		/* If we got any points, we need to return a GEOMETRYCOLLECTION */
		if ( ! homogeneous )
			type = COLLECTIONTYPE;

		olwgeom = (LWGEOM*)lwcollection_construct(type, iline->srid, NULL, g, geoms);
	}
	else
	{
		// elog(ERROR,"line_substring: 1st arg isn't a line");
		// return -1;
		printf("test3");

	}

    return olwgeom;

}
