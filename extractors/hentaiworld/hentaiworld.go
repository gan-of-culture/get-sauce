package hentaiworld

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://hentaiworld.tv/"

type extractor struct{}

// New returns a hentaiworld extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract data of provided url
func (e *extractor) Extract(url string) ([]*static.Data, error) {
	urls, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	data := []*static.Data{}
	for _, u := range urls {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}
	return data, nil
}

// parseURL for data extraction
func parseURL(url string) ([]string, error) {
	re := regexp.MustCompile(`(?:https://hentaiworld.tv/)(?:all-episodes|uncensored|3d|hentai-videos/category|hentai-videos/tag)/`)
	validMassURL := re.FindString(url)
	if validMassURL == "" {
		re := regexp.MustCompile(`hentai-videos/(?:3d/)?(?:.+episode-[0-9]*)?`)
		validEpisodeURL := re.FindString(url)
		if validEpisodeURL != "" {
			return []string{url}, nil
		}
		return []string{}, fmt.Errorf("[HentaiWorld]Invalid URL %s", url)
	}

	massHTMLPage, err := request.Get(url)
	if err != nil {
		return []string{}, fmt.Errorf("[HentaiWorld]HTTP GET URL  error %v", err)
	}

	re = regexp.MustCompile(`"display-all-posts-background"><a href="([^"]*)`)
	matchedEpisodesURLs := re.FindAllStringSubmatch(massHTMLPage, -1)
	urls := []string{}
	for _, matchedURL := range matchedEpisodesURLs {
		urls = append(urls, matchedURL[1])
	}

	return urls, nil
}

//extractData of hentai
func extractData(url string) (static.Data, error) {
	postHTMLpage, err := request.Get(url)
	if err != nil {
		return static.Data{}, nil
	}

	title := strings.TrimSuffix(utils.GetMeta(&postHTMLpage, "og:title"), " - HentaiWorld")

	if strings.Contains(title, "\u0026#8211;") {
		title = strings.ReplaceAll(title, "\u0026#8211;", "-")
	}

	re := regexp.MustCompile(`window.open\(\'([^']+\.([0-9a-zA-z]*))`)
	infoAboutFile := re.FindStringSubmatch(postHTMLpage) // 1 = dlURL 2=ext

	if len(infoAboutFile) != 3 {
		re = regexp.MustCompile(`src='(.*)\.(mp4*).*`)
		infoAboutFile = re.FindStringSubmatch(postHTMLpage) // 1 = dlURL 2=ext
		if len(infoAboutFile) != 3 {
			return static.Data{}, fmt.Errorf("[HentaiWorld] Get scrape video info for URL %s", url)
		}
	}
	infoAboutFile[1] = strings.ReplaceAll(infoAboutFile[1], " ", "%20")
	size, _ := request.Size(infoAboutFile[1], site)

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "video",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					0: {
						URL: infoAboutFile[1],
						Ext: infoAboutFile[2],
					},
				},
				Quality: "best",
				Size:    size,
			},
		},
		Url: url,
	}, nil
}
