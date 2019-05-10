package geom

import (
	"testing"
)

func TestGeomFromGeoJson(t *testing.T) {

	geom := LwGeomFromGeoJSON(JSONLinestring)

	if geom == nil {
		t.Error("Error: GeomFromGeoJson()")
	}

	geom.Free()
}

func TestLwGeomToGeoJson(t *testing.T) {

	geom := LwGeomFromGeoJSON(JSONLinestring)
	jsonString := geom.LwGeomToGeoJSON(4, 0)

	if jsonString == "" {
		t.Error("Error: LwGeomToGeoJson()")
	}

	geom.Free()
}

func TestLineSubstring(t *testing.T) {

	expectedJSON := `{"type":"LineString","coordinates":[[79.9066,6.8597],[79.9073,6.859],[79.9076,6.8588],[79.9078,6.8585],[79.908,6.8582],[79.9083,6.858],[79.9084,6.8579],[79.9085,6.8578],[79.9088,6.8575],[79.9089,6.8573],[79.9089,6.8573]]}`

	geom := LwGeomFromGeoJSON(JSONLinestring)
	geom.LineSubstring(0.5, 0.52)
	resultJSON := geom.LwGeomToGeoJSON(4, 0)

	if resultJSON != expectedJSON {
		t.Error("Error: LineSubstring()", resultJSON)
	}

	geom.Free()
}

func TestToGEOS(t *testing.T) {

	geom := LwGeomFromGeoJSON(JSONLinestring)
	geos := geom.ToGEOS()

	if geos == nil {
		t.Error("Error: LwGeomToGeoJson()")
	}

	geom.Free()
}
