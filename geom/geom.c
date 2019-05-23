#include <stdio.h>
#include <geos_c.h>
#include <liblwgeom.h>
#include "lwgeom_geos.h"
#include "lwgeom_log.h"

/*
Parse WKT into GEOS geometry
*/
GEOSGeometry *
from_wkt(char *wkt)
{
	GEOSWKTReader *wkt_reader;
	GEOSGeometry *geom;

	wkt_reader = GEOSWKTReader_create();
	geom = GEOSWKTReader_read(wkt_reader, wkt);

	GEOSWKTReader_destroy(wkt_reader);

	return geom;
}

/*
Get WKT from GEOS geometry
*/
char *
to_wkt(GEOSGeometry *g)
{
	GEOSWKTWriter *wkt_writer;
	int output_dimention;

	wkt_writer = GEOSWKTWriter_create();

	// Set output dimention
	output_dimention = GEOSGeom_getCoordinateDimension(g);
	GEOSWKTWriter_setOutputDimension(wkt_writer, output_dimention);

	char *wkt;
	wkt = GEOSWKTWriter_write(wkt_writer, g);

	GEOSWKTWriter_destroy(wkt_writer);

	return wkt;
}

/*
Simplify a given geometry and generate a buffered geometry
*/
GEOSGeometry *
simplified_buffer(GEOSGeometry *g, double width, double tolerance)
{
	GEOSGeometry *simple_geom;
	GEOSGeometry *buffered_geom;

	// Simplify geometry with a given tolerance
	simple_geom = GEOSTopologyPreserveSimplify(g, tolerance);

	//Create buffer
	buffered_geom = GEOSBuffer(simple_geom, width, 8);

	//Memory cleanup
	GEOSGeom_destroy(simple_geom);

	return buffered_geom;
}

/*
Simplify and generate a buffer for a given WKT text.
This function returns WKT.
*/
char *
simplified_buffer_from_wkt(char *inwkt, double width, double tolerance)
{
	GEOSGeometry *in_geom;
	GEOSGeometry *buffered_geom;
	char *out_wkt;

	in_geom = from_wkt(inwkt);
	buffered_geom = simplified_buffer(in_geom, width, tolerance);

	out_wkt = to_wkt(buffered_geom);

	GEOSGeom_destroy(in_geom);
	GEOSGeom_destroy(buffered_geom);

	return out_wkt;
}

