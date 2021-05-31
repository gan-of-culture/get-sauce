package hentai2read

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type gData struct {
	Title        string
	Index        int
	Images       []string
	PreloadLimit int
	MainURL      string
}

const site = "https://hentai2read.com/"
const cdn = "https://static.hentaicdn.com/hentai"

func ParseURL(URL string) []string {
	URL = strings.Split(URL, "#")[0]
	if strings.Contains(URL, "_") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`([^/]*/)" class="title"`)

	URLs := []string{}
	for _, u := range re.FindAllStringSubmatch(htmlString, -1) {
		URLs = append(URLs, site+u[1])
	}
	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Hentai2read] No scrapable URL found for %s", URL)
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
	htmlString, err := request.Get(URL + "1/")
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`{\s*'title'[\s\S]*?}`)
	jsonString := strings.ReplaceAll(re.FindString(htmlString), "'", `"`)

	galleryData := gData{}
	err = json.Unmarshal([]byte(jsonString), &galleryData)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`[^[(|]*`)
	title := html.UnescapeString(strings.TrimSpace(re.FindString(galleryData.Title)))

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "image",
		Streams: map[string]static.Stream{
			"0": {
				URLs:    buildFullImgURL(galleryData.Images),
				Quality: "unkown",
				Size:    0,
			},
		},
		Url: URL,
	}, nil
}

func buildFullImgURL(URIParts []string) []static.URL {
	out := []static.URL{}
	for _, URIPart := range URIParts {
		out = append(out, static.URL{
			URL: cdn + URIPart,
			Ext: utils.GetLastItemString(strings.Split(URIPart, ".")),
		})
	}
	return out
}
