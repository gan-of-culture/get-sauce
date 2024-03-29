package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

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

func printCaption(i int, caption *static.Caption) {
	fmt.Printf("\n     [%d]  -------------------", i)
	fmt.Printf("\n     Language:            %s\n", caption.Language)
	fmt.Printf("     # download with: ")
	fmt.Printf("get-sauce -c %d ...\n\n", i)
}

func printStream(key string, stream *static.Stream) {
	fmt.Printf("\n     [%s]  -------------------", key)
	if stream.Type == "" {
		stream.Type = static.DataTypeUnknown
	}
	if stream.Type != "" {
		fmt.Printf("\n     Type:            %s", stream.Type)
	}
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

	// for HLS streams the size is only approximated
	sizeFString := "%s\n"
	if stream.Ext != "" && stream.Size > 0 {
		sizeFString = "~ " + sizeFString
	}
	fmt.Printf("\n     Size:            ")
	fmt.Printf(sizeFString, utils.ByteCountSI(stream.Size))
	fmt.Printf("     # download with: ")
	fmt.Printf("get-sauce -s %s ...\n\n", key)
}

func printInfo(data *static.Data) {
	printHeader(data)

	if len(data.Captions) > 0 {
		fmt.Println("\n Captions:  # All available languages")
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

// mergeMediaFiles into one output file using ffmpeg | merges video + audio + subtitles
func mergeMediaFiles(files []string, outFile string) error {
	if len(files) < 2 {
		return nil
	}
	if !config.Quiet {
		fmt.Println("\nMerging files using ffmpeg...")
	}

	if strings.HasSuffix(outFile, ".webm") {
		outFile = strings.ReplaceAll(outFile, ".webm", ".mp4")
	}

	command := []string{"-y"}
	var caption string
	for _, f := range files {
		p, _ := filepath.Abs(f)
		switch utils.GetFileExt(f) {
		case "ass", "srt", "vtt":
			caption = p
		}
		command = append(command, "-i")
		command = append(command, p)
	}
	command = append(command, "-c")
	command = append(command, "copy")
	if caption != "" {
		command = append(command, "-c:s")
		command = append(command, "mov_text")
	}
	command = append(command, outFile)

	if !config.Quiet {
		fmt.Println(command)
	}
	cmd := exec.Command("ffmpeg", command...)
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}
	if !config.Quiet {
		fmt.Println("Success!")
	}

	return nil
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
