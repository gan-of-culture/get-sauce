package hstream

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	mpegdash "github.com/gan-of-culture/get-sauce/parsers/mpeg_dash"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://hstream.moe/"

var reVideoSources = regexp.MustCompile(`https://.+/\d+/[\w.]+/[\w./]+\.(?:mpd|mp4|webm)`)
var reCaptionSource = regexp.MustCompile(`https://.+/\d+/[\w.]+/[\w./]+\.ass`)

type extractor struct{}

// New returns a hstream extractor
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)

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
	if ok, _ := regexp.MatchString(site+`hentai/[\w\-]+/\d+`, URL); ok {
		return []string{URL}
	}

	if ok, _ := regexp.MatchString(site+`hentai/[\w\-]+/?`, URL); !ok {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	htmlString = strings.Split(htmlString, `class="bixbox"`)[0]

	re := regexp.MustCompile(`hentai/[\w\-]+/\d+`)
	out := []string{}
	for _, episode := range re.FindAllString(htmlString, -1) {
		out = append(out, site+episode)
	}
	return utils.RemoveAdjDuplicates(out)
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	if strings.Contains(htmlString, "<title>DDOS-GUARD</title>") {
		time.Sleep(300 * time.Millisecond)
		htmlString, _ = request.Get(URL)
	}

	videoSources := reverse(reVideoSources.FindAllString(htmlString, -1))

	streams := make(map[string]*static.Stream)
	counter := 0
	// only keep one audio stream
	foundAudioStream := false
	for _, sourceURL := range videoSources {
		if !strings.HasSuffix(sourceURL, ".mpd") {
			continue
		}

		streamsTmp, err := mpegdash.ExtractDASHManifest(sourceURL, map[string]string{"Referer": site})
		if err != nil {
			return nil, err
		}
		for _, streamTmp := range streamsTmp {
			if streamTmp.Type == static.DataTypeAudio {
				if foundAudioStream {
					continue
				} else {
					foundAudioStream = true
				}
			}

			streams[fmt.Sprint(counter)] = streamTmp
			counter += 1
		}
	}

	// skip direct file downloads if newer mpd is supplied
	if len(streams) > 0 {
		videoSources = []string{}
	}

	for i, sourceURL := range videoSources {
		size, err := request.Size(sourceURL, site)
		if err != nil {
			return nil, err
		}

		streams[fmt.Sprint(i)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: sourceURL,
					Ext: utils.GetFileExt(sourceURL),
				},
			},
			Quality: utils.GetLastItemString(strings.Split(sourceURL, "/")),
			Size:    size,
		}
	}

	captionURL := reCaptionSource.FindString(htmlString)

	return &static.Data{
		Site:    site,
		Title:   utils.GetSectionHeadingElement(&htmlString, 6, 0),
		Type:    static.DataTypeVideo,
		Streams: streams,
		Captions: []*static.Caption{
			{
				URL: static.URL{
					URL: captionURL,
					Ext: utils.GetFileExt(captionURL),
				},
				Language: "English",
			},
		},
		URL: URL,
	}, nil

}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
