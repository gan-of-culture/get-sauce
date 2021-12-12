package haho

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://haho.moe"

var reEpisodeURL = regexp.MustCompile(site + `/anime/\w+/\d+`)
var reSeriesURL = regexp.MustCompile(site + `/anime/\w{8}"`)
var reEmbedURL = regexp.MustCompile(site + `/embed\?v=[^"]+`)
var reSource = regexp.MustCompile(`<source.+src="([^"]+)" title="([^"]+)" type="([^"]+)`) // 1=URL 2=Quality 3=Type

type extractor struct{}

// New returns a hentaihd extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract data from URL
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
	if reEpisodeURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	// cut of the recommended anime on the site
	htmlString = strings.Split(htmlString, `<section id="entry-genretree">`)[0]

	// cut of hidden items
	htmlString = strings.Split(htmlString, `</main>`)[0]

	if matchedEpisodes := reEpisodeURL.FindAllString(htmlString, -1); len(matchedEpisodes) > 0 {
		return matchedEpisodes
	}

	out := []string{}
	for _, series := range reSeriesURL.FindAllString(htmlString, -1) {
		out = append(out, parseURL(strings.TrimSuffix(series, `"`))...)
	}

	return out
}
func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetH1(&htmlString, -1)

	matchedEmbedURL := reEmbedURL.FindString(htmlString)
	if matchedEmbedURL == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	embedHtml, err := request.GetWithHeaders(matchedEmbedURL, map[string]string{"Referer": URL})
	if err != nil {
		return nil, err
	}

	sources := reSource.FindAllStringSubmatch(embedHtml, -1)
	streams := map[string]*static.Stream{}
	for i, source := range sources {
		if len(source) < 4 {
			return nil, static.ErrDataSourceParseFailed
		}
		size, _ := request.Size(source[1], site)
		streams[fmt.Sprint(i)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: source[1],
					Ext: strings.Split(source[3], "/")[1],
				},
			},
			Quality: source[2],
			Size:    size,
		}
	}

	return &static.Data{
		Site:    site,
		Title:   title,
		Type:    static.DataTypeVideo,
		Streams: streams,
		URL:     URL,
	}, nil
}
