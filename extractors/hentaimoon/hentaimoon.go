package hentaimoon

import (
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/kvsplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reSingleURL = regexp.MustCompile(`https://hentai-moon.com/videos/\d+/[^/]+`)
var reSubtitles = regexp.MustCompile(`//.*?\.vtt`)

const site = "https://hentai-moon.com"

type extractor struct{}

// New returns a hentai-moon extractor.
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
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if reSingleURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, _ := request.Get(URL)

	return reSingleURL.FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	data, err := kvsplayer.ExtractFromHTML(&htmlString)
	if err != nil {
		return static.Data{}, err
	}

	data[0].Site = site
	data[0].Title = utils.GetH1(&htmlString, -1)

	subtitleURL := reSubtitles.FindString(htmlString)
	if subtitleURL != "" {
		data[0].Captions = append(data[0].Captions, &static.Caption{
			URL: static.URL{
				URL: "https:" + subtitleURL,
				Ext: utils.GetFileExt(subtitleURL),
			},
			Language: "English",
		})
	}

	data[0].Url = URL

	return *data[0], nil
}
