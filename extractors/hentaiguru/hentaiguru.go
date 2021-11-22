package hentaiguru

import (
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/jwplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reEpisodeURL = regexp.MustCompile(`https://hentai\.guru/hentai/[^/]+/episode-\d+`)
var reSeriesURL = regexp.MustCompile(`https://hentai\.guru/hentai/[^/]+`)

const site = "https://hentai.guru"

type extractor struct{}

// New returns a hentai.guru extractor.
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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if episodes := reEpisodeURL.FindAllString(htmlString, -1); len(episodes) > 0 {
		return episodes
	}

	out := []string{}
	for _, seriesURL := range reSeriesURL.FindAllString(htmlString, -1) {
		out = append(out, parseURL(seriesURL)...)
	}

	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetH1(&htmlString, -1)

	data, err := jwplayer.New().Extract(jwplayer.FindJWPlayerURL(&htmlString))
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	video := data[0]
	video.Site = site
	video.Title = title
	video.Type = static.DataTypeVideo
	video.URL = URL

	return video, nil
}
