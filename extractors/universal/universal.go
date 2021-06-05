package universal

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/extractors/imgboard"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type extractor struct{}

// New returns a universal extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract unviersal url
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	data, err := imgboard.New().Extract(URL)
	if len(data) > 0 && err == nil {
		return data, nil
	}

	u, _ := url.Parse(URL)

	re := regexp.MustCompile(`/([^/]+)\.([a-zA-z0-9]*)?\??[0-9a-zA-Z&=]*$`)
	// matches[1] = title, matches[2] = fileext
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 3 {
		return []*static.Data{
			0: {
				Site:  u.Host,
				Title: "unknown",
				Type:  "unknown",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							0: {
								URL: URL,
								Ext: utils.GetLastItemString(strings.Split(URL, ".")),
							},
						},
						Quality: "unknown",
						Size:    0,
					},
				},
				Url: URL,
			},
		}, nil
	}

	size, _ := request.Size(URL, URL)
	ext := matches[2]
	if ext == "m3u8" || ext == "txt" {
		ext = "ts"
	}

	return []*static.Data{
		0: {
			Site:  u.Host,
			Title: matches[1],
			Type:  utils.GetMediaType(matches[2]),
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						0: {
							URL: URL,
							Ext: ext,
						},
					},
					Quality: "unknown",
					Size:    size,
				},
			},
			Url: URL,
		},
	}, nil
}
