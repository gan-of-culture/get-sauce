package hentaihavenred

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/nhgroup"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaihaven.red/"

var reEpisodeURL = regexp.MustCompile(site + `hentai/[\w-%]+/`)
var reParseURLShow = regexp.MustCompile(site + `watch/[\w-%]+/`)

type extractor struct{}

// New returns a hentaihavenred extractor.
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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if strings.Contains(URL, "/watch/") {
		htmlString = strings.Split(htmlString, `<div class="bixbox"`)[0]
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

	title := strings.Split(utils.GetMeta(&htmlString, "og:title"), " - ")[0]
	title = strings.Split(title, " | ")[0]

	data, err := nhgroup.ExtractData(URL)
	if err != nil {
		return nil, err
	}

	data.Site = site
	data.Title = title
	data.URL = URL
	return data, nil
}
