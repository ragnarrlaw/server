package geolocation

/**
Add the geolocation operations related utils here
*/

import "math"

// Haversine formula, calculate distance between two points
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers
	latDelta := (lat2 - lat1) * (math.Pi / 180)
	lonDelta := (lon2 - lon1) * (math.Pi / 180)

	a := math.Sin(latDelta/2)*math.Sin(latDelta/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*math.Sin(lonDelta/2)*math.Sin(lonDelta/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c // KM distance
}
