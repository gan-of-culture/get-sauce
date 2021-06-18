package hentaistream

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type quality struct {
	codec string
	vType string
}

const site = "https://hentaistream.moe/"
const CDN = "https://01cdn.hentaistream.moe/"

var players = map[string][]quality{
	"player.html": {
		{
			codec: "av1.1080p.webm",
			vType: "video/webm",
		},
		{
			codec: "av1.720p.webm",
			vType: "video/webm",
		},
		{
			codec: "vp9.720p.webm",
			vType: "video/webm",
		},
		{
			codec: "x264.720p.mp4",
			vType: "video/mp4",
		},
	},
	"player4k.html": {
		{
			codec: "av1.2160p.webm",
			vType: "video/webm",
		},
		{
			codec: "av1.1080p.webm",
			vType: "video/webm",
		},
		{
			codec: "av1.720p.webm",
			vType: "video/webm",
		},
		{
			codec: "vp9.720p.webm",
			vType: "video/webm",
		},
		{
			codec: "x264.720p.mp4",
			vType: "video/mp4",
		},
	},
}

type extractor struct{}

// New returns a hentaistream extractor.
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
	if ok, _ := regexp.MatchString(`https://hentaistream.moe/\d*/`, URL); ok {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://hentaistream.moe/anime/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://hentaistream.moe/\d*/[^"]*`)
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	if strings.Contains(htmlString, "<title>DDOS-GUARD</title>") {
		time.Sleep(200 * time.Millisecond)
		htmlString, _ = request.Get(URL)
	}

	re := regexp.MustCompile(`<iframe[\s\S]*?(player[^#]*)#([^"]*)`)
	matchedBase64CDNURL := re.FindStringSubmatch(htmlString) // 1=player[4k].html  2 = "url=https://01cdn.hentaistream.moe/2021/02/Overflow/E08/"
	if len(matchedBase64CDNURL) < 2 {
		return static.Data{}, static.ErrDataSourceParseFailed
	}

	downloadURLBytes, err := base64.StdEncoding.DecodeString(matchedBase64CDNURL[2])
	if err != nil {
		return static.Data{}, err
	}
	baseDownloadURL := strings.Split(strings.TrimPrefix(strings.Trim(string(downloadURLBytes), `"`), "url="), ";")[0]

	streams := make(map[string]*static.Stream)
	for i, quality := range players[matchedBase64CDNURL[1]] {
		size, err := request.Size(fmt.Sprintf("%s%s", baseDownloadURL, quality.codec), site)
		if err != nil {
			return static.Data{}, err
		}

		streams[fmt.Sprint(i)] = &static.Stream{
			URLs: []*static.URL{
				{
					URL: fmt.Sprintf("%s%s", baseDownloadURL, quality.codec),
					Ext: strings.Split(quality.vType, "/")[1],
				},
			},
			Quality: quality.codec,
			Size:    size,
		}
	}

	return static.Data{
		Site:    site,
		Title:   utils.GetH1(&htmlString, -1),
		Type:    "video",
		Streams: streams,
		Url:     URL,
	}, nil

}
