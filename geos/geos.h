#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <geos_c.h>

void notice(const char *fmt, ...);

void log_and_exit(const char *fmt, ...);

GEOSContextHandle_t ctx;

GEOSContextHandle_t init_geos();
