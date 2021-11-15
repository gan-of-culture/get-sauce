package danbooru

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/v2/request"
	"github.com/gan-of-culture/get-sauce/v2/static"
	"github.com/gan-of-culture/get-sauce/v2/utils"
)

const site = "https://danbooru.donmai.us"

type extractor struct{}

// New returns a danbooru extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract for danbooru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
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
		data = append(data, &contentData)
	}

	return data, nil
}

// parseURL for danbooru pages
func parseURL(URL string) ([]string, error) {
	re := regexp.MustCompile(`page=([0-9]+)`)
	pageNo := re.FindAllString(URL, -1)
	// pageNo = url?page=number -> if it's there it means overview page otherwise single post or invalid
	if len(pageNo) == 0 {

		re := regexp.MustCompile(`/posts/[0-9]+`)
		linkToPost := re.FindString(URL)
		if linkToPost == "" {
			return nil, static.ErrURLParseFailed
		}

		return []string{linkToPost}, nil
	}

	htmlString, err := request.Get(URL)
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

func extractData(postURL string) (static.Data, error) {
	htmlString, err := request.Get(postURL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`data-width="([^"]+)"[ ]+data-height="([^"]+)".+alt="([^"]+)".+src="([^"]+)"`)
	matchedImgData := re.FindStringSubmatch(htmlString)
	if len(matchedImgData) != 5 {
		return static.Data{}, static.ErrDataSourceParseFailed
	}
	// [1] = img original width [2] image original height [3] image name [4] src url

	size, err := request.Size(matchedImgData[4], postURL)
	if err != nil {
		return static.Data{}, errors.New("no image size not found")
	}

	return static.Data{
		Site:  site,
		Title: matchedImgData[3],
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					{
						URL: matchedImgData[4],
						Ext: utils.GetLastItemString(strings.Split(matchedImgData[4], ".")),
					},
				},
				Quality: fmt.Sprintf("%s x %s", matchedImgData[1], matchedImgData[2]),
				Size:    size,
			},
		},
		Url: postURL,
	}, nil
}
