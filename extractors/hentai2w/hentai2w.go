package hentai2w

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/v2/request"
	"github.com/gan-of-culture/get-sauce/v2/static"
	"github.com/gan-of-culture/get-sauce/v2/utils"
)

const site = "https://hentai2w.com/"

type extractor struct{}

// New returns a hentai2w extractor.
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
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
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

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`<source.*src="([^"]+)"`)
	videoURL := utils.GetLastItemString(re.FindStringSubmatch(htmlString))
	if videoURL == "" || strings.HasPrefix(videoURL, "<") {
		return static.Data{}, static.ErrDataSourceParseFailed
	}
	ext := utils.GetLastItemString(strings.Split(videoURL, "."))

	size, _ := request.Size(URL, site)

	return static.Data{
		Site:  site,
		Title: utils.GetMeta(&htmlString, "og:title"),
		Type:  "video",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					{
						URL: videoURL,
						Ext: ext,
					},
				},
				Size: size,
			},
		},
		Url: URL,
	}, nil
}
