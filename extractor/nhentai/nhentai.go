package nhentai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://nhentai.net"

// Extract data from supplied url
func Extract(url string) ([]static.Data, error) {
	magicNumber, page := ParseURL(url)
	if magicNumber == "" && page == "" {
		return nil, errors.New("[NHentai]No magic number found")
	}

	htmlString, err := request.Get(fmt.Sprintf("https://nhentai.net/g/%s/", magicNumber))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("gallerythumb")
	pages := utils.NeedDownloadList(len(re.FindAllString(htmlString, -1)))

	if page != "" {
		pageNo, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		pages = []int{pageNo}
	}

	re = regexp.MustCompile("<h1>([^<]+)")
	title := re.FindStringSubmatch(htmlString)[1]

	URLs := []static.URL{}
	var quality string

	for _, page := range pages {
		pageURL := fmt.Sprintf("https://nhentai.net/g/%s/%d", magicNumber, page)
		htmlString, err := request.Get(pageURL)
		if err != nil {
			continue
		}
		// some times you need to retry
		if strings.Contains(htmlString, "<title>503 Service Temporarily Unavailable</title>") {
			time.Sleep(500 * time.Millisecond)
			htmlString, err = request.Get(pageURL)
		}

		re := regexp.MustCompile("<img src=\"([^\"]+)\" width=\"([^\"]+)\" height=\"([^\"]+)\"")
		matchedImgData := re.FindStringSubmatch(htmlString)
		if len(matchedImgData) != 4 {
			return nil, errors.New("[nhentai] Image parsing failed")
		}

		if page == 1 {
			quality = fmt.Sprintf("%s x %s", matchedImgData[2], matchedImgData[3])
		}

		URLs = append(URLs, static.URL{
			URL: matchedImgData[1],
			Ext: utils.GetLastItemString(strings.Split(matchedImgData[1], ".")),
		})
	}

	return []static.Data{
		1: static.Data{
			Site:  site,
			Title: title,
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": static.Stream{
					URLs:    URLs,
					Quality: quality,
					Size:    0,
				},
			},
			Url: url,
		},
	}, nil
}

// ParseURL data
func ParseURL(url string) (string, string) {
	re := regexp.MustCompile("[0-9]+")
	urlNumbers := re.FindAllString(url, -1)

	if len(urlNumbers) <= 0 {
		return "", ""
	}

	// if there are two "int" values it means the exact page was supplied
	var page string
	if len(urlNumbers) > 1 {
		page = urlNumbers[1]
	}

	return urlNumbers[0], page
}

func extractImageData(URL string, stream static.Stream) (map[string]static.Stream, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}
	// some times you need to retry
	if strings.Contains(htmlString, "<title>503 Service Temporarily Unavailable</title>") {
		time.Sleep(500 * time.Millisecond)
		htmlString, err = request.Get(URL)
	}

	re := regexp.MustCompile("<img src=\"([^\"]+)\" width=\"([^\"]+)\" height=\"([^\"]+)\"")
	matchedImgData := re.FindStringSubmatch(htmlString)
	if len(matchedImgData) != 4 {
		return nil, errors.New("[nhentai] Image parsing failed")
	}
	// [1] src url [2] width [3] height

	// currently not needed
	/*size, err := request.Size(matchedImgData[1], URL)
	if err != nil {
		return nil, err
	}*/
	return map[string]static.Stream{
		"0": static.Stream{
			URLs: []static.URL{
				{
					URL: matchedImgData[1],
					Ext: utils.GetLastItemString(strings.Split(matchedImgData[1], ".")),
				},
			},
			Quality: fmt.Sprintf("%s x %s", matchedImgData[2], matchedImgData[3]),
			Size:    0,
		},
	}, nil
}
