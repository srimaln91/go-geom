#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>

// Avoid using by accident non _r functions
#ifndef GEOS_USE_ONLY_R_API
#define GEOS_USE_ONLY_R_API
#endif

#include <geos_c.h>

#ifndef GEOS_INIT
#define GEOS_INIT

void notice(const char *fmt, ...);
void log_and_exit(const char *fmt, ...);

GEOSGeometry *from_wkt(GEOSContextHandle_t handle, char *wkt);
char *to_wkt(GEOSContextHandle_t handle, GEOSGeometry *g);
GEOSGeometry *simplified_buffer(GEOSContextHandle_t handle, GEOSGeometry *g, double width, double tolerance);
char *simplified_buffer_from_wkt(GEOSContextHandle_t handle, char *inwkt, double width, double tolerance);

GEOSContextHandle_t ctx;
GEOSContextHandle_t init_geos();

#endif
