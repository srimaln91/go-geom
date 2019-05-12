#include <liblwgeom.h>
#include <geos_c.h>

#ifndef LWGEOM_INIT
#define LWGEOM_INIT

LWGEOM* lwgeom_line_substring(LWGEOM *ingeom, double from, double to);

#endif