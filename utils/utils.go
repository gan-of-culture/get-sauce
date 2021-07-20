package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
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

// GetMediaType e.g. put in png get image, mp4 -> video
func GetMediaType(t string) static.DataType {
	switch t {
	case "jpg", "jpeg", "png", "gif", "webp", "avif":
		return static.DataTypeImage
	case "webm", "mp4", "mkv", "m4a", "txt", "m3u8", "avi":
		return static.DataTypeVideo
	default:
		return static.DataTypeUnknown
	}
}

// GetH1 of html file
// idx -1 = last h1 found
// if index out of range set to last h1
func GetH1(htmlString *string, idx int) string {
	re := regexp.MustCompile(`[^>]*</h1>`)
	h1s := re.FindAllString(*htmlString, -1)

	h1sLen := len(h1s)
	if idx == -1 {
		idx = h1sLen
	}

	// if index out of range set last
	if h1sLen < idx+1 {
		idx = h1sLen - 1
		if idx == -1 {
			return ""
		}
	}
	return strings.TrimSuffix(h1s[idx], "</h1>")
}

// GetMeta of html file
func GetMeta(htmlString *string, property string) string {
	re := regexp.MustCompile(fmt.Sprintf("<meta property=[\"']%s[\"'] content=[\"']([^\"']*)", property))
	metaTags := re.FindAllStringSubmatch(*htmlString, -1)
	if len(metaTags) < 1 {
		return fmt.Sprintf("no matches found for %s", property)
	}
	return metaTags[0][1]
}

// RemoveAdjDuplicates of string slice
func RemoveAdjDuplicates(slice []string) []string {
	out := []string{}
	var last string
	for _, s := range slice {
		if s != last {
			out = append(out, s)
		}
		last = s
	}

	return out
}

// ParseM3UMaster into static.Stream to prefill the structure
// returns a pre filled structure where URLs[0].URL is the media m3u URI
func ParseM3UMaster(master *string) (map[string]*static.Stream, error) {
	re := regexp.MustCompile(`#EXT-X-STREAM-INF:([^\n]*)\n([^\n]+)`) // 1=PARAMS 2=MEDIAURI
	matchedStreams := re.FindAllStringSubmatch(*master, -1)
	if len(matchedStreams) < 1 {
		return nil, fmt.Errorf("unable to parse any stream in m3u master file: %s", *master)
	}

	out := map[string]*static.Stream{}
	for i, stream := range matchedStreams {
		s := &static.Stream{}

		for _, v := range stream[1:] {

			re = regexp.MustCompile(`([A-Z\-]+=(?:"[^"]*"|[^,]*))`) // 1=list of PARAMNAME=value,
			matchedStreamParams := re.FindAllStringSubmatch(v, -1)
			if len(matchedStreamParams) == 0 {
				s.URLs = []*static.URL{
					{
						URL: strings.TrimSpace(v),
						Ext: "",
					},
				}
				continue
			}

			for _, streamParam := range matchedStreamParams[1:] {

				splitParam := strings.Split(streamParam[1], "=")
				splitParam[1] = strings.Trim(splitParam[1], `",`)
				switch splitParam[0] {
				case "BANDWIDTH":
					size, err := strconv.ParseInt(splitParam[1], 10, 64)
					if err != nil {
						return nil, err
					}
					s.Size = size
				case "RESOLUTION":
					s.Quality = splitParam[1]
				case "CODECS":
					s.Info = splitParam[1]
				}
			}
		}
		out[fmt.Sprint(i)] = s
	}

	return out, nil
}

// Wrap error with context
func Wrap(err error, ctx string) error {
	return errors.New(err.Error() + ": " + ctx)
}
