// Avoid using by accident non _r functions
// #ifndef GEOS_USE_ONLY_R_API
// #define GEOS_USE_ONLY_R_API
// #endif
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <geos_c.h>

#ifndef GEOS_INIT
#define GEOS_INIT

void log_notice(const char *fmt, ...);
void log_error(const char *fmt, ...);

GEOSContextHandle_t ctx;
GEOSContextHandle_t init_geos_r();
void init_geos();
void finish_geos();

#endif
