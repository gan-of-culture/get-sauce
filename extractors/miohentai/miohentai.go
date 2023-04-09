package miohentai

import (
	"html"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://miohentai.com/"

var reShortLink = regexp.MustCompile(site + `\?p=\d+`)
var reSourceURL = regexp.MustCompile(`[^"]*cdn\.miohentai[^"]*`)

type extractor struct{}

// New returns a miohentai extractor
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
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if matchedShortLink := reShortLink.FindString(htmlString); matchedShortLink != "" {
		return []string{matchedShortLink}
	}

	re := regexp.MustCompile(`post-\d+`)
	posts := re.FindAllString(htmlString, -1)
	URLs := []string{}
	for _, v := range utils.RemoveAdjDuplicates(posts) {
		URLs = append(URLs, site+"?p="+strings.Split(v, "-")[1])
	}
	return URLs
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := html.UnescapeString(utils.GetH1(&htmlString, -1))
	if title == "" {
		splitURL := strings.Split(utils.GetMeta(&htmlString, "og:url"), "/")
		title = splitURL[len(splitURL)-2]
	}

	srcURL := reSourceURL.FindString(htmlString)

	headers, err := request.Headers(srcURL, URL)
	if err != nil {
		return nil, err
	}
	size, err := request.GetSizeFromHeaders(&headers)
	if err != nil {
		return nil, err
	}

	ext := utils.GetLastItemString(strings.Split(srcURL, "."))
	if strings.Contains(srcURL, "index.php") {
		ext = strings.Split(headers.Get("content-type"), "/")[1]
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  utils.GetMediaType(strings.Split(headers.Get("content-type"), "/")[1]),
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					{
						URL: srcURL,
						Ext: ext,
					},
				},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}
