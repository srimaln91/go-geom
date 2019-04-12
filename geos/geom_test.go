package geos

import "testing"

func TestFromWKT(t *testing.T) {
	wkt := "POINT (0 0)"
	geom := FromWKT(wkt)

	if geom.cGeom == nil {
		t.Errorf("Error: CreateFromWKT(%q) error", wkt)
	}
}

func TestBuffer(t *testing.T) {

	geom := FromWKT("POINT (0 0)")
	geom.Buffer(1)

	if geom.cGeom == nil {
		t.Errorf("Error: Buffer() error")
	}

}

func TestToWKT(t *testing.T) {

	wkt := "POINT (0 0)"
	expectedWKT := "POINT (0.0000000000000000 0.0000000000000000)"

	geom := FromWKT(wkt)

	if geom.ToWKT() != expectedWKT {
		t.Errorf("Error: ToWKT(%s) error", wkt)
	}
}
