package hentaimama

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/parsers/hls"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
)

type setup struct {
	Sources []struct {
		Type string `json:"type"`
		File string `json:"file"`
	} `json:"sources"`
}

type source struct {
	URL     string
	Referer string
}

type sourceParser interface {
	resolve(URL *url.URL) ([]*source, error)
}

type defaultParser struct{}

func (dp *defaultParser) resolve(URL *url.URL) ([]*source, error) {
	sources := []*source{}

	b64Path, err := base64.StdEncoding.DecodeString(URL.Query().Get("p"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	b64Paths := strings.Split(string(b64Path), "?")

	HTMLString, err := request.Get(URL.String())
	if err != nil {
		return nil, err
	}

	reSrc := regexp.MustCompile(fmt.Sprintf(`[^"']*/%s[^"']*`, string(b64Paths[0])))
	videoURL := reSrc.FindString(HTMLString)
	if videoURL == "" {
		log.Printf("skipping broken source: %s", URL.String())
		return nil, nil
	}
	sources = append(sources, &source{
		URL:     videoURL,
		Referer: URL.String(),
	})

	return sources, nil
}

func NewDefaultParser() sourceParser {
	return &defaultParser{}
}

type embedParser struct {
	reSetup *regexp.Regexp
}

func (ep *embedParser) resolve(URL *url.URL) ([]*source, error) {
	HTMLString, err := request.Get(URL.String())
	if err != nil {
		return nil, err
	}

	matchedSetup := ep.reSetup.FindStringSubmatch(HTMLString)
	if len(matchedSetup) == 0 || matchedSetup[1] == "" {
		return nil, errors.WithStack(fmt.Errorf("setup string with source info not found"))
	}

	jsonString := utils.GetJSONFromRelaxedJSObjStr(matchedSetup[1])

	setup := setup{}
	err = json.Unmarshal([]byte(jsonString), &setup)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sources := []*source{}
	for _, parsedSource := range setup.Sources {
		sources = append(sources, &source{
			URL:     parsedSource.File,
			Referer: URL.String(),
		})
	}

	return sources, nil
}

func NewEmbedParser() sourceParser {
	return &embedParser{reSetup: regexp.MustCompile(`setup\(([^\)]+)`)}
}

const site = "https://hentaimama.io/"
const api = "https://hentaimama.io/wp-admin/admin-ajax.php"

var rePlayerContentURLs = regexp.MustCompile(`https:\\\/\\\/hentaimama.io[^"]+`)
var reExt = regexp.MustCompile(`([a-z][\w]*)(?:\?|$)`)
var reMimeType = regexp.MustCompile(`video/[^']*`)
var rePostID = regexp.MustCompile(`a:\s'(\d+)'`)

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

	HTMLString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://hentaimama.io/episodes[^"]*`)
	return re.FindAllString(HTMLString, -1)[1:]
}

func extractData(URL string) (*static.Data, error) {
	episodeHTMLString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedPlayerContentURLs, err := getPlayerContentURLs(&episodeHTMLString, URL)
	if err != nil {
		return nil, err
	}

	sources := []*source{}
	for _, playerContentURL := range matchedPlayerContentURLs {
		playerContentURL = strings.ReplaceAll(html.UnescapeString(playerContentURL), `\`, "")
		u, err := url.Parse(playerContentURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		sourceParser := NewDefaultParser()
		if u.Query().Has("dt_embed") {
			sourceParser = NewEmbedParser()
		}

		s, err := sourceParser.resolve(u)
		if err != nil {
			return nil, err
		}
		sources = append(sources, s...)
	}

	mirrorIdx := 0
	streams := map[string]*static.Stream{}
	// resolve all HLS URLs
	for _, src := range sources {
		ext := strings.TrimSuffix(utils.GetLastItemString(reExt.FindStringSubmatch(src.URL)), "?")
		if ext != "m3u8" {
			continue
		}

		streams, err = hls.Extract(src.URL, map[string]string{"Referer": src.Referer})
		if err != nil {
			log.Println(err)
			log.Printf("skipping broken source: %s", src.URL)
			continue
		}

		mirrorIdx += 1
		for _, v := range streams {
			v.Ext = "mp4"
			v.Info = fmt.Sprintf("Mirror %d", mirrorIdx)
		}
	}

	idx := len(streams) - 1
	if idx == -1 {
		streams = make(map[string]*static.Stream)
	}
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
		//log.Println(streams)
		//log.Println(idx)
		streams[fmt.Sprint(idx)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: src.URL,
					Ext: ext,
				},
			},
			Size:    size,
			Info:    fmt.Sprintf("Mirror %d", mirrorIdx),
			Headers: map[string]string{},
		}
		continue
	}

	streams = utils.SortStreamsBySize(streams)

	return &static.Data{
		Site:    site,
		Title:   utils.GetH1(&episodeHTMLString, -1),
		Type:    "video",
		Streams: streams,
		URL:     URL,
	}, nil

}

func getPlayerContentURLs(HTMLString *string, URL string) ([]string, error) {
	matchedID := rePostID.FindStringSubmatch(*HTMLString)
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

	return rePlayerContentURLs.FindAllString(string(buffer), -1), nil
}
