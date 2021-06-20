package damn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://www.damn.stream/"
const embed = "https://www.damn.stream/video/"

//const cdn = "https://server-one.damn.stream"

type extractor struct{}

// New returns a damn.stream extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) < 1 {
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
	if strings.HasPrefix(URL, "https://www.damn.stream/watch/hentai/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://www.damn.stream/hentai/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`[^"]*watch/hentai[^"]*`)
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := strings.TrimSuffix(utils.GetMeta(&htmlString, "og:title"), " - Damnstream")
	re := regexp.MustCompile(`"/video/([^"]*)`)
	videoID := re.FindStringSubmatch(htmlString)[1]

	htmlString, err = request.Get(fmt.Sprintf("%s%s", embed, videoID))
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`<source\s[^=]*="([^"]*\.([^"]*))"`)
	srcMeta := re.FindStringSubmatch(htmlString) //1=URL 2=ext

	srcMeta[1] = fmt.Sprintf("%s%s", "https:", srcMeta[1])

	size, _ := request.Size(srcMeta[1], site)

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "video",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					0: {
						URL: srcMeta[1],
						Ext: srcMeta[2],
					},
				},
				Size: size,
			},
		},
		Url: URL,
	}, err
}
