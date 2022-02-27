package hentaihd

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"

	"github.com/gan-of-culture/get-sauce/extractors/animestream"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type videoInfo struct {
	Success bool `json:"success"`
	Data    []struct {
		File  string `json:"file"`
		Label string `json:"label"`
		Type  string `json:"type"`
	} `json:"data"`
	Captions []interface{} `json:"captions"`
	IsVr     bool          `json:"is_vr"`
}

const site = "https://v2.hentaihd.net/"
const videoProvider = "https://amhentai.com/"
const videoProviderAPI = videoProvider + "api/source/"

var reEpisodeURL = regexp.MustCompile(site + `\d+/.+/`)
var reParseURLShow = regexp.MustCompile(site + `anime/[\w-%]+/`)
var reVideoURL = regexp.MustCompile(videoProvider + `v/([^"]+)`)

type extractor struct{}

// New returns a hentaihd extractor.
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
	return animestream.ParseURL(URL, site)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedVideo := reVideoURL.FindStringSubmatch(htmlString) // 0=URL 1=ID
	if len(matchedVideo) < 2 {
		return nil, static.ErrDataSourceParseFailed
	}

	videoJSON, err := request.PostAsBytesWithHeaders(videoProviderAPI+matchedVideo[1], map[string]string{"Referer": matchedVideo[0]})
	if err != nil {
		return nil, err
	}

	video := videoInfo{}
	err = json.Unmarshal(videoJSON, &video)
	if err != nil {
		return nil, err
	}

	streams := map[string]*static.Stream{}
	for i, s := range video.Data {
		size, _ := request.Size(s.File, site)
		streams[fmt.Sprint(len(video.Data)-i-1)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: s.File,
					Ext: s.Type,
				},
			},
			Quality: s.Label,
			Size:    size,
		}
	}

	return &static.Data{
		Site:    site,
		Title:   html.UnescapeString(utils.GetH1(&htmlString, -1)),
		Type:    static.DataTypeVideo,
		Streams: streams,
		URL:     URL,
	}, nil
}
