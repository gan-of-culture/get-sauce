package universal

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

// Extract universal link
func Extract(url string, site string) ([]static.Data, error) {

	re := regexp.MustCompile("/([^/]+)\\.([a-zA-z0-9]*)?\\??[0-9a-zA-Z&=]*$")
	// matches[1] = title, matches[2] = fileext
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return []static.Data{
			0: {
				Site:  site,
				Title: "unknown",
				Type:  "unknown",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							0: {
								URL: url,
								Ext: utils.GetLastItemString(strings.Split(url, ".")),
							},
						},
						Quality: "unknown",
						Size:    0,
					},
				},
				Url: url,
			},
		}, nil
	}

	size, _ := request.Size(url, url)
	ext := matches[2]
	if ext == "m3u8" || ext == "txt" {
		ext = "ts"
	}

	return []static.Data{
		0: {
			Site:  site,
			Title: matches[1],
			Type:  utils.GetMediaType(matches[2]),
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						0: {
							URL: url,
							Ext: ext,
						},
					},
					Quality: "unknown",
					Size:    size,
				},
			},
			Url: url,
		},
	}, nil
}
