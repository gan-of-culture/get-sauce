package danbooru

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
)

const site = "https://danbooru.donmai.us"

var reIMGData = regexp.MustCompile(`data-width="([^"]+)"[ ]+data-height="([^"]+)"[\s\S]*?alt="([^"]+)".+src="([^"]+)"`)

type extractor struct {
	client *http.Client
}

// New returns a danbooru extractor
func New() static.Extractor {
	return newForTesting()
}

func newForTesting() *extractor {
	return &extractor{
		client: request.Firefox117Client(),
	}
}

// Extract for danbooru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	posts, err := e.parseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []*static.Data{}
	for _, post := range posts {
		contentData, err := e.extractData(site + post)
		if err != nil {
			log.Println(site + post)
			return nil, err
		}
		data = append(data, contentData)
	}

	return data, nil
}

// parseURL for danbooru pages
func (e *extractor) parseURL(URL string) ([]string, error) {
	re := regexp.MustCompile(`page=([0-9]+)`)
	pageNo := re.FindAllString(URL, -1)
	// pageNo = URL?page=number -> if it's there it means overview page otherwise single post or invalid
	if len(pageNo) == 0 {

		re := regexp.MustCompile(`/posts/[0-9]+`)
		linkToPost := re.FindString(URL)
		if linkToPost == "" {
			return nil, errors.WithStack(static.ErrURLParseFailed)
		}

		return []string{linkToPost}, nil
	}

	htmlString, err := request.GetAsBytesWithClient(e.client, URL, URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	re = regexp.MustCompile(`data-id="([^"]+)`)
	matchedIDs := re.FindAllSubmatch(htmlString, -1)

	out := []string{}
	for _, submatchID := range matchedIDs {
		out = append(out, "/posts/"+string(submatchID[1]))
	}

	return out, nil
}

func (e *extractor) extractData(postURL string) (*static.Data, error) {
	htmlString, err := request.GetAsBytesWithClient(e.client, postURL, postURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	matchedImgData := reIMGData.FindStringSubmatch(string(htmlString))
	if len(matchedImgData) != 5 {
		log.Println(htmlString)
		return nil, errors.WithStack(static.ErrDataSourceParseFailed)
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
