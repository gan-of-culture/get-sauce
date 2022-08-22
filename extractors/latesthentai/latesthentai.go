package latesthentai

import (
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://latesthentai.com/"

var reEpisodeURL = regexp.MustCompile(site + `watch/[^"]+`)
var reVideoURL = regexp.MustCompile(`[^"]+htstreaming[^"]+`)

type extractor struct{}

// New returns a latesthentai extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)

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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	out := []string{}
	for _, anime := range reEpisodeURL.FindAllString(htmlString, -1) {
		out = append(out, parseURL(anime)...)
	}
	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	htStreamingURL := reVideoURL.FindString(htmlString)
	if htStreamingURL == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	data, err := htstreaming.ExtractData(htStreamingURL)
	if err != nil {
		return nil, err
	}
	data.Site = site
	data.Title = utils.GetH1(&htmlString, -1)
	data.URL = URL
	return data, nil
}
