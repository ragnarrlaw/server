package geocode

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (p *Point) UnmarshalJSON(b []byte) error {
	return nil
}

type GeometryType string

const (
	PointType       GeometryType = "Point"
	MultiPointType  GeometryType = "MultiPoint"
	MultiLineString GeometryType = "MultiLineString"
	MultiPolygon    GeometryType = "MultiPolygon"
	LineStringType  GeometryType = "LineString"
	PolygonType     GeometryType = "Polygon"
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

func (g *Geometry) UnmarshalJSON(b []byte) error {
	panic("function not implemented")
	return nil
}

type GeoJsonFeature struct {
	Type        string                 `json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"` // "bbox": [minLon, minLat, maxLon, maxLat]
	Geometry    *Geometry              `json:"geometry"`
	Properties  map[string]interface{} `json:"properties"`
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects -> leaflet EPSG:3857
}

type GeoJsonFeatureCollection struct {
	Type        string                 `json:"type"`
	BoundingBox []float64              `json:"bbox,omitempty"` // "bbox": [minLon, minLat, maxLon, maxLat]
	Features    []*GeoJsonFeature      `json:"features"`
	CRS         map[string]interface{} `json:"crs,omitempty"` // Coordinate Reference System Objects -> leaflet EPSG:3857
}
