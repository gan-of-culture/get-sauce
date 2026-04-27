package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

type MergeFile struct {
	path     string
	dataType static.DataType
}

// mergeMediaFiles into one output file using ffmpeg | merges video + audio + subtitles
func mergeMediaFiles(files []MergeFile, outFile string) error {
	if len(files) < 2 {
		return nil
	}
	if !config.Quiet {
		fmt.Println("\nMerging files using ffmpeg...")
	}

	command := []string{"-y"}
	for _, f := range files {
		p, _ := filepath.Abs(f.path)
		command = append(command, "-i")
		command = append(command, p)
	}

	// do ffmpeg mapping
	command = append(command, "-map")
	command = append(command, "0:v")
	command = append(command, getAudioMapping(files)...)
	command = append(command, getCaptionMapping(files)...)

	command = append(command, "-c")
	command = append(command, "copy")
	if len(getCaptionMapping(files)) > 0 {
		switch filepath.Ext(outFile) {
		case ".mp4":
			command = append(command, "-c:s")
			command = append(command, "mov_text")
		default:
			// Specifically ensures the subtitle codec is set correctly for the WebM container (in this case).
			// WebM currently has no support for styled subitltes like .ass or .srt. They are getting style-stripped
			// and added as webvtt.
			//
			// For other contaiers this logic might need be adjusted.
			command = append(command, "-c:s")
			command = append(command, "webvtt")
			command = append(command, "-disposition:s:0")
			command = append(command, "default")
		}
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
		err := os.Remove(f.path)
		if err != nil {
			return err
		}
	}
	if !config.Quiet {
		fmt.Println("Success!")
	}

	return nil
}

func getAudioMapping(files []MergeFile) []string {
	var cmd []string
	for idx, a := range files {
		if a.dataType != static.DataTypeAudio {
			continue
		}
		cmd = append(cmd, "-map")
		cmd = append(cmd, fmt.Sprintf("%d:a", idx))
	}
	if len(cmd) == 0 {
		// take audio from video stream instead
		cmd = append(cmd, "-map")
		cmd = append(cmd, "0:a")
	}
	return cmd
}

func getCaptionMapping(files []MergeFile) []string {
	var cmd []string
	for idx, a := range files {
		if a.dataType != static.DataTypeText {
			continue
		}
		cmd = append(cmd, "-map")
		cmd = append(cmd, fmt.Sprintf("%d:s", idx))
	}
	return cmd
}
