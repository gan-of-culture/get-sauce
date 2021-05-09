package hentais

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://www.hentais.tube/"
const infoSite = "https://www.hentais.tube/tvshows/"
const dlSite = "https://www.hentais.tube/download-hentai/"

// Extract hentai data
func Extract(URL string) ([]static.Data, error) {
	URLs, err := ParseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

// ParseURL to extract hentai data
func ParseURL(URL string) ([]string, error) {
	if strings.HasPrefix(URL, site+"episodes/") {
		return []string{URL}, nil
	}

	if !strings.HasPrefix(URL, site+"tvshows/") {
		return nil, fmt.Errorf("[Hentais] Can't parse URL %s", URL)
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`/episodes/[^"]*`)
	return re.FindAllString(htmlString, -1), nil
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := utils.GetH1(htmlString)

	re := regexp.MustCompile(`player.php[^']*`)
	playerURL := site + re.FindString(htmlString)
	if playerURL == "" {
		return static.Data{}, fmt.Errorf("[Hentais] Can't parse playerURL for %s", URL)
	}

	htmlString, err = request.Get(playerURL)
	if err != nil {
		log.Println(URL)
		return static.Data{}, err
	}

	re = regexp.MustCompile(`src="([^"]*)" type="([^"]*)" label="([^"]*)"`) // 1=videoURL 2=mimeType 3=quality
	matchedSrcTag := re.FindAllStringSubmatch(htmlString, -1)               //<-- is basically the different streams
	if len(matchedSrcTag) < 1 {
		return static.Data{}, fmt.Errorf("[Hentais] No source tags found in %s", playerURL)
	}

	quality := ""
	mimeType := ""
	streams := map[string]static.Stream{}
	for i, srcTag := range matchedSrcTag {
		quality = ""
		mimeType = ""
		switch len(srcTag) {
		case 3:
			mimeType = srcTag[2]
		case 4:
			mimeType = srcTag[2]
			quality = srcTag[3]
		}
		size, _ := request.Size(srcTag[1], site)
		streams[fmt.Sprintf("%d", len(matchedSrcTag)-i-1)] = static.Stream{
			URLs: []static.URL{
				{
					URL: srcTag[1],
					Ext: utils.GetLastItemString(strings.Split(mimeType, "/")),
				},
			},
			Quality: quality,
			Size:    size,
		}
	}
	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: streams,
		Err:     nil,
		Url:     URL,
	}, nil
}
