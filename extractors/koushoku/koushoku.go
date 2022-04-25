package koushoku

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://koushoku.org/"
const cdn = "https://cdn.koushoku.org/data/"

var reGalleryURLPart = regexp.MustCompile(`archive/(\d+)/[\w-]+/?(\d+)?`)
var reGallerySize = regexp.MustCompile(`\d+.\d`)

type extractor struct{}

// New returns a koushoku extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)

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
	matchedGallery := reGalleryURLPart.FindStringSubmatch(URL)
	switch len(matchedGallery) {
	case 2:
		return []string{site + matchedGallery[0]}
	case 3:
		config.Pages = matchedGallery[2]
		return []string{strings.TrimSuffix(site+matchedGallery[0], matchedGallery[2])}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	out := []string{}
	for _, URLPart := range reGalleryURLPart.FindAllString(htmlString, -1) {
		out = append(out, site+URLPart)
	}

	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedGalleryID := reGalleryURLPart.FindStringSubmatch(URL)
	if len(matchedGalleryID) < 2 {
		return nil, static.ErrURLParseFailed
	}

	numPagesProperty := getGalleryProperty("pages", &htmlString)
	if numPagesProperty == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	numPages, err := strconv.Atoi(numPagesProperty)
	if err != nil {
		return nil, err
	}

	URLs := []*static.URL{}
	for _, i := range utils.NeedDownloadList(numPages) {
		URLs = append(URLs, &static.URL{
			URL: fmt.Sprintf("%s%s/%d.jpg", cdn, matchedGalleryID[1], i),
			Ext: "jpg",
		})
	}

	sizeProperty := getGalleryProperty("size", &htmlString)
	sizeFloat, _ := strconv.ParseFloat(reGallerySize.FindString(sizeProperty), 64)

	return &static.Data{
		Site:  site,
		Title: utils.GetH1(&htmlString, -1),
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
				Size: utils.CalcSizeInByte(sizeFloat, sizeProperty[len(sizeProperty)-2:]),
			},
		},
		URL: URL,
	}, nil
}

func getGalleryProperty(property string, htmlString *string) string {
	var reGalleryProperty = regexp.MustCompile(fmt.Sprintf(`<td>%s</td>\n<td>([^<]+)</td>`, strings.Title(property)))
	return utils.GetLastItemString(reGalleryProperty.FindStringSubmatch(*htmlString))
}
