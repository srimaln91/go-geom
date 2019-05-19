package geom

import (
	"flag"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

var WKTLinestring string
var JSONLinestring string

func TestMain(m *testing.M) {

	lineStringBytes, err := ioutil.ReadFile("../testdata/linestring.txt")

	if err != nil {
		panic(err)
	}

	WKTLinestring = string(lineStringBytes)

	lineStringBytes, err = ioutil.ReadFile("../testdata/lwgeom_loadroute.json")

	if err != nil {
		panic(err)
	}

	JSONLinestring = string(lineStringBytes)

	flag.Parse()
	exitCode := m.Run()

	os.Exit(exitCode)
}

func GetFileContents(path string) string {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(fileBytes)
}

// Clear a set of Geoms from memory
func cleanup(geoms ...*GEOSGeom) {
	for _, geom := range geoms {
		geom.Destroy()
	}
}

func TestFromWKT(t *testing.T) {

	geom := FromWKT(WKTLinestring)

	if geom.cGeom == nil {
		t.Errorf("Error: CreateFromWKT error")
	}
}

func TestBuffer(t *testing.T) {

	geom := FromWKT("POINT (0 0)")
	geom.Buffer(1)

	if geom.cGeom == nil {
		t.Errorf("Error: Buffer() error")
	}

}

func TestSimplifiedBuffer(t *testing.T) {

	geom := FromWKT(WKTLinestring)
	geom.SimplifiedBuffer(0.0001, 0.01745675)
	t.Log(geom.ToWKT())
	if geom.cGeom == nil {
		t.Errorf("Error: Buffer() error")
	}
}

func TestBufferWithStyles(t *testing.T) {

	geom := CreatePoint(0.0, 0.0)

	initialNP, _ := geom.GetNumCoordinates()

	width := 1.0
	quadSegs := 8
	endCapStyle := GeosbufCapRound
	joinStyle := GeosbufJoinMitre
	mitreLimit := 4.0

	geom.BufferWithStyle(width, quadSegs, endCapStyle, joinStyle, mitreLimit)

	bufferNP, _ := geom.GetNumCoordinates()

	if bufferNP == initialNP {
		t.Errorf("Error: BufferWithStyles() error")
	}

	geom.Destroy()
}

func BenchmarkBuffer(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			geom := FromWKT(WKTLinestring)
			geom.SetSRID(4326)
			geom.Buffer(0.01745675)
		}

	})
}

func TestSimplifiedBufferFromWkt(t *testing.T) {
	result := SimplifiedBufferFromWkt(WKTLinestring, 0.0001, 0.01745675)

	if result == "" {
		t.Errorf("Error: SimplifiedBufferFromWkt() error")
	}
}

func BenchmarkBufferWithStyle(b *testing.B) {
	b.ReportAllocs()

	width := 0.01745675
	quadSegs := 8
	endCapStyle := GeosbufCapRound
	joinStyle := GeosbufJoinRound
	mitreLimit := 5.0

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			geom := FromWKT(WKTLinestring)
			geom.SetSRID(4326)
			geom.BufferWithStyle(width, quadSegs, endCapStyle, joinStyle, mitreLimit)
		}

	})
}
func TestToWKT(t *testing.T) {

	wkt := "POINT (0 0)"
	expectedWKT := "POINT (0.0000000000000000 0.0000000000000000)"

	geom := FromWKT(wkt)

	if geom.ToWKT() != expectedWKT {
		t.Errorf("Error: ToWKT(%s) error", wkt)
	}
}

func TestSRID(t *testing.T) {
	geom := FromWKT(WKTLinestring)
	geom.SetSRID(4326)
	srid := geom.GetSRID()

	if srid != 4326 {
		t.Errorf("Error: SRID(%s) error", WKTLinestring)
	}
}

