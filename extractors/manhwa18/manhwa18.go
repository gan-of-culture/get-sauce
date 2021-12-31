package manhwa18

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type playerData struct {
	Data struct {
		Status  bool   `json:"status"`
		Sources string `json:"sources"`
	} `json:"data"`
}

const site = "https://manhwa18.tv/"
const playerAPITemplate = site + "wp-content/themes/halimmovies/player.php?episode_slug=full&post_id=%s&nonce=%s"
const cdn = "https://cdn.manhwa18.tv/"
const m3uTemplateURL = "https://cdn.manhwa18.tv/vid/%s/index.m3u8"

var reEpisodeURL = regexp.MustCompile(`https://manhwa18\.tv/\w+-\w+-[^"./]+`)
var reAPIParams = regexp.MustCompile(`postid-(\d+).+data-nonce="([^"]+)`) //1=postid 2=nonce
var reSourceURL = regexp.MustCompile(cdn + `hls/([^.]+).+`)               //1=id

type extractor struct{}

// New returns a manhwa18.tv extractor.
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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	// remove side bars with recommended content
	htmlString = strings.Split(htmlString, "</main>")[0]

	return reEpisodeURL.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := html.UnescapeString(utils.GetMeta(&htmlString, "og:title"))

	matchedAPIParams := reAPIParams.FindStringSubmatch(htmlString)
	if len(matchedAPIParams) < 2 {
		return nil, static.ErrDataSourceParseFailed
	}

	jsonData, err := request.GetAsBytesWithHeaders(fmt.Sprintf(playerAPITemplate, matchedAPIParams[1], matchedAPIParams[2]), map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	})
	if err != nil {
		return nil, err
	}

	pData := playerData{}
	err = json.Unmarshal(jsonData, &pData)
	if err != nil {
		return nil, err
	}
	if !pData.Data.Status {
		return nil, static.ErrDataSourceParseFailed
	}

	pData.Data.Sources = strings.ReplaceAll(pData.Data.Sources, `\`, "")
	matchedSrcURL := reSourceURL.FindStringSubmatch(pData.Data.Sources)
	if len(matchedSrcURL) < 2 {
		return nil, static.ErrDataSourceParseFailed
	}

	srcURL := fmt.Sprintf(m3uTemplateURL, matchedSrcURL[1])

	m3uMedia, err := request.Get(srcURL)
	if err != nil {
		return nil, err
	}

	URLs, _, err := request.ParseHLSMediaStream(&m3uMedia, srcURL)
	if err != nil {
		return nil, err
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: URLs,
				Ext:  "mp4",
			},
		},
		URL: URL,
	}, nil
}
