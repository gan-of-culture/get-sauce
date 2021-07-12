package downloader

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func genSortedStreams(streams map[string]*static.Stream) []*static.Stream {
	index := make([]int64, 0, len(streams))
	for k := range streams {
		idx, _ := strconv.ParseInt(k, 10, 0)
		index = append(index, idx)
	}
	if len(index) < 1 {
		return nil
	}
	sort.Slice(index, func(i, j int) bool {
		return index[i] < index[j]
	})

	sortedStreams := make([]*static.Stream, 0, len(streams))
	for _, i := range index {
		sortedStreams = append(sortedStreams, streams[fmt.Sprint(i)])
	}
	return sortedStreams
}

func printHeader(data *static.Data) {
	fmt.Printf("\n Site:      %s", data.Site)
	fmt.Printf("\n Title:     %s", data.Title)
	fmt.Printf("\n Type:      %s", data.Type)
}

func printStream(key string, stream *static.Stream) {
	fmt.Printf("\n     [%s]  -------------------", key)
	if stream.Info != "" {
		fmt.Printf("\n     Info:            %s", stream.Info)
	}
	if stream.Quality == "" {
		stream.Quality = "unknown"
	}
	fmt.Printf("\n     Quality:         %s", stream.Quality)
	if len(stream.URLs) > 1 {
		fmt.Printf("\n     Parts:           %d", len(stream.URLs))
	}
	fmt.Printf("\n     Size:            ")
	fmt.Printf("%.2f MB (%d Bytes)\n", float64(stream.Size)/(1_000_000), stream.Size)
	fmt.Printf("     # download with: ")
	fmt.Printf("go-hentai-scraper -s %s ...\n\n", key)
}

func printInfo(data *static.Data) {
	printHeader(data)

	sortedStreams := genSortedStreams(data.Streams)
	fmt.Print("\n Streams:   # All available qualities")
	for i, stream := range sortedStreams {
		printStream(fmt.Sprint(i), stream)
	}
}

func printStreamInfo(data *static.Data, streamKey string) {
	printHeader(data)

	fmt.Println("\n Stream:   ")
	printStream(streamKey, data.Streams[streamKey])
}
