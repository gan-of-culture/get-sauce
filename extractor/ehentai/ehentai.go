package ehentai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://e-hentai.org/"

// Extract data
func Extract(URL string) ([]static.Data, error) {
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
	if strings.Contains(URL, "https://e-hentai.org/g/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://e-hentai.org/g/[^"]+`)
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
		return nil, errors.New("[E-Hentai] invaild URL")
	}

	if strings.Contains(htmlString, "<h1>Content Warning</h1>") {
		if config.RestrictContent {
			return []static.Data{}, errors.New("[E-hentai] Restricted content")
		}
		return extractData(URL + "?nw=session")
	}

	re := regexp.MustCompile(`([0-9]+) pages`)
	htmlNumberOfPages := re.FindStringSubmatch(htmlString)
	if len(htmlNumberOfPages) != 2 {
		return nil, errors.New("[E-Hentai] error while trying to access the gallery images")
	}
	numberOfPages, err := strconv.Atoi(htmlNumberOfPages[1])
	if err != nil {
		return nil, errors.New("[E-Hentai] couldn't get number of pages")
	}

	re = regexp.MustCompile(`https://e-hentai.org/s[^"]+-[0-9]+`)
	imgURLs := re.FindAllString(htmlString, -1)

	// if gallery has more than 40 images -> walk other pages for links aswell
	for page := 1; len(imgURLs) < numberOfPages; page++ {
		htmlString, err := request.Get(fmt.Sprintf("%s?p=%d", URL, page))
		if err != nil {
			return nil, errors.New("[E-Hentai] invaild page URL")
		}
		imgURLs = append(imgURLs, re.FindAllString(htmlString, -1)...)
	}

	data := []static.Data{}
	for _, idx := range utils.NeedDownloadList(len(imgURLs)) {
		htmlString, err := request.Get(imgURLs[idx-1])
		if err != nil {
			return nil, errors.New("[E-Hentai] unvaild image URL")
		}

		title := utils.GetH1(&htmlString, 0)
		if title == "" {
			return nil, errors.New("[E-Hentai] unvaild image title")
		}

		re := regexp.MustCompile(`<div>[^.]+\.([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Z]{2})`)
		matchedFileInfo := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("[E-Hentai] unvaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		// sometimes the "full image url is not provided"
		re = regexp.MustCompile(`<img id="img" src="([^"]+)`)
		matchedSrcURL := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedSrcURL) != 1 {
			return nil, errors.New("[E-Hentai] unvaild image src")
		}

		// size will be empty if err occurs
		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		data = append(data, static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s - %d", title, idx),
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						{
							URL: matchedSrcURL[0][1],
							Ext: fileInfo[1],
						},
					},
					Quality: fileInfo[2],
					// ex						735       KB 	== 735000Bytes
					Size: utils.CalcSizeInByte(fSize, fileInfo[4]),
				},
			},
			Url: imgURLs[idx-1],
		})

	}

	return data, nil
}
