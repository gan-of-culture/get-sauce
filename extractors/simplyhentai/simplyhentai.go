package simplyhentai

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

/*
yet another hentai/porn site group with multiple sites same content different page layout
*/

type size struct {
	Full string `json:"full"`
}

type page struct {
	PageNum uint `json:"page_num"`
	Sizes   size `json:"sizes"`
}

type language struct {
	Name string `json:"name"`
}

type data struct {
	ImageCount uint     `json:"image_count"`
	Language   language `json:"language"`
	Albums     []album  `json:"albums"`
	Images     []page   `json:"images"`
	Pages      []page   `json:"pages"`
	Title      string   `json:"title"`
}

type series struct {
	Slug string `json:"slug"`
}

type album struct {
	Series series `json:"series"`
	Slug   string `json:"slug"`
}

type objects struct {
	Albums []album `json:"albums"`
}

type initData struct {
	IsLoading bool    `json:"isLoading"`
	Objects   objects `json:"objects"`
	Data      data    `json:"data"`
}

type appState struct {
	InitialData initData `json:"initialData"`
}

var site string

var reAppState *regexp.Regexp = regexp.MustCompile(`__SERVER_APP_STATE__ =  ({[^<]+)`)
var reExt *regexp.Regexp = regexp.MustCompile(`\w+$`)

type extractor struct{}

// New returns a simply-hentai extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, static.ErrURLParseFailed
	}

	site = "https://" + u.Host + "/"

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
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	matchedAppState := reAppState.FindStringSubmatch(htmlString)
	if len(matchedAppState) < 2 {
		return nil
	}

	appStat := appState{}
	err = json.Unmarshal([]byte(matchedAppState[1]), &appStat)
	if err != nil {
		return nil
	}

	if appStat.InitialData.Data.ImageCount != 0 {
		return []string{URL}
	}

	if appStat.InitialData.Objects.Albums == nil {
		appStat.InitialData.Objects.Albums = appStat.InitialData.Data.Albums
	}

	out := []string{}
	for _, a := range appStat.InitialData.Objects.Albums {
		out = append(out, fmt.Sprintf("%s%s/%s", site, a.Series.Slug, a.Slug))
	}

	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return &static.Data{}, err
	}

	matchedAppState := reAppState.FindStringSubmatch(htmlString)
	if len(matchedAppState) < 2 {
		return &static.Data{}, errors.New("app state not found for URL")
	}

	appStat := appState{}
	err = json.Unmarshal([]byte(matchedAppState[1]), &appStat)
	if err != nil {
		return &static.Data{}, err
	}

	if appStat.InitialData.Data.ImageCount == 0 {
		return &static.Data{}, errors.New("no images found for URL")
	}

	images := appStat.InitialData.Data.Pages
	if images == nil {
		images = appStat.InitialData.Data.Images
	}

	pages := utils.NeedDownloadList(len(images))

	URLs := []*static.URL{}
	for _, p := range pages {
		ext := reExt.FindString(images[p-1].Sizes.Full)
		URLs = append(URLs, &static.URL{
			URL: images[p-1].Sizes.Full,
			Ext: ext,
		})
	}

	return &static.Data{
		Site:  site,
		Title: appStat.InitialData.Data.Title,
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: URLs,
				Info: appStat.InitialData.Data.Language.Name,
			},
		},
		Url: URL,
	}, nil
}
