package damn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

const site = "https://damn.stream"
const embed = "https://www.damn.stream/video/"
const cdn = "https://server-one.damn.stream"

func ParseURL(URL string) []string {
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

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) < 1 {
		return nil, fmt.Errorf("[Damn] No matching URL found.")
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

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := strings.TrimPrefix(URL, "https://www.damn.stream/watch/hentai/")
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
		Type:  fmt.Sprintf("video/%s", srcMeta[2]),
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					0: {
						URL: srcMeta[1],
						Ext: srcMeta[2],
					},
				},
				Quality: "best",
				Size:    size,
			},
		},
		Url: URL,
	}, err
}
