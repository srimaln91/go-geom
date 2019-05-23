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

func TestGEOSUnion(t *testing.T) {

	lwgeom1 := FromGeoJSON(`{
        "type": "Polygon",
        "coordinates": [
          [
            [
              79.87292289733887,
              6.9204625096823325
            ],
            [
              79.89575386047363,
              6.9204625096823325
            ],
            [
              79.89575386047363,
              6.9409965539246095
            ],
            [
              79.87292289733887,
              6.9409965539246095
            ],
            [
              79.87292289733887,
              6.9204625096823325
            ]
          ]
        ]
      }
	`)

	lwgeom2 := FromGeoJSON(`{
        "type": "Polygon",
        "coordinates": [
          [
            [
              79.8874282836914,
              6.93153907621163
            ],
            [
              79.90399360656738,
              6.93153907621163
            ],
            [
              79.90399360656738,
              6.951135440565165
            ],
            [
              79.8874282836914,
              6.951135440565165
            ],
            [
              79.8874282836914,
              6.93153907621163
            ]
          ]
        ]
	  }`)

	union, err := lwgeom1.Union(lwgeom2)

	if err != nil || union == nil {
		t.Error("Error: Union()")
	}

	lwgeom1.Free()
	lwgeom2.Free()
}

func TestCreateFromWKT(t *testing.T) {
	wkt := "LINESTRING(0 0,10 10)"
	geom, err := CreateFromWKT(wkt)
	defer geom.Free()

	if err != nil {
		t.Error("CreaeteFromWKT():", err)
	}
}

func TestLwgeomToWKT(t *testing.T) {
	wkt := "LINESTRING(0 0,10 10)"
	geom, err := CreateFromWKT(wkt)
	defer geom.Free()

	if err != nil {
		t.Error("ToWKT():", err)
	}

	wktResult, _ := geom.ToWKT(0)

	if string(wktResult) != wkt {
		t.Error("Error: ToWKT()")
	}
}

func TestLwGeomIntersection(t *testing.T) {

	geom1, _ := CreateFromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")
	defer geom2.Free()

	expectredWKT := "POLYGON((79.856061 6.911495,79.856116 6.911662,79.856234 6.911506,79.856184 6.911368,79.856097 6.911324,79.856061 6.911495))"

	intersection, err := geom1.Intersection(geom2)

	if err != nil || intersection == nil {
		t.Error("Error: Intersection() ", err)
	}

	resultWKT, _ := intersection.ToWKT(6)

	if string(resultWKT) != expectredWKT {
		t.Error("Error: Intersection()")
	}
}

func TestLWGeomIntersects(t *testing.T) {

	geom1, _ := CreateFromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")
	defer geom2.Free()

	intersects, err := geom1.Intersects(geom2)

	if err != nil {
		t.Error("Error: Intersects", err)
	}

	if intersects == false {
		t.Error("Error: Intersects()")
	}
}

func TestLwGeomDisjoint(t *testing.T) {

	geom1, _ := CreateFromWKT("POINT(0 0)")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("LINESTRING ( 2 0, 0 2 )")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("LINESTRING ( 0 0, 0 2 )")
	defer geom3.Free()

	disjoint, _ := geom1.Disjoints(geom2)

	if disjoint == false {
		t.Errorf("Error: Disjoint")
	}

	disjoint, _ = geom1.Disjoints(geom3)

	if disjoint == true {
		t.Errorf("Error: Disjoint")
	}
}

func TestLWGeomTouches(t *testing.T) {
	geom1, _ := CreateFromWKT("POINT(1 1)")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POINT(0 2)")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("LINESTRING(0 0, 1 1, 0 2)")
	defer geom3.Free()

	touches, _ := geom1.Touches(geom3)

	if touches == true {
		t.Errorf("Error: Touches")
	}

	touches, _ = geom2.Touches(geom1)

	if touches == false {
		t.Errorf("Error: Touches")
	}
}

func TestLwGeomWithin(t *testing.T) {

	geom1, _ := CreateFromWKT("POLYGON((20 15,23 11,26 16,20 15))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((8 52,21 50,11 46,8 52))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((10 10,25 7,40 10,33 19,24 26,15 20,10 10))")
	defer geom3.Free()

	within, _ := geom1.Within(geom3)

	if within == false {
		t.Errorf("Error: Within")
	}

	within, _ = geom2.Within(geom3)

	if within == true {
		t.Errorf("Error: Within")
	}
}

