package hentaimama

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type source struct {
	URL     string
	Referer string
}

const site = "https://hentaimama.io/"

var reMirrorURLs = regexp.MustCompile(`[^"]*new\d.php\?p=([^"]*)`)
var reExt = regexp.MustCompile(`([a-z][\w]*)(?:\?|$)`)
var reMimeType = regexp.MustCompile(`video/[^']*`)

type extractor struct{}

// New returns a hentaimama extractor.
func New() static.Extractor {
	return &extractor{}
}

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

func extractData(URL string) (*static.Data, error) {
	episodeHtmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedMirrorURLs := reMirrorURLs.FindAllStringSubmatch(episodeHtmlString, -1)
	if len(matchedMirrorURLs) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	sources := make([]source, len(matchedMirrorURLs))
	for i, u := range matchedMirrorURLs {
		b64Path, err := base64.StdEncoding.DecodeString(u[1])
		if err != nil {
			return nil, err
		}
		b64Paths := strings.Split(string(b64Path), "?")

		htmlString, err := request.Get(u[0])
		if err != nil {
			return nil, err
		}

		reSrc := regexp.MustCompile(fmt.Sprintf(`[^"']*/%s[^"']*`, string(b64Paths[0])))
		sources[i] = source{
			URL:     reSrc.FindString(htmlString),
			Referer: u[0],
		}
	}

	mirrorIdx := 0
	streams := map[string]*static.Stream{}
	// resolve all HLS URLs
	for _, src := range sources {
		ext := strings.TrimSuffix(utils.GetLastItemString(reExt.FindStringSubmatch(src.URL)), "?")
		if ext != "m3u8" {
			continue
		}

		streams, err = request.ExtractHLS(src.URL, map[string]string{"Referer": src.Referer})
		if err != nil {
			return nil, err
		}

		mirrorIdx += 1
		for _, v := range streams {
			v.Ext = "mp4"
			v.Info = fmt.Sprintf("Mirror %d", mirrorIdx)
		}
	}

	idx := len(streams) - 1
	// resolve other URLs
	for _, src := range sources {
		ext := strings.TrimSuffix(utils.GetLastItemString(reExt.FindStringSubmatch(src.URL)), "?")
		if ext == "m3u8" {
			continue
		}

		size, err := request.Size(src.URL, site)
		if err != nil {
			return nil, err
		}

		if ext == "" {
			ext = strings.Split(reMimeType.FindString(src.URL), "/")[1]
		}

		idx += 1
		mirrorIdx += 1
		streams[fmt.Sprint(idx)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: src.URL,
					Ext: ext,
				},
			},
			Size: size,
			Info: fmt.Sprintf("Mirror %d", mirrorIdx),
		}
		continue
	}

	return &static.Data{
		Site:    site,
		Title:   utils.GetMeta(&episodeHtmlString, "og:title"),
		Type:    "video",
		Streams: streams,
		URL:     URL,
	}, nil

}
