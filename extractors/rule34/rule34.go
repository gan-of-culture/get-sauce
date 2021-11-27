package rule34

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://rule34.paheal.net"

var reParsePostID = regexp.MustCompile(`data-post-id="([^"]+)`)
var rePostID = regexp.MustCompile(`[0-9]{3,}`)
var reSourceURL = regexp.MustCompile(`id='main_image' src='([^']+)`)
var reVideoSourceURL = regexp.MustCompile(`<source src='([^']+)`)
var reTagBox = regexp.MustCompile(`tag_edit__tags' value='([^']+)`)
var reQuality = regexp.MustCompile(`data-(width|height)='([0-9]+)`)
var reVideoQuality = regexp.MustCompile(`id='main_image'.+\n[^0-9]+([0-9]+)[^0-9]+([0-9]+)`)

type extractor struct{}

// New returns a rule34 extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	var data []*static.Data
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
	}
	return data, nil
}

// parseURL data
func parseURL(URL string) []string {

	// if it's single post return
	if strings.Contains(URL, "/post/view/") {
		return []string{URL}
	}

	// everything other than a overview page gets returned
	if !strings.Contains(URL, "/post/list/") {
		return nil
	}

	re := regexp.MustCompile(`(\S*)/[0-9]*?$`)
	baseURL := re.FindStringSubmatch(URL)[1]

	content := []string{}
	found := 0
	for i := 1; ; i++ {
		htmlString, err := request.Get(fmt.Sprintf("%s/%d", baseURL, i))
		if err != nil {
			return nil
		}

		matchedPosts := reParsePostID.FindAllStringSubmatch(htmlString, -1)
		if len(matchedPosts) < 1 {
			return content
		}

		for _, post := range matchedPosts {
			content = append(content, fmt.Sprintf("%s/post/view/%s", site, post[1]))
			found++
			if found >= config.Amount && config.Amount != 0 {
				return content
			}
		}
		if config.Amount == 0 {
			return content
		}
	}
}

func extractData(URL string) (*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	id := rePostID.FindStringSubmatch(URL)

	// get source of img
	matchedPostSrcURL := reSourceURL.FindStringSubmatch(htmlString)
	if len(matchedPostSrcURL) != 2 {

		// maybe it's a video - try to get source URL
		matchedPostSrcURL = reVideoSourceURL.FindStringSubmatch(htmlString)
		if len(matchedPostSrcURL) != 2 {
			return nil, static.ErrDataSourceParseFailed
		}
	}

	postSrcURL := matchedPostSrcURL[1]

	matchedTagBox := reTagBox.FindStringSubmatch(htmlString)
	if len(matchedTagBox) != 2 {
		fmt.Println(htmlString)
		return nil, errors.New("couldn't extract tags for post")
	}

	title := fmt.Sprintf("%s %s", matchedTagBox[1], id[0])

	var size int64
	if config.Amount == 0 {
		size, err = request.Size(postSrcURL, URL)
		if err != nil {
			return nil, errors.New("no image size not found")
		}
	}

	dataType := static.DataTypeImage
	if strings.Contains(htmlString, "#Videomain") {
		dataType = static.DataTypeVideo
	}

	var quality string
	if dataType == static.DataTypeVideo {
		matchedQualityProperties := reVideoQuality.FindStringSubmatch(htmlString)
		if len(matchedQualityProperties) != 3 {
			return nil, errors.New("quality not found for post ")
		}
		quality = fmt.Sprintf("%s x %s", matchedQualityProperties[1], matchedQualityProperties[2])
	} else {
		matchedQualityProperties := reQuality.FindAllStringSubmatch(htmlString, -1)
		if len(matchedQualityProperties) != 2 {
			return nil, errors.New("quality not found for post ")
		}

		quality = fmt.Sprintf("%s x %s", matchedQualityProperties[1][2], matchedQualityProperties[0][2])
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  dataType,
		Streams: map[string]*static.Stream{
			"0": {
				Type: dataType,
				URLs: []*static.URL{
					{
						URL: postSrcURL,
						Ext: utils.GetLastItemString(strings.Split(postSrcURL, ".")),
					},
				},
				Quality: quality,
				Size:    size,
			},
		},
		URL: URL,
	}, nil
}
