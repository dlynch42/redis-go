package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/geo"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

var unitMap = map[string]float64{
	"M":  geo.Meters,
	"KM": geo.Kilometers,
	"MI": geo.Miles,
	"FT": geo.Feet,
}

func GeoSearch(w io.Writer, args []string) {
	if len(args) < 8 {
		w.Write([]byte("-ERR wrong number of arguments for 'GEOSEARCH' command\r\n"))
		return
	}

	key := args[1]
	from := args[2]
	long := args[3]
	lat := args[4]
	byrad := args[5]
	radius := args[6]
	unit := args[7]

	position := geo.Position{}
	results := []string{}

	if strings.ToUpper(from) == "FROMLONLAT" {
		pos, err := fromLonLat(long, lat)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
			return
		}
		position = pos
	}

	if strings.ToUpper(byrad) == "BYRADIUS" {
		results = byRadius(key, position, radius, unit)
		response := resp.EncodeRESPArray(results)
		w.Write([]byte(response))
		return
	}

	// Fallback
	w.Write([]byte("*0\r\n"))
}

func fromLonLat(longStr, latStr string) (geo.Position, error) {
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return geo.Position{}, fmt.Errorf("invalid longitude")
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return geo.Position{}, fmt.Errorf("invalid latitude")
	}
	return geo.Position{Longitude: long, Latitude: lat}, nil
}

func byRadius(key string, position geo.Position, radiusStr string, unitStr string) []string {
	rad, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		return []string{}
	}

	// Convert radius to meters
	unit := unitMap[strings.ToUpper(unitStr)]
	radius := rad * unit

	types.SortedSetsMU.Lock()
	sortedSet, exists := types.SortedSets[key]
	types.SortedSetsMU.Unlock()

	if !exists {
		return []string{}
	}

	var results []string
	for member, entry := range sortedSet.Data {
		// Decode geohash to get member coordinates
		coords := geo.Decode(uint64(entry.Score))

		// Calculate distance using Haversine formula
		distance := geo.Haversine(
			position.Latitude, position.Longitude,
			coords.Latitude, coords.Longitude,
		)

		// If within radius, include in results
		if distance <= radius {
			results = append(results, member)
		}
	}

	return results
}
