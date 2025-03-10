package hentaiworld

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hentaiworld.tv/"
const videoProvider = "https://www.porn-d.xyz/TbLwA66UuPu4LiuOCsKr/"

var reFileInfo = regexp.MustCompile(`https://hentaiworld.tv/video-player.html\?(videos/[^.]+\.([^'"]+))`) // 1 = dlURLPart 2=ext

type extractor struct{}

// New returns a hentaiworld extractor
func New() static.Extractor {
	return &extractor{}
}

// Extract data of provided URL
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
	URLs := []string{}
	for _, matchedURL := range matchedEpisodesURLs {
		URLs = append(URLs, matchedURL[1])
	}

	return URLs
}

func extractData(URL string) (*static.Data, error) {
	postHTMLpage, err := request.Get(URL)
	if err != nil {
		return nil, nil
	}

	title := strings.TrimSuffix(utils.GetMeta(&postHTMLpage, "og:title"), " - HentaiWorld")

	if strings.Contains(title, "\u0026#8211;") {
		title = strings.ReplaceAll(title, "\u0026#8211;", "-")
	}

	infoAboutFile := reFileInfo.FindStringSubmatch(postHTMLpage)
	videoURL := videoProvider + infoAboutFile[1]
	size, _ := request.Size(videoURL, site)

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					0: {
						URL: videoURL,
						Ext: infoAboutFile[2],
					},
				},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}
