package downloader

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reSanitizeTitle = regexp.MustCompile(`["&|:?<>/*\\ ]+`)

// GenSortedStreams for stream map
func GenSortedStreams(streams map[string]*static.Stream) []*static.Stream {
	index := make([]int64, 0, len(streams))
	for k := range streams {
		idx, _ := strconv.ParseInt(k, 10, 0)
		index = append(index, idx)
	}
	if len(index) < 1 {
		return nil
	}
	slices.Sort(index)

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

func printCaption(i int, caption *static.Caption) {
	fmt.Printf("\n     [%d]  -------------------", i)
	fmt.Printf("\n     Language:            %s", caption.Language)
	fmt.Printf("\n     # download with: ")
	fmt.Printf("get-sauce -c %d ...\n\n", i)
}

func printStream(key string, stream *static.Stream) {
	fmt.Printf("\n     [%s]  -------------------", key)
	if stream.Type == "" {
		stream.Type = static.DataTypeUnknown
	}
	fmt.Printf("\n     Type:            %s", stream.Type)

	if stream.Info != "" {
		fmt.Printf("\n     Info:            %s", stream.Info)
	}
	if stream.Quality != "" {
		fmt.Printf("\n     Quality:         %s", stream.Quality)
	}
	if len(stream.URLs) > 1 {
		fmt.Printf("\n     Parts:           %d", len(stream.URLs))
	}

	if stream.Size > 0 {
		// for HLS streams the size is only approximated
		sizeFString := "%s"
		if stream.Ext != "" {
			sizeFString = "~ " + sizeFString
		}
		fmt.Printf("\n     Size:            ")
		fmt.Printf(sizeFString, utils.ByteCountSI(stream.Size))
	}
	fmt.Printf("\n     # download with: ")
	fmt.Printf("get-sauce -s %s ...\n\n", key)
}

func printInfo(data *static.Data) {
	printHeader(data)

	if len(data.Captions) > 0 {
		fmt.Print("\n Captions:  # All available languages")
	}
	for i, caption := range data.Captions {
		printCaption(i, caption)
	}

	sortedStreams := GenSortedStreams(data.Streams)
	fmt.Print("\n Streams:   # All available qualities")
	for i, stream := range sortedStreams {
		printStream(fmt.Sprint(i), stream)
	}
}

func printStreamInfo(data *static.Data, streamKey string) {
	printHeader(data)

	if len(data.Captions) > config.Caption && config.Caption > -1 {
		fmt.Println("\n Caption:   ")
		printCaption(config.Caption, data.Captions[config.Caption])
	}

	fmt.Println("\n Stream:   ")
	printStream(streamKey, data.Streams[streamKey])
}

func sanitizeVTT(fileURI string) error {
	// sometimes VTT contains weird blank lines that will cause an issue if you try to merge it later with ffmpeg
	// this routine removes said lines
	// it also contains text separated from other text by blank lines this also causes issues with ffmpeg later
	fileContent, err := os.ReadFile(fileURI)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`((?:\d{2,}:)?\d\d:\d\d\.\d{3} --> (?:\d{2,}:)?\d\d:\d\d\.\d{3})\s*((?:\n[^\n]+)+)`) // 1=TimeStamp 2=Text
	out := "WEBVTT"
	for _, match := range re.FindAllStringSubmatch(string(fileContent), -1) {
		out = fmt.Sprintf("%s\n\n%s\n%s", out, match[1], strings.TrimSpace(match[2]))
	}

	return os.WriteFile(fileURI, []byte(out), 0644)
}

func sanitizeTitle(title string) string {
	title = reSanitizeTitle.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)
	title = strings.TrimRight(title, ".")
	title = strings.ReplaceAll(title, "  ", " ")
	return title
}
