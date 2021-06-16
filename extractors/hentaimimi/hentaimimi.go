package hentaimimi

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://hentaimimi.com/"

var reTitle *regexp.Regexp
var reImgData *regexp.Regexp
var reImgExt *regexp.Regexp

func init() {
	reTitle = regexp.MustCompile(`<meta name="title" content="([^"]*)`)
	reImgData = regexp.MustCompile(`\["uploads.*?]`)
	reImgExt = regexp.MustCompile(`\w+$`)
}

type extractor struct{}

// New returns a hentaimimi extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	IDs := parseURL(URL)

	data := []*static.Data{}
	for _, u := range IDs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	re := regexp.MustCompile(`/view/(\d+)`)
	matchedID := re.FindStringSubmatch(URL)
	if len(matchedID) == 2 {
		return []string{matchedID[1]}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	IDs := []string{}
	for _, v := range re.FindAllStringSubmatch(htmlString, -1) {
		IDs = append(IDs, v[1])
	}
	return utils.RemoveAdjDuplicates(IDs)
}

func extractData(ID string) (*static.Data, error) {
	URL := fmt.Sprintf("%sview/%s", site, ID)

	htmlString, err := request.Get(URL)
	if err != nil {
		return &static.Data{}, err
	}

	title := reTitle.FindStringSubmatch(htmlString)
	if len(title) < 1 {
		return &static.Data{}, fmt.Errorf("no title found for: %s", URL)
	}

	jsonStr := reImgData.FindString(htmlString)
	if jsonStr == "" {
		return &static.Data{}, fmt.Errorf("no image links found for: %s", URL)
	}

	imgURLs := []string{}
	err = json.Unmarshal([]byte(jsonStr), &imgURLs)
	if err != nil {
		return &static.Data{}, err
	}

	URLs := []*static.URL{}

	pages := utils.NeedDownloadList(len(imgURLs))
	for _, i := range pages {
		URLs = append(URLs, &static.URL{
			URL: site + imgURLs[i-1],
			Ext: reImgExt.FindString(imgURLs[i-1]),
		})
	}

	return &static.Data{
		Site:  site,
		Title: html.UnescapeString(title[1]),
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: URLs,
				Info: fmt.Sprintf("Pages: %d", len(imgURLs)),
			},
		},
		Url: URL,
	}, nil
}
