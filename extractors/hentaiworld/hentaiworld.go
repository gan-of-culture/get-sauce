package hentaiworld

import (
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
	URLs := parseURL(url)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

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

// parseURL for data extraction
func parseURL(URL string) []string {
	re := regexp.MustCompile(`(?:https://hentaiworld.tv/)(?:all-episodes|uncensored|3d|hentai-videos/category|hentai-videos/tag)/`)
	validMassURL := re.FindString(URL)
	if validMassURL == "" {
		re := regexp.MustCompile(`hentai-videos/(?:3d/)?(?:.+episode-[0-9]*)?`)
		validEpisodeURL := re.FindString(URL)
		if validEpisodeURL != "" {
			return []string{URL}
		}
		return []string{}
	}

	massHTMLPage, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re = regexp.MustCompile(`"display-all-posts-background"><a href="([^"]*)`)
	matchedEpisodesURLs := re.FindAllStringSubmatch(massHTMLPage, -1)
	urls := []string{}
	for _, matchedURL := range matchedEpisodesURLs {
		urls = append(urls, matchedURL[1])
	}

	return urls
}

//extractData of hentai
func extractData(URL string) (static.Data, error) {
	postHTMLpage, err := request.Get(URL)
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
			return static.Data{}, static.ErrDataSourceParseFailed
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
				Size: size,
			},
		},
		Url: URL,
	}, nil
}
