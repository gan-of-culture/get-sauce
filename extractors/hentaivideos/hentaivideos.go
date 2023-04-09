package hentaivideos

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaivideos.net/"

var reEpisodeURL = regexp.MustCompile(site + `[^/]*?episode-\d+`)
var reVideoSource = regexp.MustCompile(`[^"]+\.mp4`)

type extractor struct{}

// New returns a hentaivideos extractor
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
	if reEpisodeURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	return reEpisodeURL.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	videoSource := reVideoSource.FindString(htmlString)
	if videoSource == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	size, err := request.Size(videoSource, site)
	if err != nil {
		return nil, err
	}

	return &static.Data{
		Site:  site,
		Title: strings.TrimSpace(utils.GetH1(&htmlString, 0)),
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					{
						URL: videoSource,
						Ext: utils.GetFileExt(videoSource),
					},
				},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}
