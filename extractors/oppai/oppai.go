package oppai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://oppai.stream/"
const episodeURLTemplate = "https://oppai.stream/watch.php?e="

var reSources = regexp.MustCompile(`var availableres = ({[^}]+})`)           // 1=srcURL 2=Resolution
var reCaptions = regexp.MustCompile(`<track.+label="([^"]+)"\ssrc="([^"]+)`) // 1=Language 2=srcURL

type extractor struct{}

// New returns a oppai.stream extractor
func New() static.Extractor {
	return &extractor{}
}

// Extract data from URL
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	episodeSlugs := parseURL(URL)
	if len(episodeSlugs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range episodeSlugs {
		d, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	u, err := url.Parse(URL)
	if err != nil {
		return nil
	}

	if episodeSlug := u.Query().Get("e"); episodeSlug != "" {
		return []string{episodeSlug}
	}

	return nil
}

func extractData(episodeSlug string) (*static.Data, error) {
	URL := fmt.Sprintf("%s%s", episodeURLTemplate, episodeSlug)

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedSources := map[string]string{}
	sourcesKeys := []int{}
	err = json.Unmarshal([]byte(reSources.FindStringSubmatch(htmlString)[1]), &matchedSources)
	if err != nil {
		return nil, err
	}

	sources := map[int]string{}
	for res, matchedURL := range matchedSources {
		resKey := res
		if resKey == "" {
			continue
		}
		if resKey == "4k" {
			resKey = "2160"
		}
		resKeyAsInt, err := strconv.Atoi(resKey)
		if err != nil {
			return nil, err
		}
		sources[resKeyAsInt] = matchedURL
		sourcesKeys = append(sourcesKeys, resKeyAsInt)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sourcesKeys)))

	streams := map[string]*static.Stream{}
	for idx, sourceKey := range sourcesKeys {
		srcURL := sources[sourceKey]
		size, _ := request.Size(srcURL, site)

		streams[fmt.Sprint(idx)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: srcURL,
					Ext: utils.GetFileExt(srcURL),
				},
			},
			Quality: fmt.Sprintf("%dp", sourceKey),
			Size:    size,
		}
	}

	caption := &static.Caption{}
	matchedCaption := reCaptions.FindStringSubmatch(htmlString)
	if len(matchedCaption) == 3 {
		caption.URL.URL = matchedCaption[2]
		caption.URL.Ext = utils.GetFileExt(matchedCaption[2])
		caption.Language = matchedCaption[1]
	}

	return &static.Data{
		Site:     site,
		Title:    episodeSlug,
		Type:     static.DataTypeVideo,
		Streams:  streams,
		Captions: []*static.Caption{caption},
		URL:      URL,
	}, nil
}
