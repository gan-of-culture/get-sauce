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

	re = regexp.MustCompile(`<h1.*?pretty">([^>]*)<`)
	title := re.FindStringSubmatch(htmlString)[1]
	title = strings.ReplaceAll(title, "?", "")

	URLs := []static.URL{}
	var quality string

	for _, page := range pages {
		pageURL := fmt.Sprintf("https://nhentai.net/g/%s/%d", magicNumber, page)
		time.Sleep(250 * time.Millisecond)
		htmlString, err := request.Get(pageURL)
		if err != nil {
			continue
		}
		// some times you need to retry
		if strings.Contains(htmlString, "<title>503 Service Temporarily Unavailable</title>") || strings.Contains(htmlString, "<title>429 Too Many Requests</title>") {
			time.Sleep(400 * time.Millisecond)
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
		0: {
			Site:  site,
			Title: title,
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": {
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
