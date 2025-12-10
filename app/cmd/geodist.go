package cmd

import (
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/geo"
	"github.com/codecrafters-io/redis-starter-go/app/types"
)

func GeoDist(w io.Writer, args []string) {
	if len(args) < 4 {
		w.Write([]byte("-ERR wrong number of arguments for 'GEODIST' command\r\n"))
		return
	}

	key := args[1]
	member1 := args[2]
	member2 := args[3]

	types.SortedSetsMU.Lock()
	sortedSet, exists := types.SortedSets[key]
	types.SortedSetsMU.Unlock()

	if !exists {
		w.Write([]byte("$-1\r\n"))
		return
	}

	data1, exists1 := sortedSet.Data[member1]
	data2, exists2 := sortedSet.Data[member2]

	if !exists1 || !exists2 {
		w.Write([]byte("$-1\r\n"))
		return
	}

	coord1 := geo.Decode(uint64(data1.Score))
	coord2 := geo.Decode(uint64(data2.Score))

	distance := geo.Haversine(coord1.Latitude, coord1.Longitude, coord2.Latitude, coord2.Longitude)

	distanceStr := strconv.FormatFloat(distance, 'f', -1, 64)
	response := "$" + strconv.Itoa(len(distanceStr)) + "\r\n" + distanceStr + "\r\n"
	w.Write([]byte(response))
}
