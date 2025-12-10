package geo

const (
	MIN_LATITUDE  = -85.05112878
	MAX_LATITUDE  = 85.05112878
	MIN_LONGITUDE = -180.0
	MAX_LONGITUDE = 180.0

	LATITUDE_RANGE  = MAX_LATITUDE - MIN_LATITUDE
	LONGITUDE_RANGE = MAX_LONGITUDE - MIN_LONGITUDE
)

const radius = 6372797.560856 // Earth radius in kilometers

const (
	Meters     float64 = 1.0
	Kilometers float64 = 1000.0
	Miles      float64 = 1609.344
	Feet       float64 = 0.3048
)

type Position struct {
	Longitude float64
	Latitude  float64
}
