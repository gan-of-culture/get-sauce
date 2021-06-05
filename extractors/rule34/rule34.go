package rule34

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://rule34.paheal.net"

type extractor struct{}

// New returns a rule34 extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	urls := parseURL(URL)
	if len(urls) == 0 {
		return nil, errors.New("[Rule34] Can't parse URL")
	}

	var data []*static.Data
	for _, u := range urls {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
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

		re := regexp.MustCompile(`data-post-id="([^"]+)`)
		matchedPosts := re.FindAllStringSubmatch(htmlString, -1)
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

func extractData(URL string) (static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`[0-9]{3,}`)
	id := re.FindStringSubmatch(URL)

	// get source of img
	re = regexp.MustCompile(`id='main_image' src='([^']+)`)
	matchedPostSrcURL := re.FindStringSubmatch(htmlString)
	if len(matchedPostSrcURL) != 2 {

		// maybe it's a video - try to get source URL
		re = regexp.MustCompile(`<source src='([^']+)`)
		matchedPostSrcURL = re.FindStringSubmatch(htmlString)
		if len(matchedPostSrcURL) != 2 {
			return static.Data{}, errors.New("[Rule34] src URL not found for post " + URL)
		}
	}

	postSrcURL := matchedPostSrcURL[1]

	re = regexp.MustCompile(`tag_edit__tags' value='([^']+)`)
	matchedTagBox := re.FindStringSubmatch(htmlString)
	if len(matchedTagBox) != 2 {
		fmt.Println(htmlString)
		return static.Data{}, errors.New("[Rule34] couldn't extract tags for post " + URL)
	}

	title := fmt.Sprintf("%s %s", matchedTagBox[1], id[0])

	var size int64
	if config.Amount == 0 {
		size, err = request.Size(postSrcURL, URL)
		if err != nil {
			return static.Data{}, errors.New("[Rule34]No image size not found")
		}
	}

	dataType := "image"
	if strings.HasSuffix(postSrcURL, ".gif") {
		dataType = "gif"
	}
	if strings.Contains(htmlString, "#Videomain") {
		dataType = "video"
	}

	var quality string
	if dataType == "video" {
		re := regexp.MustCompile(`id='main_image'.+\n[^0-9]+([0-9]+)[^0-9]+([0-9]+)`)
		matchedQualityProperties := re.FindStringSubmatch(htmlString)
		if len(matchedQualityProperties) != 3 {
			return static.Data{}, errors.New("[Rule34] quality not found for post " + URL)
		}
		quality = fmt.Sprintf("%s x %s", matchedQualityProperties[1], matchedQualityProperties[2])
	} else {
		re := regexp.MustCompile(`data-(width|height)='([0-9]+)`)
		matchedQualityProperties := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedQualityProperties) != 2 {
			return static.Data{}, errors.New("[Rule34] quality not found for post " + URL)
		}

		quality = fmt.Sprintf("%s x %s", matchedQualityProperties[1][2], matchedQualityProperties[0][2])
	}

	return static.Data{
		Site:  site,
		Title: title,
		Type:  dataType,
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: postSrcURL,
						Ext: utils.GetLastItemString(strings.Split(postSrcURL, ".")),
					},
				},
				Quality: quality,
				Size:    size,
			},
		},
		Url: URL,
	}, nil
}
