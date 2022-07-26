package cartoonpornvideos

import (
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://www.cartoonpornvideos.com/"

var reEpisodeURL = regexp.MustCompile(`https://www\.cartoonpornvideos\.com/(?:click/\d-\d+/)*video/.+\w{11}\.html`)
var reHLSURL = regexp.MustCompile(`https://hls\.cartoonpornvideos\.com/[^"]+`)

type extractor struct{}

// New returns a cartoonpornvideos extractor.
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
	if reEpisodeURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	re := regexp.MustCompile(`click/\d-\d/`)
	out := []string{}
	for _, URL := range utils.RemoveAdjDuplicates(reEpisodeURL.FindAllString(htmlString, -1)) {
		out = append(out, re.ReplaceAllString(URL, ""))
	}
	return out
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.GetWithHeaders(URL, map[string]string{
		"Referer": site,
	})
	if err != nil {
		return nil, err
	}

	masterURL := reHLSURL.FindString(htmlString)
	if masterURL == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	streams, err := request.ExtractHLS(masterURL, map[string]string{"Referer": URL})
	if err != nil {
		return nil, err
	}

	var ext string
	for _, stream := range streams {
		for _, u := range stream.URLs {
			u.Ext = "ts"
		}

		ext = stream.URLs[0].Ext

		if strings.Contains(stream.Info, "mp4a") {
			ext = "mp4"
		}

		stream.Ext = ext
	}

	return &static.Data{
		Site:    site,
		Title:   utils.GetH1(&htmlString, -1),
		Type:    static.DataTypeVideo,
		Streams: streams,
		URL:     URL,
	}, nil
}
