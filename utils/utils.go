package utils

import (
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

// NeedDownloadList return the indices of playlist that need download
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
