package latesthentai

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://latesthentai.com/"

var reEpisodeURL = regexp.MustCompile(site + `watch/[^"]+`)
var reParseURLShow = regexp.MustCompile(site + `anime/[\w-%]+`)
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

	if strings.Contains(URL, "/anime/") {
		htmlString = strings.Split(htmlString, `<div class="eplister"`)[0]
		return utils.RemoveAdjDuplicates(reEpisodeURL.FindAllString(htmlString, -1))
	}

	// contains list of show that need to be derefenced to episode level
	htmlString = strings.Split(htmlString, `<div id="sidebar">`)[0]

	out := []string{}
	for _, anime := range reParseURLShow.FindAllString(htmlString, -1) {
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
