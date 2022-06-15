package hanime

import (
	"github.com/gan-of-culture/get-sauce/extractors/animestream"
	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hanime.io/"

type extractor struct{}

// New returns a hanime extractor.
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
	return animestream.ParseURL(URL, site)
}

func extractData(URL string) (*static.Data, error) {

	data, err := htstreaming.ExtractData(URL)
	if err != nil {
		return nil, err
	}
	data.Site = site
	return data, nil
}
