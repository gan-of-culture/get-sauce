package hentai2read

import (
	"encoding/json"
	"html"
	"regexp"
	"slices"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type gData struct {
	MangaID      int
	Title        string
	Index        int
	Images       []string
	PreloadLimit int
	MainURL      string
}

const site = "https://hentai2read.com/"
const cdn = "https://static.hentaicdn.com/hentai"

var reJSONString = regexp.MustCompile(`{\s*'mangaID'[\s\S]*?}`)
var reTitle = regexp.MustCompile(`[^[(|]*`)

type extractor struct{}

// New returns a hentai2read extractor.
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
	URL = strings.Split(URL, "#")[0]
	if strings.Contains(URL, "_") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`([^/]*/)" class="title"`)

	URLs := []string{}
	for _, u := range re.FindAllStringSubmatch(htmlString, -1) {
		URLs = append(URLs, site+u[1])
	}
	return URLs
}

func extractData(URL string) (*static.Data, error) {
	URLs, err := fetchChapterURLs(URL)
	if err != nil {
		return nil, err
	}

	var imagesParts []string
	var title string
	for _, u := range URLs {
		htmlString, err := request.Get(u)
		if err != nil {
			return nil, err
		}

		jsonString := strings.ReplaceAll(reJSONString.FindString(htmlString), "'", `"`)

		galleryData := gData{}
		err = json.Unmarshal([]byte(jsonString), &galleryData)
		if err != nil {
			return nil, err
		}
		imagesParts = append(imagesParts, galleryData.Images...)

		if title == "" {
			title = html.UnescapeString(strings.TrimSpace(reTitle.FindString(galleryData.Title)))
		}
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: buildFullImgURL(imagesParts),
			},
		},
		URL: URL,
	}, nil
}

func buildFullImgURL(URIParts []string) []*static.URL {
	out := []*static.URL{}

	for _, idx := range utils.NeedDownloadList(len(URIParts)) {
		URIPart := URIParts[idx]
		out = append(out, &static.URL{
			URL: cdn + URIPart,
			Ext: utils.GetLastItemString(strings.Split(URIPart, ".")),
		})
	}
	return out
}

func fetchChapterURLs(URL string) ([]string, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	chapterURLs := regexp.MustCompile(URL+`[\d\.]+/`).FindAllString(htmlString, -1)
	fallbackURL := URL + "1/"
	// the chapter URLs come sorted, but are preceeded by a couple of duplicates -> remove them
	chapterURLs = slices.DeleteFunc(chapterURLs, func(s string) bool { return s == fallbackURL })
	chapterURLs = append(chapterURLs, fallbackURL)
	// chapters come sorted newest -> oldest so we reverse it
	slices.Reverse(chapterURLs)

	return chapterURLs, nil
}
