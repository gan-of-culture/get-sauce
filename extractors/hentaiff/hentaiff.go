package hentaiff

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/v2/request"
	"github.com/gan-of-culture/get-sauce/v2/static"
	"github.com/gan-of-culture/get-sauce/v2/utils"
)

type videoData struct {
	Success bool `json:"success"`
	Data    []struct {
		File  string `json:"file"`
		Label string `json:"label"`
		Type  string `json:"type"`
	} `json:"data"`
	Captions []interface{} `json:"captions"`
	IsVr     bool          `json:"is_vr"`
}

var reParseURLMass = regexp.MustCompile(`https://hentaiff.com/(anime|genres|raw|sub|uncensored|censored|bookmark|studio|director)`)
var reParseURLShow = regexp.MustCompile(`https://hentaiff.com/anime/[\w-%]+/`)
var reParseURL = regexp.MustCompile(`https://hentaiff.com/[\w-%]+(?:raw|english-subbed|english-dubbed|previews)/`)

var reParseAmHentaiID = regexp.MustCompile(`https://amhentai.com/v/([^"]+)`)

const site = "https://hentaiff.com/"
const videoInfoAPI = "https://amhentai.com/api/source/"

type extractor struct{}

// New returns a hentaiff.com extractor.
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
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	matchedGroup := reParseURLMass.FindStringSubmatch(URL)
	if len(matchedGroup) < 2 {
		if !reParseURL.MatchString(URL) {
			return []string{}
		}
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	switch matchedGroup[1] {
	case "anime":
		htmlString = strings.Split(htmlString, `<div class="bixbox"`)[0]
		return utils.RemoveAdjDuplicates(reParseURL.FindAllString(htmlString, -1))
	default:
		// contains list of show that need to be derefenced to episode level
		htmlString = strings.Split(htmlString, `<div id="sidebar">`)[0]

		out := []string{}
		for _, anime := range reParseURLShow.FindAllString(htmlString, -1) {
			out = append(out, parseURL(anime)...)
		}
		return out
	}
}

func extractData(URL string) (static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := html.UnescapeString(utils.GetH1(&htmlString, -1))

	matchedVideoID := reParseAmHentaiID.FindStringSubmatch(htmlString) //1=VID
	if len(matchedVideoID) < 2 {
		return static.Data{}, static.ErrDataSourceParseFailed
	}

	jsonData, err := request.PostAsBytesWithHeaders(videoInfoAPI+matchedVideoID[1], map[string]string{"Referer": matchedVideoID[0]})
	if err != nil {
		return static.Data{}, err
	}

	videoInfo := videoData{}
	err = json.Unmarshal(jsonData, &videoInfo)
	if err != nil {
		return static.Data{}, err
	}

	streams := map[string]*static.Stream{}
	dataLen := len(videoInfo.Data)
	for i, stream := range videoInfo.Data {
		url, size, err := resolveVideoSource(stream.File)
		if err != nil {
			return static.Data{}, err
		}

		streams[fmt.Sprint(dataLen-i-1)] = &static.Stream{
			URLs: []*static.URL{
				{
					URL: url,
					Ext: stream.Type,
				},
			},
			Quality: stream.Label,
			Size:    size,
			Info:    "Downloading concurrently with the -w is not possible for this source",
		}
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    static.DataTypeVideo,
		Streams: streams,
		Url:     URL,
	}, nil
}

func resolveVideoSource(URL string) (string, int64, error) {
	res, err := request.Request(http.MethodHead, URL, nil, nil)
	if err != nil {
		return "", 0, nil
	}

	size, _ := strconv.ParseInt(res.Header.Get("Content-Length"), 0, 64)

	return res.Request.URL.String(), size, nil
}
