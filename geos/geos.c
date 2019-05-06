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

GEOSContextHandle_t init_geos()
{
    ctx = GEOS_init_r();

    // Attach error/ notice handler
    GEOSContext_setErrorHandler_r(ctx, log_error);
    GEOSContext_setNoticeHandler_r(ctx, log_notice);

    return ctx;
}