func TestSimplify(t *testing.T) {

	geom := FromWKT(WKTLinestring)

	geom.Simplify(0.01)

	resultWKT := geom.ToWKT()
	expectedWKT := "LINESTRING (79.8561779999999999 6.9118529999999998, 79.8716299999999961 6.8786810000000003, 79.9299139999999966 6.8451360000000001, 79.9798380000000009 6.8385800000000003)"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Simplify(%s) error", resultWKT)
	}
}

func TestSimplifyPreserveTopology(t *testing.T) {

	geom := FromWKT(WKTLinestring)

	geom.SimplifyPreserveTopology(0.01)

	resultWKT := geom.ToWKT()
	expectedWKT := "LINESTRING (79.8561779999999999 6.9118529999999998, 79.8716299999999961 6.8786810000000003, 79.9299139999999966 6.8451360000000001, 79.9798380000000009 6.8385800000000003)"

	if resultWKT != expectedWKT {
		t.Errorf("Error: SimplifyPreserveTopology(%s) error", resultWKT)
	}
}

// func TestReverse(t *testing.T) {
// 	wkt := "LINESTRING (30 10, 10 30, 40 40)"
// 	geom := FromWKT(wkt)

// 	geom.Reverse()

// 	resultWKT := geom.ToWKT()

// 	if resultWKT != "LINESTRING (40.0000000000000000 40.0000000000000000, 10.0000000000000000 30.0000000000000000, 30.0000000000000000 10.0000000000000000)" {
// 		t.Errorf("Error: ToWKT(%s) error", resultWKT)
// 	}
// }

func TestVersion(t *testing.T) {
	version := Version()

	matched, err := regexp.MatchString(`^3\.\d+\.\d+-CAPI-\d+\.\d+\.\d+.+$`, version)

	if err != nil {
		t.Errorf("Error: Version fetch error, %s", err)
	}

	if !matched {
		t.Errorf("Error: Version %s is invalid", version)
	}
}

func TestUnion(t *testing.T) {

	geom1 := FromWKT("LINESTRING(79.856178 6.911853, 79.856382 6.911449, 79.856645 6.910948)")
	geom2 := FromWKT("LINESTRING(79.856178 6.911853, 79.856382 6.911449, 79.856645 6.910948)")

	union := geom1.Union(geom2)

	resultWKT := union.ToWKT()
	expectedWKT := "MULTILINESTRING ((79.8561779999999999 6.9118529999999998, 79.8563819999999964 6.9114490000000002), (79.8563819999999964 6.9114490000000002, 79.8566450000000003 6.9109480000000003))"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Invalid Buffer %s", resultWKT)
	}

	//Cleanup
	cleanup(geom1, geom2, union)
}

func TestIntersection(t *testing.T) {

	geom1 := FromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	geom2 := FromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")

	intersection := geom1.Intersection(geom2)

	resultWKT := intersection.ToWKT()
	expectedWKT := "POLYGON ((79.8560614850391488 6.9114947536352744, 79.8561157638370389 6.9116616436370890, 79.8562337699231648 6.9115058819179929, 79.8561836990276532 6.9113681394729278, 79.8560968252871675 6.9113237218799171, 79.8560614850391488 6.9114947536352744))"

	if resultWKT != expectedWKT {
		t.Errorf("Error: Invalid Intersection %s", resultWKT)
	}

	//Cleanup
	cleanup(geom1, geom2, intersection)
}

func TestIntersects(t *testing.T) {
	geom1 := FromWKT("POLYGON((79.856178 6.911853,79.85598771527475 6.911267935124347,79.85636717179295 6.911461947090437,79.856178 6.911853))")
	geom2 := FromWKT("POLYGON((79.85599371221633 6.911822745245366,79.85623376992316 6.911505881917993,79.85612245824905 6.911199669256946,79.85599371221633 6.911822745245366))")

	intersects, _ := geom1.Intersects(geom2)

	if intersects == false {
		t.Errorf("Error: Intersects")
	}

	//Cleanup
	cleanup(geom1, geom2)

}

