package hentais

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://www.hentais.tube/"
const infoSite = "https://www.hentais.tube/tvshows/"
const dlSite = "https://www.hentais.tube/download-hentai/"

// Extract hentai data
func Extract(url string) ([]static.Data, error) {
	url, err := ParseURL(url)
	if err != nil {
		return []static.Data{}, err
	}

	return ExtractData(url)
}

// ParseURL to extract hentai data
func ParseURL(url string) (string, error) {
	return url, nil
}

// ExtractData of hentai
func ExtractData(url string) ([]static.Data, error) {

	re := regexp.MustCompile("/([^/0-9]*)?([0-9])*/?$")
	seriesInfo := re.FindStringSubmatch(url)

	seriesName := seriesInfo[1]
	if strings.Contains(seriesName, "episode") {
		seriesName = strings.TrimSuffix(seriesName, "episode-")
	}

	seriesName = strings.TrimSuffix(seriesName, "-")
	wantedEpisodes := 0
	if seriesInfo[2] != "" {
		wantedEpisodes, _ = strconv.Atoi(seriesInfo[2])
	}

	htmlInfoPage, err := request.Get(fmt.Sprintf("%s%s", infoSite, seriesName))
	if err != nil {
		return []static.Data{}, err
	}

	re = regexp.MustCompile("<h1>([^<]*)")
	title := re.FindStringSubmatch(htmlInfoPage)

	re = regexp.MustCompile("n>?\\s*([0-9]*-[0-9]*-[0-9]*)")
	dates := re.FindAllStringSubmatch(htmlInfoPage, -1)

	/*re = regexp.MustCompile("\"/genre/[^\"]*\">?\\s*([^<]*)")
	tags := re.FindAllStringSubmatch(htmlInfoPage, -1)

	re = regexp.MustCompile("info1\"[^>]*>?\\s*<[^>]*>?\\s*<p>([^<]*)")
	descr := re.FindStringSubmatch(htmlInfoPage)

	re = regexp.MustCompile("/studio[^>]*>([^<]*)")
	studio := re.FindStringSubmatch(htmlInfoPage)

	re = regexp.MustCompile("Episodes[^0-9]*([0-9]*)")
	numberOfEpisodes := re.FindStringSubmatch(htmlInfoPage)*/

	re = regexp.MustCompile("/episodes/([^/\"]*)")
	matchedEpisodeURLs := re.FindAllStringSubmatch(htmlInfoPage, -1)

	episodeURLs := []string{}
	if wantedEpisodes != 0 {
		if wantedEpisodes-1 > len(matchedEpisodeURLs) {
			return []static.Data{}, fmt.Errorf("[Hentai] wanted episode %d not in availabel episodes %d", wantedEpisodes, len(episodeURLs))
		}
		episodeURLs = []string{matchedEpisodeURLs[wantedEpisodes-1][1]}
	} else {
		for _, url := range matchedEpisodeURLs {
			episodeURLs = append(episodeURLs, url[1])
		}
	}

	data := []static.Data{}
	for idx, episodeURL := range episodeURLs {
		downloadSite, err := request.Get(fmt.Sprintf("%s%s/", dlSite, episodeURL))
		if err != nil {
			return []static.Data{}, err
		}

		re = regexp.MustCompile("https://www\\.hentais\\.tube/baixar\\.php\\?enc=[^\"]*")
		redirectURL := re.FindString(downloadSite)

		headers, err := request.Headers(redirectURL, dlSite)
		if err != nil {
			return []static.Data{}, err
		}

		ext := utils.GetLastItemString(strings.Split(headers.Get("content-type"), "/"))
		if ext == "" {
			ext = "mp4"
		}
		size, _ := strconv.Atoi(headers.Get("content-length"))

		currentEp := wantedEpisodes
		if wantedEpisodes == 0 {
			currentEp = idx
		}
		data = append(data, static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s Episode %d", title[1], currentEp),
			Type:  "video",
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						0: {
							URL: redirectURL,
							Ext: ext,
						},
					},
					Quality: "unknown",
					Size:    int64(size),
					Info:    fmt.Sprintf("Date: %s", dates[idx][1]),
				},
			},
			Err: nil,
			Url: fmt.Sprintf("%s/tvshows/%s", site, seriesName),
		})
	}

	return data, nil

}
