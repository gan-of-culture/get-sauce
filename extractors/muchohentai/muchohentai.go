package muchohentai

/*
	Muchohentai uses a really sensitve cloudflare configuration.
	To prevent blocking I've added the tls config of https://github.com/DaRealFreak/cloudflare-bp-go to the default client of request.go
	&tls.Config{
		PreferServerCipherSuites: false,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
	}
	Although this also won't work 100% since sometimes they increase cloudflare protection to the maximum
*/

import (
	"html"
	"math/rand"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/parsers/hls"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://muchohentai.com/"

var reEpisodeURL = regexp.MustCompile(site + `aBo4Rk/\d+/`)
var reServerPrefixes = regexp.MustCompile(`va\d\d`)
var reSrcURLBaseParts = regexp.MustCompile(`(https://)"[^"]+"([^"]+)`) // 1=https 2=second part of domain
var reSrcURLParts = regexp.MustCompile(`\\/wp-content[^"]+`)           // 1=m3u8 2=subs 3=thumbs
var reCaptionLangu = regexp.MustCompile(`subName="([^"]+)`)            // 1=Langu

type extractor struct{}

// New returns a muchohentai extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract data from URL
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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	return reEpisodeURL.FindAllString(htmlString, -1)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := html.UnescapeString(utils.GetH1(&htmlString, -1))
	captionLangu := utils.GetLastItemString(reCaptionLangu.FindStringSubmatch(htmlString))

	servers := reServerPrefixes.FindAllString(htmlString, -1) //va01 doesn't work - might also be a honey pot
	if len(servers) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}
	servers = servers[1:]

	baseURLParts := reSrcURLBaseParts.FindStringSubmatch(htmlString)
	if len(baseURLParts) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	baseURL := baseURLParts[1] + servers[rand.Intn(len(servers))] + baseURLParts[2]

	srcURLParts := reSrcURLParts.FindAllString(htmlString, -1)
	if len(srcURLParts) < 3 {
		return nil, static.ErrDataSourceParseFailed
	}

	for idx, v := range srcURLParts {
		srcURLParts[idx] = strings.ReplaceAll(v, `\`, ``)
	}

	m3uMasterURL := baseURL + srcURLParts[0]

	streams, err := hls.Extract(m3uMasterURL, map[string]string{"Referer": site})
	if err != nil {
		return nil, err
	}

	var ext string
	for _, stream := range streams {
		ext = stream.URLs[0].Ext

		if strings.Contains(stream.Info, "mp4a") {
			ext = "mp4"
		}

		stream.Ext = ext
	}

	return &static.Data{
		Site:    site,
		Title:   title,
		Type:    static.DataTypeVideo,
		Streams: streams,
		Captions: []*static.Caption{
			{
				URL: static.URL{
					URL: baseURL + srcURLParts[1],
					Ext: utils.GetFileExt(srcURLParts[1]),
				},
				Language: captionLangu,
			},
		},
		URL: URL,
	}, nil
}
