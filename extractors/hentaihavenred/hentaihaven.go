package hentaihavenred

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type embed struct {
	EmbedURL string `json:"embed_url"`
	Type     string `json:"type"`
}

const site = "https://hentaihaven.red/"
const api = "https://hentaihaven.red/wp-admin/admin-ajax.php"

var reEpisodeURL = regexp.MustCompile(site + `hentai/[\w-%]+/`)
var reParseURLShow = regexp.MustCompile(site + `watch/[\w-%]+/`)
var rePostID = regexp.MustCompile(site + `\?p=(\d+)`)
var rePlayer = regexp.MustCompile(`https://htstreaming.com/(?:(?:player/index.php\?data=[^\\"]+)|(?:video/[^\\"]+))`)

type extractor struct{}

// New returns a hentaihavenred extractor.
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
	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if strings.Contains(URL, "/watch/") {
		htmlString = strings.Split(htmlString, `<div class="bixbox"`)[0]
		return utils.RemoveAdjDuplicates(reEpisodeURL.FindAllString(htmlString, -1))
	}

	// contains list of show that need to be derefenced to episode level
	htmlString = strings.Split(htmlString, `<div id="sidebar">`)[0]

	out := []string{}
	for _, anime := range reParseURLShow.FindAllString(htmlString, -1) {
		out = append(out, parseURL(anime)...)
	}
	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := strings.Split(utils.GetMeta(&htmlString, "og:title"), " - ")[0]
	title = strings.Split(title, " | ")[0]

	params := url.Values{}
	params.Add("action", "doo_player_ajax")
	params.Add("post", utils.GetLastItemString(rePostID.FindStringSubmatch(htmlString)))
	params.Add("nume", "1")
	params.Add("type", "movie")

	res, err := request.Request(http.MethodPost, api, map[string]string{
		"Referer": URL,
		//		"Content-Length": len(params.Encode()),
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

	embedData := embed{}
	err = json.Unmarshal(buffer, &embedData)
	if err != nil {
		return nil, err
	}

	data, err := htstreaming.ExtractData(rePlayer.FindString(embedData.EmbedURL))
	if err != nil {
		return nil, err
	}

	data.Site = site
	data.Title = title
	data.URL = URL
	return data, nil
}
