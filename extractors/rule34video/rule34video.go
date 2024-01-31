package rule34video

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/kvsplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://rule34video.com/"

type extractor struct{}

// New returns a rule34video extractor
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
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if strings.HasPrefix(URL, site+"videos/") {
		return []string{URL}
	}

	re := regexp.MustCompile(site + `video/\d+/[^?"]+`)
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	data, err := kvsplayer.ExtractFromHTML(&htmlString)
	if err != nil {
		return nil, err
	}

	data[0].Site = site
	data[0].Title = utils.GetSectionHeadingElement(&htmlString, 1, -1)
	data[0].URL = URL

	return data[0], nil

}
