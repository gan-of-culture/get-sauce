package orzqwq

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://orzqwq.com/"

var reImageURL = regexp.MustCompile(`image-0" src="([^"]+/)\d+\.(\w{3})`)
var reNumPages = regexp.MustCompile(`(\d+) pages`)

type extractor struct{}

// New returns a orzqwq extractor.
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
	if strings.HasPrefix(URL, site+"manga") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://orzqwq\.com/manga/[^/]+`)
	return utils.RemoveAdjDuplicates(re.FindAllString(htmlString, -1))
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	imageTemplateURL := reImageURL.FindStringSubmatch(htmlString)
	if len(imageTemplateURL) < 3 {
		return nil, static.ErrDataSourceParseFailed
	}

	matchedNumPages := reNumPages.FindStringSubmatch(htmlString)
	if len(matchedNumPages) < 2 {
		return nil, static.ErrDataSourceParseFailed
	}

	numberOfPages, err := strconv.Atoi(matchedNumPages[1])
	if err != nil {
		return nil, err
	}

	URLs := []*static.URL{}
	for page := 1; page <= numberOfPages; page++ {
		URLs = append(URLs, &static.URL{
			URL: fmt.Sprintf("%s%03d.%s", imageTemplateURL[1], page, imageTemplateURL[2]),
			Ext: imageTemplateURL[2],
		})
	}

	return &static.Data{
		Site:  site,
		Title: strings.TrimSpace(utils.GetH1(&htmlString, -1)),
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: URL,
	}, nil
}
