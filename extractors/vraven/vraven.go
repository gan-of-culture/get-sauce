package vraven

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/jwplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var site string
var reJWPlayerURL = regexp.MustCompile(`https://[^/]+.*/player[\./](?:php\?data=)?[^"\s]+`)

type extractor struct{}

// New returns a vraven extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, static.ErrURLParseFailed
	}

	site = fmt.Sprintf("https://%s/", u.Host)

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
	if ok, _ := regexp.MatchString(`/episode-\d*/?$`, URL); ok {
		return []string{URL}
	}

	if !strings.Contains(URL, site+"watch/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}
	slug := strings.Split(URL, "watch/")[1]
	re := regexp.MustCompile(fmt.Sprintf("[^\"]*%sepisode-\\d*", slug))
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(utils.GetH1(&htmlString, -1))

	data, err := jwplayer.New().Extract(reJWPlayerURL.FindString(htmlString))
	if err != nil {
		return nil, err
	}
	if len(data) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	video := data[0]
	video.Site = site
	video.Title = title
	video.Type = static.DataTypeVideo
	video.URL = URL

	return video, nil

}
