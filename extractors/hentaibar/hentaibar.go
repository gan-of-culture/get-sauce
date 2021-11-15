package hentaibar

import (
	"regexp"

	"github.com/gan-of-culture/get-sauce/v2/extractors/kvsplayer"
	"github.com/gan-of-culture/get-sauce/v2/request"
	"github.com/gan-of-culture/get-sauce/v2/static"
	"github.com/gan-of-culture/get-sauce/v2/utils"
)

var reSingleURL = regexp.MustCompile(`https://hentaibar.com/videos/\d+/[^/]+`)

const site = "https://hentaibar.com"

type extractor struct{}

// New returns a hentaibar extractor.
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
	data[0].Title = utils.GetMeta(&htmlString, "og:title")
	data[0].Url = URL

	return *data[0], nil
}
