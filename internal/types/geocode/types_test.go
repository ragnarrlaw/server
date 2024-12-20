package geocode

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	// unmarshal point geometry
	point := []byte(`{"type":"Point","coordinates":[100.0,0.0]}`)
	geometry := &Geometry{}
	if err := geometry.UnmarshalJSON(point); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal line string geometry
	lineString := []byte(`{"type":"LineString","coordinates":[[100.0,0.0],[101.0,1.0]]}`)
	geometry = &Geometry{}
	if err := geometry.UnmarshalJSON(lineString); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal multi point geometry
	multiPoint := []byte(`{"type":"MultiPoint","coordinates":[[100.0,0.0],[101.0,1.0]]}`)
	geometry = &Geometry{}
	if err := json.Unmarshal(multiPoint, geometry); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal polygon geometry
	polygon := []byte(`
	{
    	"type": "Polygon",
    	"coordinates": [
        	[
            	[100.0, 0.0],
            	[101.0, 0.0],
            	[101.0, 1.0],
            	[100.0, 1.0],
            	[100.0, 0.0]
        	]
    	]
	}
	`)
	geometry = &Geometry{}
	if err := json.Unmarshal(polygon, geometry); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal multi polygon geometry
	multiPolygon := []byte(`
	{
		"type": "MultiPolygon",
		"coordinates": [
			[
				[
					[100.0, 0.0],
					[101.0, 0.0],
					[101.0, 1.0],
					[100.0, 1.0],
					[100.0, 0.0]
				]
			],
			[
				[
					[102.0, 2.0],
					[103.0, 2.0],
					[103.0, 3.0],
					[102.0, 3.0],
					[102.0, 2.0]
				]
			]
		]
	}
	`)
	geometry = &Geometry{}
	if err := json.Unmarshal(multiPolygon, geometry); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal geometry
	geometries := []byte(`
	{
		"type": "GeometryCollection",
		"geometries": [
			{
				"type": "Point",
				"coordinates": [100.0, 0.0]
			},
			{
				"type": "LineString",
				"coordinates": [
					[101.0, 0.0],
					[102.0, 1.0]
				]
			},
			{
				"type": "Polygon",
				"coordinates": [
					[
						[102.0, 2.0],
						[103.0, 2.0],
						[103.0, 3.0],
						[102.0, 3.0],
						[102.0, 2.0]
					]
				]
			}
		]
	}
	`)
	geometry = &Geometry{}
	if err := json.Unmarshal(geometries, geometry); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// unmarshal invalid geometry
	invalid := []byte(`{"type":"Invalid","coordinates":[100.0,0.0]}`)
	geometry = &Geometry{}
	if err := json.Unmarshal(invalid, geometry); err == nil {
		t.Error("expected error but got nil")
	}
}
