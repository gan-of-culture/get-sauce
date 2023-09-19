package hentaihavenred

import (
	"github.com/gan-of-culture/get-sauce/extractors/jwplayer"
	"github.com/gan-of-culture/get-sauce/extractors/nhgroup"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaihaven.red/"

type extractor struct{}

// New returns a hentaihaven.red extractor
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
	return nhgroup.ParseURL(URL)
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	data, err := jwplayer.New().Extract(jwplayer.FindJWPlayerURL(&htmlString))
	if err != nil {
		return nil, err
	}

	data[0].Site = site
	data[0].Title = utils.GetSectionHeadingElement(&htmlString, 1, -1)
	data[0].URL = URL

	return data[0], nil

}
