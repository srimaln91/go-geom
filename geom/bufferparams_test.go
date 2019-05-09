package geom

import (
	"testing"
)

func TestBufferWithParams(t *testing.T) {
	bParams, _ := CreateBufferParams()
	defer bParams.Destroy()

	quadsegs := 8
	width := 10.0

	bParams.SetEndCapStyle(GeosbufCapRound)
	bParams.SetJoinStyle(GeosbufJoinRound)
	bParams.SetQuadrantSegments(quadsegs)
	bParams.SetMitreLimit(5)

	geom := FromWKT("POINT(0 0)")

	geom.BufferWithParams(width, bParams)

	ptCount, _ := geom.GetNumCoordinates()

	if geom == nil || ptCount != 33 {
		t.Errorf("Error: CreateBufferWithParams()")
	}

	geom.Destroy()
}

func TestSetSingleSided(t *testing.T) {
	bParams, _ := CreateBufferParams()
	defer bParams.Destroy()

	width := 10.0

	bParams.SetSingleSided()

	geom := FromWKT("LINESTRING(0 0, 1 1)")

	geom.BufferWithParams(width, bParams)

	ptCount, _ := geom.GetNumCoordinates()

	if geom == nil || ptCount != 5 {
		t.Errorf("Error: SetSingleSided()")
	}

	geom.Destroy()
}
