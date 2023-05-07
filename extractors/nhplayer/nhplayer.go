package nhplayer

import (
	"encoding/base64"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reNHPlayerURL = regexp.MustCompile(`https://nhplayer\.com/v/[^/"]+`)
var rePlayerURL = regexp.MustCompile(`/player.php\?u=([^"&]+)`)
var reHTStreamingVideoURL = regexp.MustCompile(`https://htstreaming.com/video/([^"]*)`)

type extractor struct{}

// New returns a nhplayer extractor
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
	if !reNHPlayerURL.MatchString(URL) {
		return nil
	}
	return []string{URL}
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	videoURL := reHTStreamingVideoURL.FindString(htmlString)
	if videoURL != "" {
		return htstreaming.ExtractData(videoURL)
	}

	playerURL := rePlayerURL.FindStringSubmatch(htmlString)
	if len(playerURL) < 2 {
		return nil, static.ErrDataSourceParseFailed
	}

	b64Path, err := base64.StdEncoding.DecodeString(playerURL[1])
	if err != nil {
		return nil, err
	}
	videoURL = string(b64Path)

	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetLastItemString(strings.Split(videoURL, "/"))
	title = strings.Split(title, ".")[0]

	size, _ := request.Size(videoURL, URL)

	return &static.Data{
		Site:  baseURL.Host,
		Title: title,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					{
						URL: videoURL,
						Ext: utils.GetFileExt(videoURL),
					},
				},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}
