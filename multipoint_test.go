package geojson

import "testing"

func TestMultiPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3]]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"MultiPoint","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2]],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"MultiPoint"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPoint","coordinates":null}`, errCoordinatesInvalid)
}
