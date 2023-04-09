package simplyhentai

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type state struct {
	Props struct {
		PageProps struct {
			Status int `json:"status"`
			Tag    struct {
				Albums []struct {
					CommentCount int         `json:"comment_count"`
					Description  interface{} `json:"description"`
					ID           int         `json:"id"`
					ImageCount   int         `json:"image_count"`
					Language     struct {
						Name     string `json:"name"`
						Slug     string `json:"slug"`
						FlagCode string `json:"flag_code"`
					} `json:"language"`
					Series struct {
						AlbumCount   int    `json:"album_count"`
						CommentCount int    `json:"comment_count"`
						Slug         string `json:"slug"`
						Title        string `json:"title"`
						Type         string `json:"type"`
					} `json:"series"`
					Slug  string `json:"slug"`
					Title string `json:"title"`
					Type  string `json:"type"`
				}
			}
			Manga struct {
				Description interface{} `json:"description"`
				ID          int         `json:"id"`
				ImageCount  int         `json:"image_count"`
				Language    struct {
					Name     string `json:"name"`
					Slug     string `json:"slug"`
					FlagCode string `json:"flag_code"`
				} `json:"language"`
				Images []struct {
					ID      int `json:"id"`
					PageNum int `json:"page_num"`
					Sizes   struct {
						Full       string `json:"full"`
						SmallThumb string `json:"small_thumb"`
						Thumb      string `json:"thumb"`
						GiantThumb string `json:"giant_thumb"`
					} `json:"sizes"`
				} `json:"images"`
				Slug  string `json:"slug"`
				Title string `json:"title"`
			} `json:"manga"`
			Meta struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"meta"`
		} `json:"pageProps"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Series string `json:"series"`
		Slug   string `json:"slug"`
	} `json:"query"`
}

const site = "https://www.simply-hentai.com/"

var reAppState *regexp.Regexp = regexp.MustCompile(`__NEXT_DATA__.*?({[^<]+)`)

type extractor struct{}

// New returns a simply-hentai extractor
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
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	matchedAppState := reAppState.FindStringSubmatch(htmlString)
	if len(matchedAppState) < 2 {
		return nil
	}

	appStat := state{}
	err = json.Unmarshal([]byte(matchedAppState[1]), &appStat)
	if err != nil {
		return nil
	}

	if appStat.Props.PageProps.Manga.ImageCount != 0 {
		return []string{URL}
	}

	out := []string{}
	for _, a := range appStat.Props.PageProps.Tag.Albums {
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

	appStat := state{}
	err = json.Unmarshal([]byte(matchedAppState[1]), &appStat)
	if err != nil {
		return &static.Data{}, err
	}

	if appStat.Props.PageProps.Manga.ImageCount == 0 {
		return &static.Data{}, errors.New("no images found for URL")
	}

	images := appStat.Props.PageProps.Manga.Images

	pages := utils.NeedDownloadList(len(images))

	URLs := []*static.URL{}
	for _, p := range pages {
		ext := utils.GetFileExt(images[p-1].Sizes.Full)
		URLs = append(URLs, &static.URL{
			URL: images[p-1].Sizes.Full,
			Ext: ext,
		})
	}

	return &static.Data{
		Site:  site,
		Title: appStat.Props.PageProps.Manga.Title,
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
				Info: appStat.Props.PageProps.Manga.Language.Name,
			},
		},
		URL: URL,
	}, nil
}
