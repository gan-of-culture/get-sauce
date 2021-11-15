package hentaipulse

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/v2/request"
	"github.com/gan-of-culture/get-sauce/v2/static"
	"github.com/gan-of-culture/get-sauce/v2/utils"
)

const site = "https://hentaipulse.com"

var reTitle = regexp.MustCompile(`<link rel="canonical" href="https://hentaipulse.com/([^/]+)`)
var reEpisodes = regexp.MustCompile(`post-\d+"`)
var reSourceURL = regexp.MustCompile(`main_video_url"[^"]+"([^"]+)`)

type extractor struct{}

// New returns a hentaipulse extractor.
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
	if strings.Count(URL, "/") == 4 {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	out := []string{}
	for _, post := range utils.RemoveAdjDuplicates(reEpisodes.FindAllString(htmlString, -1)) {
		out = append(out, site+"?p="+utils.GetLastItemString(strings.Split(post, "-")))
	}

	return out
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := utils.GetLastItemString(reTitle.FindStringSubmatch(htmlString))

	sourceURL := utils.GetLastItemString(reSourceURL.FindStringSubmatch(htmlString))
	if sourceURL == "" {
		return static.Data{}, static.ErrDataSourceParseFailed
	}

	size, _ := request.Size(sourceURL, site+"/")

	return static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					{
						URL: sourceURL,
						Ext: utils.GetFileExt(sourceURL),
					},
				},
				Size: size,
			},
		},
		Url: site + "/" + title,
	}, nil
}
