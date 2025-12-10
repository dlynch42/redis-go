package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/geo"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func GeoPos(w io.Writer, args []string) {
	if len(args) < 3 {
		w.Write([]byte("-ERR wrong number of arguments for 'GEOPOS' command\r\n"))
		return
	}

	key := args[1]
	members := args[2:]

	types.SortedSetsMU.Lock()
	sortedSet, keyExists := types.SortedSets[key]
	types.SortedSetsMU.Unlock()

	response := fmt.Sprintf("*%d\r\n", len(members))

	for _, member := range members {
		if !keyExists {
			response += "*-1\r\n"
			continue
		}
		data, exists := sortedSet.Data[member]
		if !exists {
			response += "*-1\r\n"
			continue
		}

		// Decode the geohash to get latitude and longitude
		coordinates := geo.Decode(uint64(data.Score))

		latStr := strconv.FormatFloat(coordinates.Latitude, 'f', -1, 64)
		lonStr := strconv.FormatFloat(coordinates.Longitude, 'f', -1, 64)

		response += "*2\r\n"
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(lonStr), lonStr)
		response += fmt.Sprintf("$%d\r\n%s\r\n", len(latStr), latStr)

	}
	w.Write([]byte(response))
}
