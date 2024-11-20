package types

type SearchPayload struct {
	List     string   `json:"list"`
	Location GeoPoint `json:"location"`
}