func TestLwGeomContains(t *testing.T) {

	geom1, _ := CreateFromWKT("POLYGON((20 15,23 11,26 16,20 15))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((8 52,21 50,11 46,8 52))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((10 10,25 7,40 10,33 19,24 26,15 20,10 10))")
	defer geom3.Free()

	contains, _ := geom3.Contains(geom1)

	if contains == false {
		t.Errorf("Error: Contains")
	}

	contains, _ = geom3.Contains(geom2)

	if contains == true {
		t.Errorf("Error: Contains")
	}
}

func TestLwGeomOverlaps(t *testing.T) {
	geom1, _ := CreateFromWKT("POLYGON((17 19,28 19,28 8,17 8,17 19))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((24 22,41 22,41 8,24 8,24 22))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((57 31,61 31,61 26,57 26,57 31))")
	defer geom3.Free()

	isOVerlaps, _ := geom1.Overlaps(geom2)

	if isOVerlaps == false {
		t.Errorf("Error: Overlaps")
	}

	isOVerlaps, _ = geom3.Overlaps(geom1)

	if isOVerlaps == true {
		t.Errorf("Error: Overlaps")
	}
}

func TestLwGeomGEOSEquals(t *testing.T) {
	geom1, _ := CreateFromWKT("POLYGON((27 19,38 19,38 9,27 9,27 19))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((27 19,38 19,38 9,27 9,27 19))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((26 21,41 18,40 6,25 8,26 21))")
	defer geom3.Free()

	isequal, _ := geom1.GEOSEquals(geom2)

	if !isequal {
		t.Errorf("Error: Equals")
	}

	isequal, _ = geom1.GEOSEquals(geom3)

	if isequal {
		t.Errorf("Error: Equals")
	}
}

func TestLwGeomEqualsExact(t *testing.T) {
	geom1, _ := CreateFromWKT("LINESTRING(0 0, 1 1, 0 2)")
	geom2, _ := CreateFromWKT("LINESTRING(0 0, 1 1.001, 0 2)")
	// geom3 := FromWKT("POLYGON((26 21,41 18,40 6,25 8,26 21))")

	isEqualEx, _ := geom1.GEOSEqualsExact(geom2, 0.01)

	if !isEqualEx {
		t.Errorf("Error: Equals")
	}
}

func TestLwGeomCovers(t *testing.T) {
	geom1, _ := CreateFromWKT("POLYGON((20 39,39 39,39 23,20 23,20 39))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((30 32,36 32,36 27,30 27,30 32))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((41 39,45 39,45 37,41 37,41 39))")
	defer geom3.Free()

	covers, _ := geom1.Covers(geom2)

	if covers != true {
		t.Errorf("Error: Covers()")
	}

	covers, _ = geom1.Covers(geom3)

	if covers != false {
		t.Errorf("Error: Covers()")
	}
}

func TestLwGeomCoveredBy(t *testing.T) {
	geom1, _ := CreateFromWKT("POLYGON((20 39,39 39,39 23,20 23,20 39))")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("POLYGON((30 32,36 32,36 27,30 27,30 32))")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("POLYGON((41 39,45 39,45 37,41 37,41 39))")
	defer geom3.Free()

	covered, _ := geom2.CoveredBy(geom1)

	if covered == false {
		t.Errorf("Error: CoveredBy()")
	}

	covered, _ = geom3.CoveredBy(geom1)

	if covered == true {
		t.Errorf("Error: CoveredBy()")
	}
}

func TestLwGeomCrosses(t *testing.T) {
	geom1, _ := CreateFromWKT("LINESTRING(76 25,79 21,77 19,78 16,77 10)")
	defer geom1.Free()

	geom2, _ := CreateFromWKT("LINESTRING(73 19,75 17,79 18,81 17)")
	defer geom2.Free()

	geom3, _ := CreateFromWKT("LINESTRING(64 28,73 28)")
	defer geom3.Free()

	crosses, _ := geom1.Crosses(geom2)

	if !crosses {
		t.Errorf("Error: Crosses()")
	}

	crosses, _ = geom1.Crosses(geom3)

	if crosses {
		t.Errorf("Error: Crosses()")
	}
}

func TestLwGeomArea(t *testing.T) {
	geom, _ := CreateFromWKT("POLYGON((0 0,0 1,1 1,1 0,0 0))")
	defer geom.Free()

	area, _ := geom.Area()

	if area != 1 {
		t.Errorf("Error: Area()")

	}
}

func TestLwGeomLength(t *testing.T) {
	geom, _ := CreateFromWKT("LINESTRING(0 0, 0 1)")
	defer geom.Free()

	len, _ := geom.Length()

	if len != 1 {
		t.Errorf("Error: Length()")

	}
}