/*
Return a linestring being a substring of the input one starting and ending at the given fractions of total 2d length.
Second and third arguments are float8 values between 0 and 1. This only works with LINESTRINGs

Copied from PostGIS source and modified accordingly.
*/
LWGEOM *
lwgeom_line_substring(LWGEOM *ingeom, double from, double to)
{

	LWGEOM *olwgeom = NULL;
	POINTARRAY *ipa, *opa;

	int type = lwgeom_get_type(ingeom);

	if ( from < 0 || from > 1 )
	{
		// elog(ERROR,"line_interpolate_point: 2nd arg isn't within [0,1]");
		return NULL;
	}

	if ( to < 0 || to > 1 )
	{
		// elog(ERROR,"line_interpolate_point: 3rd arg isn't within [0,1]");
		return NULL;
	}

	if ( from > to )
	{
		// elog(ERROR, "2nd arg must be smaller then 3rd arg");
		return NULL;
	}

	if (type == LINETYPE)
	{
		LWLINE *iline = lwgeom_as_lwline(ingeom);

		if (lwgeom_is_empty((LWGEOM *)iline))
		{
			/* TODO return empty line */
			lwline_release(iline);

			return NULL;
		}

		ipa = iline->points;

		opa = ptarray_substring(ipa, from, to, 0);

		if (opa->npoints == 1) /* Point returned */
			olwgeom = (LWGEOM *)lwpoint_construct(iline->srid, NULL, opa);
		else
			olwgeom = (LWGEOM *)lwline_construct(iline->srid, NULL, opa);
	}
	else if (type == MULTILINETYPE)
	{
		LWMLINE *iline;
		uint32_t i = 0, g = 0;
		int homogeneous = LW_TRUE;
		LWGEOM **geoms = NULL;
		double length = 0.0, sublength = 0.0, minprop = 0.0, maxprop = 0.0;

		iline = lwgeom_as_lwmline(ingeom);

		if (lwgeom_is_empty((LWGEOM *)iline))
		{
			/* TODO return empty collection */
			lwmline_release(iline);
		}

		/* Calculate the total length of the mline */
		for (i = 0; i < iline->ngeoms; i++)
		{
			LWLINE *subline = (LWLINE *)iline->geoms[i];
			if (subline->points && subline->points->npoints > 1)
				length += ptarray_length_2d(subline->points);
		}

		geoms = lwalloc(sizeof(LWGEOM *) * iline->ngeoms);

		/* Slice each sub-geometry of the multiline */
		for (i = 0; i < iline->ngeoms; i++)
		{
			LWLINE *subline = (LWLINE *)iline->geoms[i];
			double subfrom = 0.0, subto = 0.0;

			if (subline->points && subline->points->npoints > 1)
				sublength += ptarray_length_2d(subline->points);

			/* Calculate proportions for this subline */
			minprop = maxprop;
			maxprop = sublength / length;

			/* This subline doesn't reach the lowest proportion requested
			   or is beyond the highest proporton */
			if (from > maxprop || to < minprop)
				continue;

			if (from <= minprop)
				subfrom = 0.0;
			if (to >= maxprop)
				subto = 1.0;

			if (from > minprop && from <= maxprop)
				subfrom = (from - minprop) / (maxprop - minprop);

			if (to < maxprop && to >= minprop)
				subto = (to - minprop) / (maxprop - minprop);

			opa = ptarray_substring(subline->points, subfrom, subto, 0);
			if (opa && opa->npoints > 0)
			{
				if (opa->npoints == 1) /* Point returned */
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
		if (!homogeneous)
			type = COLLECTIONTYPE;

		olwgeom = (LWGEOM *)lwcollection_construct(type, iline->srid, NULL, g, geoms);
	}
	else
	{
		return NULL;
	}

	return olwgeom;
}

// Create a buffer around a geometry
LWGEOM *
buffer(LWGEOM *lwg, double width, int quadsegs)
{
	GEOSGeometry *geos_geom;
	GEOSGeometry *buffered_geos;
	LWGEOM *buffered_lwgeom;

	geos_geom = LWGEOM2GEOS(lwg, 0);

	//create buffer
	buffered_geos = GEOSBuffer(geos_geom, width, quadsegs);

	if (!buffered_geos)
		return NULL;

	GEOSGeom_destroy(geos_geom);

	if (lwgeom_has_srid(lwg))
	{
		GEOSSetSRID(buffered_geos, lwgeom_get_srid(lwg));
	}

	buffered_lwgeom = GEOS2LWGEOM(buffered_geos, 0);
	if (!buffered_lwgeom)
		return NULL;

	GEOSGeom_destroy(buffered_geos);

	return buffered_lwgeom;
}

// Create a buffer around a geometry using GEOSBufferParams
LWGEOM *
buffer_with_params(LWGEOM *lwg, double width, GEOSBufferParams *buffer_params)
{
	GEOSGeometry *geos_geom;
	GEOSGeometry *buffered_geos;
	LWGEOM *buffered_lwgeom;

	geos_geom = LWGEOM2GEOS(lwg, 0);

	//create buffer
	buffered_geos = GEOSBufferWithParams(geos_geom, buffer_params, width);

	if (!buffered_geos)
		return NULL;

	GEOSGeom_destroy(geos_geom);

	if (lwgeom_has_srid(lwg))
	{
		GEOSSetSRID(buffered_geos, lwgeom_get_srid(lwg));
	}

	buffered_lwgeom = GEOS2LWGEOM(buffered_geos, 0);
	if (!buffered_lwgeom)
		return NULL;

	GEOSGeom_destroy(buffered_geos);

	return buffered_lwgeom;
}

LWGEOM *
closest_point(LWGEOM *lwg1, LWGEOM *lwg2)
{
	error_if_srid_mismatch(lwg1->srid, lwg2->srid);	

	LWGEOM *lwg_point;
	lwg_point = lwgeom_closest_point(lwg1, lwg2);

	if (lwgeom_is_empty(lwg_point))
	{
		return NULL;
	}
	
	return lwg_point;
}

LWGEOM *
split(LWGEOM *lwg_in, LWGEOM *blade)
{
	LWGEOM *lwgeom_out;

	error_if_srid_mismatch(lwg_in->srid, blade->srid);

	lwgeom_out = lwgeom_split(lwg_in, blade);

	if ( ! lwgeom_out )
	{
		return NULL;
	}

	return lwgeom_out;
}

LWGEOM *
get_subgeom(LWGEOM *lwg, int index)
{
	LWGEOM *sub_geom;
	LWCOLLECTION *collection;

	if ( ! lwgeom_is_collection(lwg) )
	{
		lwerror("Input should be a geometry collection");
		return NULL;
	}

	collection = lwgeom_as_lwcollection(lwg);
	sub_geom = lwcollection_getsubgeom(collection, index);

	if ( ! sub_geom )
	{
		return NULL;
	}

	return sub_geom;
}

double
line_locate_point(LWGEOM *linestring, LWGEOM *point)
{
	LWLINE *lwline;
	LWPOINT *lwpoint;
	POINTARRAY *pa;
	POINT4D p, p_proj;
	double ret;

	if ( linestring->type != LINETYPE )
	{
		lwerror("line_locate_point: 1st arg isn't a line");
		return 0;
	}
	if ( point->type != POINTTYPE )
	{
		lwerror("line_locate_point: 2st arg isn't a point");
		return 0;
	}

	error_if_srid_mismatch(linestring->srid, point->srid);

	lwline = lwgeom_as_lwline(linestring);
	lwpoint = lwgeom_as_lwpoint(point);

	pa = lwline->points;
	lwpoint_getPoint4d_p(lwpoint, &p);

	ret = ptarray_locate_point(pa, &p, NULL, &p_proj);

	return ret;
}

LWGEOM *
geos_union(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;

	GEOSGeometry *geos_union;
	LWGEOM *lwgeom_union;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//create buffer
	geos_union = GEOSUnion(geos_geom1, geos_geom2);

	if (!geos_union)
	{
		return NULL;
	}

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	if (lwgeom_has_srid(lwg1))
	{
		GEOSSetSRID(geos_union, lwgeom_get_srid(lwg1));
	}

	lwgeom_union = GEOS2LWGEOM(geos_union, 0);
	if (!lwgeom_union)
	{
		return NULL;
	}

	GEOSGeom_destroy(geos_union);

	return lwgeom_union;
}

LWGEOM *
geos_intersection(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;

	GEOSGeometry *geos_intersection;
	LWGEOM *lwgeom_intersection;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//create buffer
	geos_intersection = GEOSIntersection(geos_geom1, geos_geom2);

	if (!geos_intersection)
	{
		return NULL;
	}

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	if (lwgeom_has_srid(lwg1))
	{
		GEOSSetSRID(geos_intersection, lwgeom_get_srid(lwg1));
	}

	lwgeom_intersection = GEOS2LWGEOM(geos_intersection, 0);
	if (!lwgeom_intersection)
	{
		return NULL;
	}

	GEOSGeom_destroy(geos_intersection);

	return lwgeom_intersection;
}

// Binary predicate, returns 2 on exception, 1 on true, 0 on false
char
geos_intersects(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char intersects;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	intersects = GEOSIntersects(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return intersects;
}

char
geos_disjoints(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char disjoints;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	disjoints = GEOSDisjoint(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return disjoints;
}

char
geos_touches(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char touches;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	touches = GEOSDisjoint(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return touches;
}

char
geos_within(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char within;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	within = GEOSWithin(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return within;
}

char
geos_contains(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char contains;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	contains = GEOSContains(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return contains;
}

char
geos_overlaps(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char overlaps;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	overlaps = GEOSOverlaps(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return overlaps;
}

char
geos_equals(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char equals;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	equals = GEOSEquals(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return equals;
}

char
geos_equals_exact(LWGEOM *lwg1, LWGEOM *lwg2, double tolerance)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char equals;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	equals = GEOSEqualsExact(geos_geom1, geos_geom2, tolerance);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return equals;
}

char
geos_covers(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char covers;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	covers = GEOSCovers(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return covers;
}

char
geos_covered_by(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char covered_by;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	covered_by = GEOSCoveredBy(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return covered_by;
}

char
geos_crosses(LWGEOM *lwg1, LWGEOM *lwg2)
{
	GEOSGeometry *geos_geom1;
	GEOSGeometry *geos_geom2;
	char crosses;

	geos_geom1 = LWGEOM2GEOS(lwg1, 0);
	geos_geom2 = LWGEOM2GEOS(lwg2, 0);

	//GEOSIntersects return 2 on exception, 1 on true, 0 on false
	crosses = GEOSCrosses(geos_geom1, geos_geom2);

	GEOSGeom_destroy(geos_geom1);
	GEOSGeom_destroy(geos_geom2);

	return crosses;
}
