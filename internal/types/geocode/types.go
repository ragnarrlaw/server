package geocode

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (p *Point) UnmarshalJSON(b []byte) error {
	return nil
}

type GeoFeatureType string

const (
	FeatureType           GeoFeatureType = "Feature"
	FeatureCollectionType GeoFeatureType = "FeatureCollection"
)

type GeometryType string

const (
	PointType          GeometryType = "Point"
	MultiPointType     GeometryType = "MultiPoint"
	LineStringType     GeometryType = "LineString"
	MultiLineString    GeometryType = "MultiLineString"
	PolygonType        GeometryType = "Polygon"
	MultiPolygonType   GeometryType = "MultiPolygon"
	GeometryCollection GeometryType = "GeometryCollection"
)

type Geometry struct {
	Type            GeometryType `json:"type"`
	BoundingBox     []float64    `json:"bbox,omitempty"` // bounding box -> "bbox": [minLon, minLat, maxLon, maxLat]
	Point           []float64    // -> [lon(x). lat(y)]
	MultiPoint      [][]float64  // -> [[lon(x). lat(y)], [lon(x). lat(y)]]
	LineString      [][]float64
	MultiLineString [][][]float64
	Polygon         [][][]float64
	MultiPolygon    [][][][]float64
	Geometries      []*Geometry
	CRS             map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects -> leaflet EPSG:3857
}

/*
Example marshalled geojson geometry

	{
		"type": "Point",
		"coordinates": [100.0, 0.0],
	}
*/
func (g Geometry) MarshalJSON() ([]byte, error) {
	var p struct {
		Type        GeometryType           `json:"type"`
		BoundingBox []float64              `json:"bbox,omitempty"`
		Coordinates interface{}            `json:"coordinates"`
		CRS         map[string]interface{} `json:"crs,omitempty"`
	}
	p.Type = g.Type
	if len(g.BoundingBox) > 0 {
		p.BoundingBox = g.BoundingBox
	}

	if len(g.CRS) > 0 {
		p.CRS = g.CRS
	}

	switch g.Type {
	case PointType:
		p.Coordinates = g.Point
	case MultiPointType:
		p.Coordinates = g.MultiPoint
	case LineStringType:
		p.Coordinates = g.LineString
	case MultiLineString:
		p.Coordinates = g.MultiLineString
	case PolygonType:
		p.Coordinates = g.Polygon
	case MultiPolygonType:
		p.Coordinates = g.MultiPolygon
	case GeometryCollection:
		p.Coordinates = g.Geometries
	}
	return json.Marshal(p)
}

func (g *Geometry) UnmarshalJSON(b []byte) error {
	geom := make(map[string]interface{})
	if err := json.Unmarshal(b, &geom); err != nil {
		return err
	}
	return unmarshalGeometry(g, geom)
}

func unmarshalPoint(point interface{}) ([]float64, error) {
	p, ok := point.([]interface{})
	if !ok {
		return nil, fmt.Errorf("point is not valid")
	}
	cs := make([]float64, 0, len(p))
	for _, c := range p {
		f, ok := c.(float64)
		if !ok {
			return nil, fmt.Errorf("%v is not a valid coordinate", c)
		}
		cs = append(cs, f)
	}
	return cs, nil
}

func unmarshalMultiPoint(points interface{}) ([][]float64, error) {
	ps, ok := points.([]interface{})
	if !ok {
		return nil, fmt.Errorf("points are not valid")
	}
	cs := make([][]float64, 0, len(ps))
	for _, p := range ps {
		fs, err := unmarshalPoint(p)
		if err != nil {
			return nil, err
		}
		cs = append(cs, fs)
	}
	return cs, nil
}

func unmarshalLineString(line interface{}) ([][]float64, error) {
	l, ok := line.([]interface{})
	if !ok {
		return nil, fmt.Errorf("line is invalid")
	}
	return unmarshalMultiPoint(l)
}

func unmarshalMultiLineString(lines interface{}) ([][][]float64, error) {
	ls, ok := lines.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid multiple lines")
	}
	mls := make([][][]float64, 0, len(ls))
	for _, sl := range ls {
		l, err := unmarshalLineString(sl)
		if err != nil {
			return nil, err
		}
		mls = append(mls, l)
	}
	return mls, nil
}

func unmarshalPolygon(polygon interface{}) ([][][]float64, error) {
	ps, ok := polygon.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid polygon type")
	}
	polygons := make([][][]float64, 0, len(ps))
	for _, plygon := range ps {
		ply, ok := plygon.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid polygon type")
		}
		p, err := unmarshalLineString(ply)
		if err != nil {
			return nil, err
		}
		polygons = append(polygons, p)
	}
	return polygons, nil
}

func unmarshalMultiPolygon(polygons interface{}) ([][][][]float64, error) {
	mps, ok := polygons.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid multi-polygon type")
	}
	ps := make([][][][]float64, 0, len(mps))
	for _, mp := range mps {
		ply, ok := mp.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid polygon type")
		}
		p, err := unmarshalPolygon(ply)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func unmarshalGeometry(geometry *Geometry, geom interface{}) error {
	g, ok := geom.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid geometry type")
	}
	t, exists := g["type"]
	if !exists {
		return fmt.Errorf("type field is not defined on the geometry")
	}
	str, ok := t.(string)
	if !ok {
		return fmt.Errorf("geometry type field is not of string type")
	}
	geometry.Type = GeometryType(str)

	var err error
	switch geometry.Type {
	case PointType:
		geometry.Point, err = unmarshalPoint(g["coordinates"])
		if err != nil {
			return err
		}
	case MultiPointType:
		geometry.MultiPoint, err = unmarshalMultiPoint(g["coordinates"])
		if err != nil {
			return err
		}
	case LineStringType:
		geometry.LineString, err = unmarshalLineString(g["coordinates"])
		if err != nil {
			return err
		}
	case MultiLineString:
		geometry.MultiLineString, err = unmarshalMultiLineString(g["coordinates"])
		if err != nil {
			return err
		}
	case PolygonType:
		geometry.Polygon, err = unmarshalPolygon(g["coordinates"])
		if err != nil {
			return err
		}
	case MultiPolygonType:
		geometry.MultiPolygon, err = unmarshalMultiPolygon(g["coordinates"])
		if err != nil {
			return err
		}
	case GeometryCollection:
		geometries, exists := g["geometries"]
		if !exists {
			return fmt.Errorf("geometries field is not defined on the geometry")
		}
		geometry.Geometries, err = unmarshalGeometries(geometries)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported geometry type")
	}
	return nil
}

func unmarshalGeometries(geoms interface{}) ([]*Geometry, error) {
	gms, ok := geoms.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid geometries types")
	}
	gs := make([]*Geometry, 0, len(gms))
	for _, g := range gms {
		geom := &Geometry{}
		gm, ok := g.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid geometric type %v", g)
		}
		err := unmarshalGeometry(geom, gm)
		if err != nil {
			return nil, err
		}
		gs = append(gs, geom)
	}
	return gs, nil
}

type GeoJsonFeature struct {
	Type        GeoFeatureType         `json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"` // "bbox": [minLon, minLat, maxLon, maxLat]
	Geometry    *Geometry              `json:"geometry"`
	Properties  map[string]interface{} `json:"properties"`
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects -> leaflet EPSG:3857
}

type GeoJsonFeatureCollection struct {
	Type        GeoFeatureType         `json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"` // "bbox": [minLon, minLat, maxLon, maxLat]
	Features    []*GeoJsonFeature      `json:"features"`
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects -> leaflet EPSG:3857
}
