package nhgroup

import (
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/gan-of-culture/get-sauce/extractors/animestream"
	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/extractors/nhplayer"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reNHPlayerURL = regexp.MustCompile(`https:\\?/\\?/nhplayer\.com\\?/v\\?/[^/"]+`)

type extractor struct{}

// New returns a nhgroup extractor
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := ExtractData(u)
		if err != nil {
			if strings.Contains(err.Error(), "video not found") || strings.Contains(err.Error(), "player URL not found") {
				log.Println(utils.Wrap(err, u).Error())
				continue
			}
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
	}

	return data, nil
}

func ParseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`episode-\d+[/_\-]*`, URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https[^"\s]*?episode-\d*(?:/*|[-\w]*)"`)
	matchedURLs := re.FindAllString(htmlString, -1)
	if len(matchedURLs) == 0 {
		matchedURLs = animestream.ParseURLwoSite(URL)
	}

	sort.Strings(matchedURLs)
	matchedURLs = utils.RemoveAdjDuplicates(matchedURLs)

	out := []string{}
	for _, u := range matchedURLs {
		out = append(out, strings.Trim(u, `"`))
	}

	return out
}

// ExtractData of a nhplayer
func ExtractData(URL string) (*static.Data, error) {

	htmlString, err := request.GetWithHeaders(URL, map[string]string{
		"Cookie": "inter=1",
	})
	if err != nil {
		return nil, err
	}

	playerURL := reNHPlayerURL.FindString(htmlString)
	playerURL = strings.ReplaceAll(playerURL, `\`, "")
	if playerURL != "" {
		data, err := nhplayer.New().Extract(playerURL)
		if err != nil {
			return nil, err
		}
		return data[0], err
	}

	return htstreaming.ExtractData(URL)
}
