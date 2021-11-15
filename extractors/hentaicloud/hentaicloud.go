package hentaicloud

import (
	"fmt"
	"regexp"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://www.hentaicloud.com/"

type extractor struct{}

// New returns a hentaicloud extractor.
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
	if ok, _ := regexp.MatchString(`https://www.hentaicloud.com/video/\d*/[^/]*/episode\d*/*`, URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}
	re := regexp.MustCompile(`video/\d*/[^/]*/episode\d*/[^"]*`)
	URLs := []string{}
	for i, v := range re.FindAllString(htmlString, -1) {
		if i%2 == 0 {
			URLs = append(URLs, site+v)
		}
	}

	return URLs
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}
	title := utils.GetMeta(&htmlString, "og:title")

	re := regexp.MustCompile(`source src="(https://www.hentaicloud.com/media/videos/[^.]+([^"]+)).+res="([^"]*)`)
	srcTags := re.FindAllStringSubmatch(htmlString, -1) //1=URL 2=ext 3=resolution
	if len(srcTags) < 1 {
		return static.Data{}, static.ErrDataSourceParseFailed
	}

	streams := map[string]*static.Stream{}
	dataLen := len(srcTags)
	for i, source := range srcTags {
		if len(source) != 4 {
			return static.Data{}, static.ErrDataSourceParseFailed
		}

		size, _ := request.Size(source[1], URL)

		streams[fmt.Sprint(dataLen-i-1)] = &static.Stream{
			URLs: []*static.URL{
				{
					URL: source[1],
					Ext: source[2],
				},
			},
			Quality: source[3] + "p",
			Size:    size,
		}
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: streams,
		Url:     URL,
	}, nil
}
