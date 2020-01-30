package ehentai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const ehentai = "https://e-hentai.org/"
const exhentai = "https://exhentai.org/"

var site string

// Extract data
func Extract(URL string) ([]static.Data, error) {
	if strings.Contains(URL, "https://ex") {
		site = exhentai
	} else {
		site = ehentai
	}

	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, errors.New("[E-Hentai] no vaild URL found")
	}

	data := []static.Data{}
	for _, URL := range URLs {
		rData, err := extractData(URL)
		if err != nil {
			return nil, err
		}
		data = append(data, rData...)
	}
	return data, nil
}

// ParseURL to gallery URL
func ParseURL(URL string) []string {
	if strings.Contains(URL, "https://ex") {
		site = exhentai
	} else {
		site = ehentai
	}

	if strings.Contains(URL, site+"g/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(fmt.Sprintf("%sg/[^\"]+", site))
	galleries := re.FindAllStringSubmatch(htmlString, -1)
	if len(galleries) == 0 {
		return []string{}
	}

	out := []string{}

	for _, gallery := range galleries {
		out = append(out, gallery[0])
	}
	return out
}

func extractData(URL string) ([]static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, errors.New("[E-Hentai] unvaild URL")
	}

	re := regexp.MustCompile("([0-9]+) pages")
	htmlNumberOfPages := re.FindStringSubmatch(htmlString)
	if len(htmlNumberOfPages) != 2 {
		return nil, errors.New("[E-Hentai] error while trying to access the gallery images")
	}
	numberOfPages, err := strconv.Atoi(htmlNumberOfPages[1])
	if err != nil {
		return nil, errors.New("[E-Hentai] couldn't get number of pages")
	}

	re = regexp.MustCompile(fmt.Sprintf("%ss[^\"]+-[0-9]+", site))
	matchedImgURLs := re.FindAllStringSubmatch(htmlString, -1)
	imgURLs := []string{}
	for _, imgURL := range matchedImgURLs {
		imgURLs = append(imgURLs, imgURL[0])
	}

	for page := 1; len(imgURLs) < numberOfPages; page++ {
		htmlString, err := request.Get(fmt.Sprintf("%s?p=%b", URL, page))
		if err != nil {
			return nil, errors.New("[E-Hentai] unvaild page URL")
		}
		imgURLs = append(imgURLs, re.FindStringSubmatch(htmlString)...)
	}

	data := []static.Data{}
	for idx, URL := range imgURLs {
		htmlString, err := request.Get(URL)
		if err != nil {
			return nil, errors.New("[E-Hentai] unvaild image URL")
		}

		re := regexp.MustCompile("<h1>([^<]+)")
		matchedTitle := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedTitle) == 0 {
			return nil, errors.New("[E-Hentai] unvaild image title")
		}

		re = regexp.MustCompile("<div>[^.]+([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Z]{2})")
		matchedFileInfo := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("[E-Hentai] unvaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		re = regexp.MustCompile(fmt.Sprintf("%sfullimg[^\"]+", site))
		srcURL := re.FindStringSubmatch(htmlString)
		if len(srcURL) != 1 {

			// sometimes the "full image url is not provided"
			re = regexp.MustCompile("<img id=\"img\" src=\"([^\"]+)")
			matchedSrcURL := re.FindAllStringSubmatch(htmlString, -1)
			if len(matchedSrcURL) != 1 {
				return nil, errors.New("[E-Hentai] unvaild image src")
			}
			srcURL = []string{matchedSrcURL[0][1]}
		}

		// size will be empty if err occurs
		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		data = append(data, static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s - %b", matchedTitle[0][1], idx+1),
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						{
							URL: srcURL[0],
							Ext: fileInfo[1],
						},
					},
					Quality: fileInfo[2],
					// ex						735       KB 	== 735000Bytes
					Size: utils.CalcSizeInByte(fSize, fileInfo[4]),
				},
			},
			Url: URL,
		})

	}

	return data, nil
}
