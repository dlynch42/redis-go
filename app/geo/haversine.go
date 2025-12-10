package geo

import "math"

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	return hsDist(degPos(lat1, lon1), degPos(lat2, lon2))
}

func degPos(lat float64, lon float64) Position {
	return Position{
		Latitude:  lat * math.Pi / 180,
		Longitude: lon * math.Pi / 180,
	}
}

func hv(theta float64) float64 {
	return 0.5 * (1 - math.Cos(theta))
}

func hsDist(p1, p2 Position) float64 {
	return 2 * radius * math.Asin(math.Sqrt(hv(p2.Latitude-p1.Latitude)+math.Cos(p1.Latitude)*math.Cos(p2.Latitude)*hv(p2.Longitude-p1.Longitude)))
}
