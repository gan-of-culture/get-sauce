package hentaiyes

import (
	"fmt"
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaiyes.com/"

type extractor struct{}

// New returns a hentaiyes extractor.
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
	if ok, _ := regexp.MatchString(`episode-\d+[/_\-\&]`, URL); ok {
		return []string{URL}
	}

	//check if it's an overview/series page maybe
	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`<h6><a href="/(watch/[^"\s]*?episode-\d*/)`)
	matchedEpisodes := re.FindAllStringSubmatch(htmlString, -1)
	URLs := []string{}
	for _, u := range matchedEpisodes {
		URLs = append(URLs, site+u[1])
	}
	return URLs
}

func extractData(URL string) (static.Data, error) {
	re := regexp.MustCompile(`/watch/([^"\s]*?episode-\d*)/`)
	slug := re.FindStringSubmatch(URL)
	embedURL := fmt.Sprintf("%sembed.php?name=%s&source=1", site, slug[1])

	data, err := htstreaming.ExtractData(embedURL)
	if err != nil {
		return static.Data{}, err
	}
	data.Site = site
	data.Title = slug[1]
	data.Url = URL
	return data, nil
}
