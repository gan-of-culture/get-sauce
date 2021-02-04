package universal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

// Extract universal link
func Extract(url string, site string) ([]static.Data, error) {

	re := regexp.MustCompile("/([^/]+)\\.([a-zA-z]*)?\\??[0-9]*$")
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

	dataType := ""
	switch matches[2] {
	case "jpg", "jpeg", "png", "gif", "webp":
		dataType = fmt.Sprintf("%s/%s", "image", matches[2])
	case "webm", "mp4", "mkv", "m4a":
		dataType = fmt.Sprintf("%s/%s", "video", matches[2])
	default:
		dataType = fmt.Sprintf("%s/%s", "unknown", matches[2])
	}

	size, _ := request.Size(url, url)

	return []static.Data{
		0: {
			Site:  site,
			Title: matches[1],
			Type:  dataType,
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						0: {
							URL: url,
							Ext: matches[2],
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