func TestDisjoint(t *testing.T) {

	geom1 := FromWKT("POINT(0 0)")
	geom2 := FromWKT("LINESTRING ( 2 0, 0 2 )")
	geom3 := FromWKT("LINESTRING ( 0 0, 0 2 )")

	disjoint, _ := geom1.Disjoints(geom2)

	if disjoint == false {
		t.Errorf("Error: Disjoint")
	}

	disjoint, _ = geom1.Disjoints(geom3)

	if disjoint == true {
		t.Errorf("Error: Disjoint")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}

func TestTouches(t *testing.T) {
	geom1 := FromWKT("POINT(1 1)")
	geom2 := FromWKT("POINT(0 2)")
	geom3 := FromWKT("LINESTRING(0 0, 1 1, 0 2)")
	geom4 := FromWKT("LINESTRING(0 0, 1 1, 0 2)")

	touches, _ := geom1.Touches(geom3)

	if touches == true {
		t.Errorf("Error: Touches")
	}

	touches, _ = geom2.Touches(geom4)

	if touches == false {
		t.Errorf("Error: Touches")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3, geom4)
}

func TestWithin(t *testing.T) {

	geom1 := FromWKT("POLYGON((20 15,23 11,26 16,20 15))")
	geom2 := FromWKT("POLYGON((8 52,21 50,11 46,8 52))")
	geom3 := FromWKT("POLYGON((10 10,25 7,40 10,33 19,24 26,15 20,10 10))")

	disjoint, _ := geom1.Within(geom3)

	if disjoint == false {
		t.Errorf("Error: Within")
	}

	disjoint, _ = geom2.Within(geom3)

	if disjoint == true {
		t.Errorf("Error: Within")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}

func TestContains(t *testing.T) {

	geom1 := FromWKT("POLYGON((20 15,23 11,26 16,20 15))")
	geom2 := FromWKT("POLYGON((8 52,21 50,11 46,8 52))")
	geom3 := FromWKT("POLYGON((10 10,25 7,40 10,33 19,24 26,15 20,10 10))")

	contains, _ := geom3.Contains(geom1)

	if contains == false {
		t.Errorf("Error: Contains")
	}

	contains, _ = geom3.Contains(geom2)

	if contains == true {
		t.Errorf("Error: Contains")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}

func TestOverlaps(t *testing.T) {
	geom1 := FromWKT("POLYGON((17 19,28 19,28 8,17 8,17 19))")
	geom2 := FromWKT("POLYGON((24 22,41 22,41 8,24 8,24 22))")
	geom3 := FromWKT("POLYGON((57 31,61 31,61 26,57 26,57 31))")

	isOVerlaps, _ := geom1.Overlaps(geom2)

	if isOVerlaps == false {
		t.Errorf("Error: Overlaps")
	}

	isOVerlaps, _ = geom3.Overlaps(geom1)

	if isOVerlaps == true {
		t.Errorf("Error: Overlaps")
	}

	//Cleanup
	cleanup(geom1, geom2, geom3)
}

func TestEquals(t *testing.T) {
	geom1 := FromWKT("POLYGON((27 19,38 19,38 9,27 9,27 19))")
	geom2 := FromWKT("POLYGON((27 19,38 19,38 9,27 9,27 19))")
	geom3 := FromWKT("POLYGON((26 21,41 18,40 6,25 8,26 21))")

	isequal, _ := geom1.Equals(geom2)

	if !isequal {
		t.Errorf("Error: Equals")
	}

	isequal, _ = geom1.Equals(geom3)

	if isequal {
		t.Errorf("Error: Equals")
	}

	cleanup(geom1, geom2, geom3)
}

func TestEqualsExact(t *testing.T) {
	geom1 := FromWKT("LINESTRING(0 0, 1 1, 0 2)")
	geom2 := FromWKT("LINESTRING(0 0, 1 1.001, 0 2)")
	// geom3 := FromWKT("POLYGON((26 21,41 18,40 6,25 8,26 21))")

	isEqualEx, _ := geom1.EqualsExact(geom2, 0.01)

	if !isEqualEx {
		t.Errorf("Error: Equals")
	}
}

func TestCreatePoint(t *testing.T) {
	geom1 := CreatePoint(10, 30)
	geom2 := FromWKT("POINT(10 30)")

	if geom1.ToWKT() != geom2.ToWKT() {
		t.Errorf("Error: CreatePoint")
	}

	cleanup(geom1, geom2)
}

func TestCovers(t *testing.T) {
	geom1 := FromWKT("POLYGON((20 39,39 39,39 23,20 23,20 39))")
	geom2 := FromWKT("POLYGON((30 32,36 32,36 27,30 27,30 32))")
	geom3 := FromWKT("POLYGON((41 39,45 39,45 37,41 37,41 39))")

	covers, _ := geom1.Covers(geom2)

	if covers != true {
		t.Errorf("Error: Covers()")
	}

	covers, _ = geom1.Covers(geom3)

	if covers != false {
		t.Errorf("Error: Covers()")
	}

	cleanup(geom1, geom2, geom3)
}

func TestCoveredBy(t *testing.T) {
	geom1 := FromWKT("POLYGON((20 39,39 39,39 23,20 23,20 39))")
	geom2 := FromWKT("POLYGON((30 32,36 32,36 27,30 27,30 32))")
	geom3 := FromWKT("POLYGON((41 39,45 39,45 37,41 37,41 39))")

	covered, _ := geom2.CoveredBy(geom1)

	if covered != true {
		t.Errorf("Error: CoveredBy()")
	}

	covered, _ = geom3.CoveredBy(geom1)

	if covered != false {
		t.Errorf("Error: CoveredBy()")
	}

	cleanup(geom1, geom2, geom3)

}

func TestCrosses(t *testing.T) {
	geom1 := FromWKT("LINESTRING(76 25,79 21,77 19,78 16,77 10)")
	geom2 := FromWKT("LINESTRING(73 19,75 17,79 18,81 17)")
	geom3 := FromWKT("LINESTRING(64 28,73 28)")

	crosses, _ := geom1.Crosses(geom2)

	if !crosses {
		t.Errorf("Error: Crosses()")
	}

	crosses, _ = geom1.Crosses(geom3)

	if crosses {
		t.Errorf("Error: Crosses()")
	}

	cleanup(geom1, geom2, geom3)
}

func TestGetNumCoordinates(t *testing.T) {
	geom := FromWKT("LINESTRING(76 25,79 21,77 19,78 16,77 10)")

	numberCoords, _ := geom.GetNumCoordinates()

	if numberCoords != 5 {
		t.Errorf("Error: GetNumCoordinates()")
	}

	cleanup(geom)
}

func TestArea(t *testing.T) {
	geom := FromWKT("POLYGON((0 0,0 1,1 1,1 0,0 0))")

	area, _ := geom.Area()

	if area != 1 {
		t.Errorf("Error: Area()")

	}
}

func TestLength(t *testing.T) {
	geom := FromWKT("LINESTRING(0 0, 0 1)")

	len, _ := geom.Length()

	if len != 1 {
		t.Errorf("Error: Length()")

	}
}

func TestDistance(t *testing.T) {
	geom1 := FromWKT("LINESTRING(1 0, 1 1)")
	geom2 := FromWKT("LINESTRING(0 0, 0 1)")

	dist, _ := geom1.Distance(geom2)

	if dist != 1 {
		t.Errorf("Error: Distance()")
	}
}

func TestNumPoints(t *testing.T) {
	geom1 := FromWKT("LINESTRING(1 0, 1 1, 25 33)")
	geom2 := FromWKT("POLYGON((20 39,39 39,39 23,20 23,20 39))")

	numPoints, _ := geom1.NumPoints()

	if numPoints != 3 {
		t.Errorf("Error: NumPoints()")
	}

	numPoints, err := geom2.NumPoints()

	if err == nil {
		t.Errorf("Error: NumPoints()")
	}
}
