package nhplayer

import (
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var rePlayerURL = regexp.MustCompile(`https://nhplayer\.com/v/[^/"]+`)
var reHTStreamingVideoURL = regexp.MustCompile(`https://htstreaming.com/video/([^"]*)`)

type extractor struct{}

// New returns a nhplayer extractor.
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
	if !rePlayerURL.MatchString(URL) {
		return nil

	}
	return []string{URL}
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	videoURL := reHTStreamingVideoURL.FindString(htmlString)
	if videoURL == "" {
		return nil, static.ErrURLParseFailed
	}

	return htstreaming.ExtractData(videoURL)
}
