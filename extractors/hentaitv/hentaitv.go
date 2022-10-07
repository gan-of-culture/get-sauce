package hentaitv

import (
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaitv.fun/"

var reVideoURL = regexp.MustCompile(`https://hentaitv.fun/ads/[^"]+`)
var reVideoSource = regexp.MustCompile(`file[^\w]+([^"]+)`)

type extractor struct{}

// New returns a hentaitv extractor.
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
	reEpisodeURL := regexp.MustCompile(site + `episode/\d+/[^"/]+`)
	reParseURLShow := regexp.MustCompile(site + `hentai/\d+/[^/]+/`)

	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if reParseURLShow.MatchString(URL) {
		htmlString = strings.Split(htmlString, `<div id="comments"`)[0]
		return utils.RemoveAdjDuplicates(reEpisodeURL.FindAllString(htmlString, -1))
	}

	// contains list of show that need to be derefenced to episode level
	htmlString = strings.Split(htmlString, `<div id="sidebar">`)[0]

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

	title := utils.GetH1(&htmlString, -1)

	videoURL := reVideoURL.FindString(htmlString)
	if videoURL == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	base64String := utils.GetLastItemString(strings.Split(videoURL, "/"))
	videoURLAsBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	htmlString, err = request.Get(string(videoURLAsBytes))
	if err != nil {
		return nil, err
	}

	videoURL = utils.GetLastItemString(reVideoSource.FindStringSubmatch(htmlString))
	if videoURL == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	size, _ := request.Size(videoURL, URL)

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					{
						URL: videoURL,
						Ext: "mp4",
					},
				},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}
