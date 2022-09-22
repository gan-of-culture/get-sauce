package hentaimoon

import (
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/extractors/kvsplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reSingleURL = regexp.MustCompile(`https://hentai-moon.com/videos/\d+/[^/]+`)
var reSubtitles = regexp.MustCompile(`//.*?\.vtt`)

const site = "https://hentai-moon.com"

type extractor struct{}

// New returns a hentai-moon extractor.
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
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if reSingleURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, _ := request.Get(URL)

	return reSingleURL.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {

	req, _ := http.NewRequest(http.MethodGet, URL, nil)

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	htmlString := string(body)

	data, err := kvsplayer.ExtractFromHTML(&htmlString)
	if err != nil {
		return nil, err
	}

	data[0].Site = site
	data[0].Title = utils.GetH1(&htmlString, -1)

	matchedSubtitleURL := reSubtitles.FindString(htmlString)
	if matchedSubtitleURL != "" {
		subtitleURL, err := url.Parse(matchedSubtitleURL)
		if err != nil {
			return nil, err
		}
		data[0].Captions = append(data[0].Captions, &static.Caption{
			URL: static.URL{
				URL: "https:" + subtitleURL.String(),
				Ext: utils.GetFileExt(matchedSubtitleURL),
			},
			Language: "English",
		})
	}

	data[0].URL = URL

	return data[0], nil
}
