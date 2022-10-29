package rokuhentai

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://rokuhentai.com/"
const cdn = "https://rokuhentai.com/_images/pages/%s/%d.jpg"

var reID = regexp.MustCompile(fmt.Sprintf("%s(\\w*)[/^\"]*", site))
var rePageInfo = regexp.MustCompile(`(\d+) images`)

type extractor struct{}

// New returns a rokuhentai extractor.
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
	if ok, _ := regexp.MatchString(fmt.Sprintf("%s\\w*[/^\"]*", site), URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(fmt.Sprintf("%s\\w*[/^\"]*\">", site))

	URLs := []string{}
	for _, v := range re.FindAllString(htmlString, -1) {
		URLs = append(URLs, strings.TrimSuffix(v, ">"))
	}
	return URLs
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	ID := utils.GetLastItemString(reID.FindStringSubmatch(URL))

	matchedPageInfo := rePageInfo.FindStringSubmatch(htmlString)
	pages, _ := strconv.Atoi(matchedPageInfo[1])

	URLs := []*static.URL{}
	for _, pageNumber := range utils.NeedDownloadList(pages) {
		URLs = append(URLs, &static.URL{
			URL: fmt.Sprintf(cdn, ID, pageNumber-1),
			Ext: "jpg",
		})
	}

	return &static.Data{
		Site:  site,
		Title: strings.TrimSpace(strings.Split(utils.GetMeta(&htmlString, "og:title"), " - Roku Hentai")[0]),
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
