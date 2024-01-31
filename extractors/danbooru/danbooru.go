package danbooru

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request/webdriver"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://danbooru.donmai.us"

var reIMGData = regexp.MustCompile(`data-width="([^"]+)"[ ]+data-height="([^"]+)"[\s\S]*?alt="([^"]+)".+src="([^"]+)"`)

type extractor struct{}

// New returns a danbooru extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract for danbooru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	config.FakeHeaders["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0"

	posts, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []*static.Data{}
	for _, post := range posts {
		contentData, err := extractData(site + post)
		if err != nil {
			return nil, utils.Wrap(err, site+post)
		}
		data = append(data, contentData)
	}

	return data, nil
}

// parseURL for danbooru pages
func parseURL(URL string) ([]string, error) {
	re := regexp.MustCompile(`page=([0-9]+)`)
	pageNo := re.FindAllString(URL, -1)
	// pageNo = URL?page=number -> if it's there it means overview page otherwise single post or invalid
	if len(pageNo) == 0 {

		re := regexp.MustCompile(`/posts/[0-9]+`)
		linkToPost := re.FindString(URL)
		if linkToPost == "" {
			return nil, static.ErrURLParseFailed
		}

		return []string{linkToPost}, nil
	}

	wd, err := webdriver.New()
	if err != nil {
		return nil, err
	}
	defer wd.Close()

	htmlString, err := wd.Get(URL)
	if err != nil {
		return nil, err
	}

	re = regexp.MustCompile(`data-id="([^"]+)`)
	matchedIDs := re.FindAllStringSubmatch(htmlString, -1)

	out := []string{}
	for _, submatchID := range matchedIDs {
		out = append(out, "/posts/"+submatchID[1])
	}

	return out, nil
}

func extractData(postURL string) (*static.Data, error) {
	wd, err := webdriver.New()
	if err != nil {
		return nil, err
	}
	defer wd.Close()

	htmlString, err := wd.Get(postURL)
	if err != nil {
		return nil, err
	}

	matchedImgData := reIMGData.FindStringSubmatch(htmlString)
	if len(matchedImgData) != 5 {
		log.Println(htmlString)
		return nil, static.ErrDataSourceParseFailed
	}
	// [1] = img original width [2] image original height [3] image name [4] src URL

	return &static.Data{
		Site:  site,
		Title: matchedImgData[3],
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: []*static.URL{
					{
						URL: matchedImgData[4],
						Ext: utils.GetLastItemString(strings.Split(matchedImgData[4], ".")),
					},
				},
				Quality: fmt.Sprintf("%s x %s", matchedImgData[1], matchedImgData[2]),
				Size:    0,
			},
		},
		URL: postURL,
	}, nil
}
