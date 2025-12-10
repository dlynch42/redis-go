package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/geo"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func GeoAdd(w io.Writer, args []string) {
	if len(args) < 5 || (len(args)-2)%3 != 0 {
		w.Write([]byte("-ERR wrong number of arguments for 'GEOADD' command\r\n"))
		return
	}

	key := args[1]
	lonStr := args[2]
	latStr := args[3]
	member := args[4]

	types.SortedSetsMU.Lock()
	sortedSet, exists := types.SortedSets[key]
	if !exists {
		sortedSet = types.SortedSet{
			Data: make(map[string]types.SortedSetEntry),
		}
	}

	addedCount := 0
	lon, _ := strconv.ParseFloat(lonStr, 64)
	lat, _ := strconv.ParseFloat(latStr, 64)

	if !validateLongitude(lon) && !validateLatitude(lat) {
		types.SortedSetsMU.Unlock()
		w.Write([]byte(fmt.Sprintf("-ERR invalid longitude,latitude pair %s,%s\r\n", lonStr, latStr)))
		return
	}

	if !validateLatitude(lat) {
		types.SortedSetsMU.Unlock()
		w.Write([]byte(fmt.Sprintf("-ERR invalid latitude %s\r\n", latStr)))
		return
	}

	if !validateLongitude(lon) {
		types.SortedSetsMU.Unlock()
		w.Write([]byte(fmt.Sprintf("-ERR invalid longitude %s\r\n", lonStr)))
		return
	}

	score := float64(geo.Encode(lat, lon))

	_, memberExists := sortedSet.Data[member]
	sortedSet.Data[member] = types.SortedSetEntry{
		Score: score,
		Value: member,
	}
	if !memberExists {
		addedCount++
	}

	types.SortedSets[key] = sortedSet
	types.SortedSetsMU.Unlock()

	response := fmt.Sprintf(":%d\r\n", addedCount)
	w.Write([]byte(response))
}

func validateLongitude(lon float64) bool {
	// Longitude: -180 to +180 (inclusive)
	return lon >= -180 && lon <= 180
}

func validateLatitude(lat float64) bool {
	// Latitude: -85.05112878 to +85.05112878 (inclusive)
	return lat >= -85.05112878 && lat <= 85.05112878
}
