#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include "geos.h"

void log_notice(const char *fmt, ...)
{
    va_list ap;
    fprintf(stdout, "NOTICE: ");
    va_start(ap, fmt);
    vfprintf(stdout, fmt, ap);
    va_end(ap);
    fprintf(stdout, "\n");
}

void log_error(const char *fmt, ...)
{
    va_list ap;
    fprintf(stdout, "ERROR: ");
    va_start(ap, fmt);
    vfprintf(stdout, fmt, ap);
    va_end(ap);
    fprintf(stdout, "\n");
}

GEOSContextHandle_t ctx;

// Thread safe versions of GEOS initialization and destruct functions
GEOSContextHandle_t init_geos_r()
{
    ctx = GEOS_init_r();

    // Attach error/ notice handler
    GEOSContext_setErrorHandler_r(ctx, log_error);
    GEOSContext_setNoticeHandler_r(ctx, log_notice);

    return ctx;
}

/*
Initialization/ destruction function for non R version which will be used by liblwgeom
*/
void init_geos()
{
	initGEOS(log_notice, log_error);
}

void finish_geos()
{
    finishGEOS();
}
