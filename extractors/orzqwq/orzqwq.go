package orzqwq

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://orzqwq.com/"

var reImageURL = regexp.MustCompile(`image-\d+"\s.+data-src="([^"]+)`)

type extractor struct{}

// New returns a orzqwq extractor
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

	re := regexp.MustCompile(strings.ToLower(URL) + `p(\d+)`)
	pages := []int{}
	for _, pageNum := range re.FindAllStringSubmatch(htmlString, -1) {
		numAsInt, err := strconv.Atoi(pageNum[1])
		if err != nil {
			return nil, err
		}
		pages = append(pages, numAsInt)
	}

	sort.Ints(pages)
	allImageURLs := []*static.URL{}
	for _, pageNum := range utils.RemoveAdjDuplicates(pages) {
		pageHtml, err := request.Get(fmt.Sprintf("%s/p%d", URL, pageNum))
		if err != nil {
			return nil, err
		}
		for _, imgUrl := range reImageURL.FindAllStringSubmatch(pageHtml, -1) {
			allImageURLs = append(allImageURLs, &static.URL{
				URL: imgUrl[1],
				Ext: utils.GetFileExt(imgUrl[1]),
			})
		}
	}

	wantedImages := utils.NeedDownloadList(len(allImageURLs))
	URLs := []*static.URL{}
	for _, i := range wantedImages {
		URLs = append(URLs, allImageURLs[i-1])
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
