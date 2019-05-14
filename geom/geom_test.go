package geom

import (
	"testing"
)

func TestLwGeomBufferWithParams(t *testing.T) {
	lwgeom := FromGeoJSON(string(JSONLinestring))
	lwgeom.SetSRID(4326)
	defer lwgeom.Free()

	lwgeom.LineSubstring(0.5, 0.9)

	fromSRS := SRS["EPSG:4326"]
	toSRS := SRS["EPSG:3857"]

	lwgeom.Project(fromSRS, toSRS)

	bufferParams, _ := CreateBufferParams()
	defer bufferParams.Destroy()

	bufferParams.SetJoinStyle(GeosbufJoinRound)
	bufferParams.SetEndCapStyle(GeosbufCapFlat)
	bufferParams.SetQuadrantSegments(8)

	lwgeom.BufferWithParams(bufferParams, 200)

	lwgeom.Project(toSRS, fromSRS)

	bufJSON := lwgeom.ToGeoJSON(4, 0)

	if bufJSON == "" {
		t.Error("Error: BufferWithParams()")
	}
}

func TestLwGeomBuffer(t *testing.T) {
	lwgeom := FromGeoJSON(string(JSONLinestring))
	lwgeom.SetSRID(4326)

	defer lwgeom.Free()

	lwgeom.LineSubstring(0.5, 0.9)

	fromSRS := SRS["EPSG:4326"]
	toSRS := SRS["EPSG:3857"]

	lwgeom.Project(fromSRS, toSRS)
	lwgeom.Buffer(200)
	lwgeom.Project(toSRS, fromSRS)

	bufJSON := lwgeom.ToGeoJSON(4, 0)

	if bufJSON == "" {
		t.Error("Error: Buffer()")
	}
}
