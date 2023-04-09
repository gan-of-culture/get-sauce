package hentaimama

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/parsers/hls"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type source struct {
	URL     string
	Referer string
}

const site = "https://hentaimama.io/"
const api = "https://hentaimama.io/wp-admin/admin-ajax.php"

var reMirrorURLs = regexp.MustCompile(`[^"]*new\d.php\?p=([^"]*)`)
var reExt = regexp.MustCompile(`([a-z][\w]*)(?:\?|$)`)
var reMimeType = regexp.MustCompile(`video/[^']*`)
var rePostID = regexp.MustCompile(`a:'(\d+)'`)

type extractor struct{}

// New returns a hentaimama extractor
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

	matchedMirrorURLs, err := getMirrorURLs(&episodeHtmlString, URL)
	if err != nil {
		return nil, err
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

		streams, err = hls.ExtractHLS(src.URL, map[string]string{"Referer": src.Referer})
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

func getMirrorURLs(htmlString *string, URL string) ([][]string, error) {
	matchedID := rePostID.FindStringSubmatch(*htmlString)
	if len(matchedID) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	params := url.Values{}
	params.Add("action", "get_player_contents")
	params.Add("a", matchedID[1])

	res, err := request.Request(http.MethodPost, api, map[string]string{
		"Referer":      URL,
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
	}, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, errors.New("api request failed")
	}
	defer res.Body.Close()

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resString := string(buffer)
	resString = strings.ReplaceAll(resString, `\`, "")

	return reMirrorURLs.FindAllStringSubmatch(resString, -1), nil
}
