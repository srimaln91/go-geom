package geom

import (
	"testing"
)

func TestGeomFromGeoJson(t *testing.T) {

	geom := FromGeoJSON(JSONLinestring)

	if geom == nil {
		t.Error("Error: GeomFromGeoJson()")
	}

	geom.Free()
}

func TestToGeoJson(t *testing.T) {

	geom := FromGeoJSON(JSONLinestring)
	jsonString := geom.ToGeoJSON(4, 0)

	if jsonString == "" {
		t.Error("Error: LwGeomToGeoJson()")
	}

	geom.Free()
}

func TestLineSubstring(t *testing.T) {

	expectedJSON := `{"type":"LineString","coordinates":[[79.9066,6.8597],[79.9073,6.859],[79.9076,6.8588],[79.9078,6.8585],[79.908,6.8582],[79.9083,6.858],[79.9084,6.8579],[79.9085,6.8578],[79.9088,6.8575],[79.9089,6.8573],[79.9089,6.8573]]}`

	geom := FromGeoJSON(JSONLinestring)
	geom.LineSubstring(0.5, 0.52)
	resultJSON := geom.ToGeoJSON(4, 0)

	if resultJSON != expectedJSON {
		t.Error("Error: LineSubstring()", resultJSON)
	}

	geom.Free()
}

func TestToGEOS(t *testing.T) {

	geom := FromGeoJSON(JSONLinestring)
	geos := geom.ToGEOS()

	coords, _ := geos.GetNumCoordinates()

	if coords == 0 {
		t.Error("Error: LwGeomToGeoJson()")
	}

	geom.Free()
	geos.Destroy()
}

func TestLwGeomFromGEOS(t *testing.T) {

	geos := FromWKT(WKTLinestring)
	LwGeom := LwGeomFromGEOS(geos.cGeom)

	coords, _ := geos.GetNumCoordinates()

	if coords == 0 {
		t.Error("Error: LwGeomToGeoJson()")
	}

	LwGeom.Free()
	geos.Destroy()
}

func TestProject(t *testing.T) {

	geom := FromGeoJSON(JSONLinestring)
	geom.SetSRID(4326)

	fromSRS := SRS["EPSG:4326"]
	toSRS := SRS["EPSG:3857"]

	geom.Project(fromSRS, toSRS)
	geom.SetSRID(3857)

	geos := geom.ToGEOS()
	geos.Simplify(0.001)
	geos.Buffer(200.00)

	geom2 := LwGeomFromGEOS(geos.cGeom)
	geom2.Project(toSRS, fromSRS)

	geom2.SetSRID(4326)

	if geom2.GetSRID() != 4326 {
		t.Error("Error: Project()")
	}

	geom.Free()
	geos.Destroy()
}

func TestLwGeomVersion(t *testing.T) {
	version := LwGeomVersion()

	if version == "" {
		t.Error("Error: LwGeomVersion()")
	}
}

func TestGEOSVersion(t *testing.T) {
	version := GEOSVersion()

	if version == "" {
		t.Error("Error: GEOSVersion()")
	}
}

func TestClosestPoint(t *testing.T) {
	geom1 := FromGeoJSON(JSONLinestring)
	geom2 := FromGeoJSON(`{
        "type": "Point",
        "coordinates": [
          79.92603331804276,
          6.84914291895139
        ]
	  }`)

	closestPoint, err := geom1.ClosestPoint(geom2)

	if err != nil || closestPoint == nil {
		t.Error("Error: ClosestPoint()")
	}

	geoJSON := closestPoint.ToGeoJSON(6, 0)

	if geoJSON != `{"type":"Point","coordinates":[79.925546,6.848402]}` {
		t.Error("Error: ClosestPoint()")
	}

	geom1.Free()
	geom2.Free()
	closestPoint.Free()

}

func TestSplit(t *testing.T) {
	geom1 := FromGeoJSON(GetFileContents("../testdata/split/source.json"))
	blade := FromGeoJSON(GetFileContents("../testdata/split/blade.json"))

	// expectedResult := GetFileContents("../testdata/split/result.json")

	collection, _ := geom1.Split(blade)

	if collection == nil {
		t.Error("Error: SplitAndSubGeom()")
	}

}

func TestSubGeom(t *testing.T) {
	geom1 := FromGeoJSON(GetFileContents("../testdata/split/source.json"))
	geom1.SetSRID(4326)

	blade := FromGeoJSON(GetFileContents("../testdata/split/blade.json"))
	blade.SetSRID(4326)

	expectedResult := FromGeoJSON(GetFileContents("../testdata/split/result.json"))

	collection, _ := geom1.Split(blade)

	selectedGeom, _ := collection.GetSubGeom(0)

	if selectedGeom.ToGeoJSON(4, 0) != expectedResult.ToGeoJSON(4, 0) {
		t.Error("Error: SplitAndSubGeom()")
	}
}

func TestLwGeomEquals(t *testing.T) {
	geom := FromGeoJSON(JSONLinestring)

	if !geom.Equals(geom) {
		t.Error("Error: Equals()")
	}
}

func TestLineLocatePoint(t *testing.T) {
	linestring := FromGeoJSON(JSONLinestring)
	defer linestring.Free()

	point := FromGeoJSON(`{
        "type": "Point",
        "coordinates": [
          79.91254448890686,
          6.856021573114641
        ]
	}`)
	defer point.Free()

	ret, err := linestring.LineLocatePoint(point)

	if err != nil {
		t.Error("Error: Equals()", err)
	}

	t.Log("Line locate point: ", ret)
}
