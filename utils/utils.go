package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

// GetLastItemString of slice
func GetLastItemString(slice []string) string {
	if len(slice) <= 0 {
		return ""
	}
	return slice[len(slice)-1]
}

// CalcSizeInByte func
func CalcSizeInByte(number float64, unit string) int64 {
	switch unit {
	case "KB":
		return int64(number) * 1000
	case "MB":
		return int64(number) * 1000000
	case "GB":
		return int64(number) * 10000000000
	default:
		return int64(number)
	}
}

// NeedDownloadList return the indices of gallery that need download
func NeedDownloadList(length int) []int {
	if config.Pages != "" {
		var items []int
		var selStart, selEnd int
		temp := strings.Split(config.Pages, ",")

		for _, i := range temp {
			selection := strings.Split(i, "-")
			selStart, _ = strconv.Atoi(strings.TrimSpace(selection[0]))

			if len(selection) >= 2 {
				selEnd, _ = strconv.Atoi(strings.TrimSpace(selection[1]))
			} else {
				selEnd = selStart
			}

			for item := selStart; item <= selEnd; item++ {
				items = append(items, item)
			}
		}
		return items
	}
	out := []int{}
	for i := 1; i <= length; i++ {
		out = append(out, i)
	}
	return out
}

// IsInTests to limit run time of some extractors
func IsInTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.run") {
			return true
		}
	}
	return false
}

// GetMediaType e.x put in png get image/png, mp4 -> video/mp4
func GetMediaType(t string) string {
	switch t {
	case "jpg", "jpeg", "png", "gif", "webp":
		return fmt.Sprintf("%s/%s", "image", t)
	case "webm", "mp4", "mkv", "m4a":
		return fmt.Sprintf("%s/%s", "video", t)
	case "txt", "m3u8":
		return fmt.Sprintf("%s/%s", "application", "x-mpegurl")
	default:
		return fmt.Sprintf("%s/%s", "unknown", t)
	}
}

// GetH1 of html file
func GetH1(htmlString string) string {
	re := regexp.MustCompile(`[^>]*</h1>`)
	h1s := re.FindAllString(htmlString, -1)
	return strings.TrimSuffix(GetLastItemString(h1s), "</h1>")
}

// GetMeta of html file
func GetMeta(htmlString, property string) string {
	re := regexp.MustCompile(fmt.Sprintf("<meta property=[\"']%s[\"'] content=[\"']([^\"']*)", property))
	metaTags := re.FindAllStringSubmatch(htmlString, -1)
	if len(metaTags) < 1 {
		log.Println(htmlString)
		return fmt.Sprintf("No matches found for %s", property)
	}
	return metaTags[0][1]
}
