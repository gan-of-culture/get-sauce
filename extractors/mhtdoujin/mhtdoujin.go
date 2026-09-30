package mhtdoujin

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
)

type extractor struct{}

type page struct {
	Page   int    `json:"page"`
	Ext    string `json:"ext"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

var reGalleryID = regexp.MustCompile(`<a href="(/gallery/\d+)`)
var reImageBaseURL = regexp.MustCompile(`data-reader-image-base="([^"]+)`)
var reTitle = regexp.MustCompile(`data-reader-seo-title="([^"]+)`)
var rePagesJsonStr = regexp.MustCompile(`id="readerPagesJson">([^<]+)`)

// New returns a modern (layout) htdoujin extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	galleryURLs, err := parseURL(u)
	if err != nil {
		return nil, err
	}

	data := []*static.Data{}
	for _, galleryURL := range galleryURLs {
		d, err := extractData(galleryURL)
		if err != nil {
			return nil, utils.Wrap(err, galleryURL.String())
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL *url.URL) ([]*url.URL, error) {
	// https://hentaienvy.com/gallery/1615288/
	URLPath := strings.TrimPrefix(URL.Path, "/")
	URLPathElements := strings.Split(URLPath, "/")
	if len(URLPathElements) > 1 {
		_, err := strconv.Atoi(URLPathElements[1])
		if URLPathElements[0] == "gallery" && err == nil {
			return []*url.URL{URL}, nil
		}
	}

	HTMLString, err := request.Get(URL.String())
	if err != nil {
		return nil, err
	}

	matchedIDs := reGalleryID.FindAllStringSubmatch(HTMLString, -1)
	var out []*url.URL
	for _, URLPath := range matchedIDs {
		galleryURL, err := URL.Parse(URLPath[1])
		if err != nil {
			return nil, errors.WithStack(err)
		}
		out = append(out, galleryURL)
	}

	return out, nil
}

func extractData(URL *url.URL) (*static.Data, error) {
	URLPath := strings.Trim(URL.Path, "/")
	URLPathElements := strings.Split(URLPath, "/")
	URLPathElements[0] = "g"
	URLPathElements = append(URLPathElements, "1")

	readerURL, err := URL.Parse(fmt.Sprintf("/%s", path.Join(URLPathElements...)))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	readerHTML, err := request.Get(readerURL.String())
	if err != nil {
		return nil, err
	}

	matchedImageBaseURL := reImageBaseURL.FindStringSubmatch(readerHTML)
	if len(matchedImageBaseURL) < 2 {
		return nil, errors.WithStack(errors.New("unable to parse image base URL"))
	}
	imageBaseURL, err := url.Parse(matchedImageBaseURL[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	matchedTitle := reTitle.FindStringSubmatch(readerHTML)
	if len(matchedTitle) < 2 {
		return nil, errors.WithStack(errors.New("unable to parse gallery title"))
	}

	matchedPagesJsonStr := rePagesJsonStr.FindStringSubmatch(readerHTML)
	if len(matchedPagesJsonStr) < 2 {
		return nil, errors.WithStack(errors.New("unable to parse pages json struct"))
	}

	var pages []page
	err = json.Unmarshal([]byte(matchedPagesJsonStr[1]), &pages)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var URLs []*static.URL
	for _, page := range pages {
		pageURL := imageBaseURL.Clone()
		pageURL = pageURL.JoinPath(fmt.Sprintf("%d.%s", page.Page, page.Ext))
		URLs = append(URLs, &static.URL{
			URL: pageURL.String(),
			Ext: page.Ext,
		})
	}

	return &static.Data{
		Site:  URL.Host,
		Title: html.UnescapeString(matchedTitle[1]),
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": &static.Stream{
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: URL.String(),
	}, nil
}
