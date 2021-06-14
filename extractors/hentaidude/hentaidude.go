package hentaidude

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type streams struct {
	Success bool              `json:"success"`
	Sources map[string]string `json:"sources"`
}

const site = "https://hentaidude.com/"
const api = "https://hentaidude.com/wp-admin/admin-ajax.php"
const apiPost = "https://hentaidude.com/?p="

type extractor struct{}

// New returns a hentaidude extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`hentaidude\.com/.*(?:(?:/tag/)|(?:/3d-hentai[^/]*/)|(?:page/\d/)|(?:\?*orderby=)|(?:\?*tid=))`, URL); ok || URL == site {
		htmlString, err := request.Get(URL)
		if err != nil {
			return []string{}
		}
		re := regexp.MustCompile(`post-([^"]*)`)
		URLs := []string{}
		for _, v := range re.FindAllStringSubmatch(htmlString, -1) {
			URLs = append(URLs, apiPost+v[1])
		}
		return URLs
	}

	return []string{URL}
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}
	title := utils.GetMeta(&htmlString, "og:title")
	title = strings.TrimSuffix(title, " | Hentaidude.com")

	re := regexp.MustCompile(`id: '(\d*)',\s*nonce: '([^']*)`)
	matchedSourceReq := re.FindStringSubmatch(htmlString) // 1=id  2=nonce
	if len(matchedSourceReq) < 3 {
		return static.Data{}, fmt.Errorf("can't locate json params in URL: %s", URL)
	}

	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      site,
	}

	params := url.Values{}
	params.Add("action", "msv-get-sources")
	params.Add("id", matchedSourceReq[1])
	params.Add("nonce", matchedSourceReq[2])

	res, err := request.Request(http.MethodPost, api, headers, strings.NewReader(params.Encode()))
	if err != nil {
		return static.Data{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return static.Data{}, err
	}

	sources := &streams{}
	err = json.Unmarshal(body, &sources)
	if err != nil {
		return static.Data{}, err
	}
	if !sources.Success {
		return static.Data{}, fmt.Errorf("the api request for the streams did not return successful for %s", URL)
	}

	streams := make(map[string]*static.Stream)
	for _, source := range sources.Sources {
		headers, err := request.Headers(source, source)
		if err != nil {
			return static.Data{}, err
		}

		size, err := request.GetSizeFromHeaders(&headers)
		if err != nil {
			return static.Data{}, err
		}

		streams[fmt.Sprint(len(streams))] = &static.Stream{
			URLs: []*static.URL{
				{
					URL: source,
					Ext: strings.Split(headers.Get("content-type"), "/")[1],
				},
			},
			Size: size,
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
