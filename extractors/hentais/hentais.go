package hentais

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://www.hentais.tube/"

var rePlayerURL = regexp.MustCompile(`player.php[^']*`)
var reSourceTag = regexp.MustCompile(`src="([^"]*)" type="([^"]*)"(?: label="([^"]*)")?`) // 1=videoURL 2=mimeType 3=quality

type extractor struct{}

// New returns a hentais extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract hentai data
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

// parseURL to extract hentai data
func parseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"episodes/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, site+"tvshows/") {
		return nil
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	re := regexp.MustCompile(`/episodes/[^"]*`)
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetH1(&htmlString, -1)

	playerURL := site + rePlayerURL.FindString(htmlString)
	if playerURL == "" {
		return nil, errors.New("can't parse playerURL for")
	}

	htmlString, err = request.Get(playerURL)
	if err != nil {
		log.Println(URL)
		return nil, err
	}

	matchedSrcTag := reSourceTag.FindAllStringSubmatch(htmlString, -1) //<-- is basically the different streams
	if len(matchedSrcTag) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	u := ""
	quality := ""
	mimeType := ""
	streams := map[string]*static.Stream{}
	for i, srcTag := range matchedSrcTag {
		quality = ""
		mimeType = ""

		u = srcTag[1]
		if !strings.Contains(srcTag[1], "http") {
			u = site + srcTag[1][1:] //remove extra slash
		}

		switch len(srcTag) {
		case 3:
			mimeType = srcTag[2]
		case 4:
			mimeType = srcTag[2]
			quality = srcTag[3]
		}
		size, _ := request.Size(u, playerURL)
		streams[fmt.Sprintf("%d", len(matchedSrcTag)-i-1)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: u,
					Ext: utils.GetLastItemString(strings.Split(mimeType, "/")),
				},
			},
			Quality: quality,
			Size:    size,
		}
	}
	return &static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: streams,
		URL:     URL,
	}, nil
}
