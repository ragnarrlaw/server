package types

type GeoPoint struct {
	GeoCode   string  `json:"geoCode"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
