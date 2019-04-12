package geos

/*
#cgo CFLAGS: -g -Wall -I/usr/local/include -I/usr/include
#cgo LDFLAGS: -L/usr/local/lib -L/usr/lib -lgeos_c
#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <geos_c.h>

void notice(const char *fmt, ...) {
	va_list ap;
    fprintf( stdout, "NOTICE: ");
	va_start (ap, fmt);
    vfprintf( stdout, fmt, ap);
    va_end(ap);
    fprintf( stdout, "\n" );
}

void log_and_exit(const char *fmt, ...) {
	va_list ap;
    fprintf( stdout, "ERROR: ");
	va_start (ap, fmt);
    vfprintf( stdout, fmt, ap);
    va_end(ap);
    fprintf( stdout, "\n" );
	exit(1);
}

GEOSContextHandle_t ctx;

GEOSContextHandle_t init_geos(){
	ctx = GEOS_init_r();

    // Attach error/ notice handler
    GEOSContext_setErrorHandler_r(ctx, log_and_exit);
	GEOSContext_setNoticeHandler_r(ctx, notice);

	return ctx;
}
*/
import "C"

var (
	ctxHandler C.GEOSContextHandle_t
)

func init() {
	ctxHandler = C.init_geos()
}

// GoFinishGEOS destroys the GEOS lib
func GoFinishGEOS() {
	C.finishGEOS_r(ctxHandler)
}
