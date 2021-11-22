package hentai2read

import (
	"encoding/json"
	"html"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
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

var reJSONString = regexp.MustCompile(`{\s*'title'[\s\S]*?}`)
var reTitle = regexp.MustCompile(`[^[(|]*`)

type extractor struct{}

// New returns a hentai2read extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
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

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL + "1/")
	if err != nil {
		return nil, err
	}

	jsonString := strings.ReplaceAll(reJSONString.FindString(htmlString), "'", `"`)

	galleryData := gData{}
	err = json.Unmarshal([]byte(jsonString), &galleryData)
	if err != nil {
		return nil, err
	}

	title := html.UnescapeString(strings.TrimSpace(reTitle.FindString(galleryData.Title)))

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: buildFullImgURL(galleryData.Images),
				Size: 0,
			},
		},
		URL: URL,
	}, nil
}

func buildFullImgURL(URIParts []string) []*static.URL {
	out := []*static.URL{}
	for _, URIPart := range URIParts {
		out = append(out, &static.URL{
			URL: cdn + URIPart,
			Ext: utils.GetLastItemString(strings.Split(URIPart, ".")),
		})
	}
	return out
}
