package hentaimama

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://hentaimama.io/"

type extractor struct{}

// New returns a hentaimama extractor.
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
	if strings.HasPrefix(URL, "https://hentaimama.io/episodes") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://hentaimama.io/tvshows/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://hentaimama.io/episodes[^"]*`)
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {
	episodeHtmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`[^"]*new\d.php\?p=([^"]*)`)
	matchedMirrorURLs := re.FindAllStringSubmatch(episodeHtmlString, -1)
	if len(matchedMirrorURLs) < 1 {
		return static.Data{}, fmt.Errorf("can't locate video src URL for: %s", URL)
	}

	b64Path, err := base64.StdEncoding.DecodeString(matchedMirrorURLs[1][1])
	if err != nil {
		return static.Data{}, fmt.Errorf("error decoding string: %s ", err.Error())
	}

	streams := make(map[string]*static.Stream)
	reSrc := regexp.MustCompile(fmt.Sprintf("[^\"']*/%s[^\"']*", string(b64Path)))
	for i, u := range matchedMirrorURLs {
		htmlString, err := request.Get(u[0])
		if err != nil {
			return static.Data{}, err
		}
		srcURL := reSrc.FindString(htmlString)
		size, err := request.Size(srcURL, site)
		if err != nil {
			return static.Data{}, err
		}

		re = regexp.MustCompile(`\.([\d\w]*)\?`)
		ext := strings.TrimSuffix(re.FindStringSubmatch(srcURL)[1], "?")

		if ext == "" {
			re = regexp.MustCompile(`video/[^']*`)
			ext = strings.Split(re.FindString(srcURL), "/")[1]
		}

		streams[fmt.Sprint(i)] = &static.Stream{
			URLs: []*static.URL{
				{
					URL: srcURL,
					Ext: ext,
				},
			},
			Size: size,
			Info: fmt.Sprintf("Mirror %d", i+1),
		}

	}

	return static.Data{
		Site:    site,
		Title:   utils.GetMeta(&episodeHtmlString, "og:title"),
		Type:    "video",
		Streams: streams,
		Url:     URL,
	}, nil

}
