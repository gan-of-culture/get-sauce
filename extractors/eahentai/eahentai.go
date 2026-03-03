package eahentai

import (
	"encoding/json"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type extractor struct{}

type gallery struct {
	ImageID      int    `json:"imageID"`
	AddDt        string `json:"addDt"`
	ImageURI     string `json:"imageUri"`
	ThumbnailURI string `json:"thumbnailUri"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	AlbumID      int    `json:"albumID"`
	Views        int    `json:"views"`
	Sort         int    `json:"sort"`
	Album        any    `json:"album"`
}

const site = "https://eahentai.com"
const CDN = "https://i.eahentai.com/file/ea-gallery/"

var reGalleryURLPart = regexp.MustCompile(`"/(a/\d+)"`)
var reGallery = regexp.MustCompile(`"images\\":(\[{[^\]]+\])`)

func parseURL(URL string) ([]string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(u.Path, "/a/") {
		return []string{URL}, nil
	}

	body, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, URLPart := range reGalleryURLPart.FindAllStringSubmatch(body, -1) {
		uPart := utils.GetLastItemString(URLPart)
		u, err = url.Parse(site)
		if err != nil {
			return nil, err
		}
		u.Path = uPart
		out = append(out, u.String())
	}

	return out, nil
}

func extractData(URL string) (*static.Data, error) {
	body, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := strings.Split(html.UnescapeString(utils.GetMeta(&body, "og:title")), " - EAHentai")[0]

	matchedJSON := utils.GetLastItemString(reGallery.FindStringSubmatch(body))
	if matchedJSON == "" {
		return nil, static.ErrDataSourceParseFailed
	}
	matchedJSON = utils.GetJSONFromJSObjStr(matchedJSON)

	var gallery []gallery
	err = json.Unmarshal([]byte(matchedJSON), &gallery)
	if err != nil {
		return nil, err
	}

	CDNURL, err := url.Parse(CDN)
	if err != nil {
		return nil, err
	}

	var URLs []*static.URL
	for _, page := range gallery {
		imgURL, err := CDNURL.Parse(page.ImageURI)
		if err != nil {
			return nil, err
		}
		URLs = append(URLs, &static.URL{
			URL: imgURL.String(),
			Ext: utils.GetFileExt(page.ImageURI),
		})
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{"0": {
			Type: static.DataTypeImage,
			URLs: URLs,
		}},
		URL: URL,
	}, nil
}

// Extract implements [static.Extractor].
func (e extractor) Extract(URL string) ([]*static.Data, error) {
	URLs, err := parseURL(URL)
	if err != nil {
		return nil, err
	}
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, id := range URLs {
		d, err := extractData(id)
		if err != nil {
			return nil, utils.Wrap(err, id)
		}
		data = append(data, d)
	}

	return data, nil
}

// New returns a eahentai extractor
func New() static.Extractor {
	return extractor{}
}
