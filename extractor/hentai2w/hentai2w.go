package hentai2w

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://hentai2w.com/"

func ParseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"video/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`[^"]*/video/[^"]*`)
	return re.FindAllString(htmlString, -1)
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Hentai2w] No scrapable URL found for %s", URL)
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`<source src="([^"]+)"`)
	videoURL := utils.GetLastItemString(re.FindStringSubmatch(htmlString))
	if videoURL == "" || strings.HasPrefix(videoURL, "<") {
		return static.Data{}, fmt.Errorf("[Hentai2w] No videoURL found for %s", URL)
	}
	ext := utils.GetLastItemString(strings.Split(videoURL, "."))

	size, _ := request.Size(URL, site)

	return static.Data{
		Site:  site,
		Title: utils.GetMeta(&htmlString, "og:title"),
		Type:  utils.GetMediaType(ext),
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: videoURL,
						Ext: ext,
					},
				},
				Quality: "unknown",
				Size:    size,
			},
		},
		Url: URL,
	}, nil
}
