package hentaifoundry

import (
	"fmt"
	"regexp"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://www.hentai-foundry.com/"
const enterBlockURL = "https://www.hentai-foundry.com/?enterAgree=1&size=0"

var rePost = regexp.MustCompile(`pictures/user/[^/]+/\d+/[^"]+`)
var reImg = regexp.MustCompile(`<img width="(\d+)" height="(\d+)" [^/]+//(pictures\.hentai-foundry\.com[^\d]+[^/]+/([^\.]+)\.([^"]+))`)
var jar request.Myjar

type extractor struct{}

// New returns a hentai-foundry extractor.
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
	jar.New()
	_, err := request.GetWithCookies(enterBlockURL, &jar)
	if err != nil {
		return nil
	}

	if rePost.MatchString(URL) {
		return []string{URL}
	}

	var out []string
	var idx int
	for {
		idx++

		htmlString, err := request.GetWithCookies(fmt.Sprintf("%s?page=%d", URL, idx), &jar)
		if err != nil {
			return nil
		}

		URLsPerpage := utils.RemoveAdjDuplicates(rePost.FindAllString(htmlString, -1))
		for _, match := range URLsPerpage {
			out = append(out, site+match)
		}

		if len(out) >= config.Amount || len(URLsPerpage) < 25 {
			break
		}

	}

	return out
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.GetWithCookies(URL, &jar)
	if err != nil {
		return nil, err
	}

	image := reImg.FindStringSubmatch(htmlString) // 1=width 2=height 3=srcURL 4=filename 5=ext
	if len(image) == 0 {
		return nil, static.ErrDataSourceParseFailed
	}

	size, _ := request.Size(URL, site)

	return &static.Data{
		Site:  site,
		Title: image[4],
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: []*static.URL{
					{
						URL: "https://" + image[3],
						Ext: image[5],
					},
				},
				Quality: fmt.Sprintf("%sx%s", image[1], image[2]),
				Size:    size,
			},
		},
		URL: URL,
	}, nil
}
