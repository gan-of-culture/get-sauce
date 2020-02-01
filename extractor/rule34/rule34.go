package rule34

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://rule34.paheal.net"

// Extract for data
func Extract(URL string) ([]static.Data, error) {

	data := []static.Data{}

	for _, postURL := range ParseURL(URL) {

		postData, err := extractData(postURL)
		if err != nil {
			return nil, err
		}
		data = append(data, postData)
	}

	return data, nil
}

// ParseURL data
func ParseURL(URL string) []string {

	// if it's single post return
	if strings.Contains(URL, "/post/view/") {
		return []string{URL}
	}

	// everything other than a overview page gets returned
	if !strings.Contains(URL, "/post/list/") {
		return nil
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	re := regexp.MustCompile("data-post-id=\"([^\"]+)")
	matchedPosts := re.FindAllStringSubmatch(htmlString, -1)
	if len(matchedPosts) < 1 {
		return nil
	}

	content := []string{}
	for _, post := range matchedPosts {
		content = append(content, fmt.Sprintf("%s/post/view/%s", site, post[1]))
	}

	return content
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile("tag_edit__tags' value='([^']+)")
	matchedTagBox := re.FindStringSubmatch(htmlString)
	if len(matchedTagBox) != 2 {
		return static.Data{}, errors.New("[Rule34] couldn't extract tags for post " + URL)
	}
	title := matchedTagBox[1]

	// get source of img
	re = regexp.MustCompile("id='main_image' src='([^']+)")
	matchedPostSrcURL := re.FindStringSubmatch(htmlString)
	if len(matchedPostSrcURL) != 2 {

		// maybe it's a video - try to get source URL
		re = regexp.MustCompile("<source src='([^']+)")
		matchedPostSrcURL = re.FindStringSubmatch(htmlString)
		if len(matchedPostSrcURL) != 2 {
			return static.Data{}, errors.New("[Rule34] src URL not found for post " + URL)
		}
	}

	postSrcURL := matchedPostSrcURL[1]

	size, err := request.Size(postSrcURL, URL)
	if err != nil {
		return static.Data{}, errors.New("[Rule34]No image size not found")
	}

	postType := "image"
	if strings.HasSuffix(postSrcURL, ".gif") {
		postType = "gif"
	}
	if strings.Contains(htmlString, "#Videomain") {
		postType = "video"
	}

	var postQuality string
	if postType == "video" {
		re = regexp.MustCompile("id='main_image'.+\n[^0-9]+([0-9]+)[^0-9]+([0-9]+)")
		matchedQualityProperties := re.FindStringSubmatch(htmlString)
		if len(matchedQualityProperties) != 3 {
			return static.Data{}, errors.New("[Rule34] quality not found for post " + URL)
		}
		postQuality = fmt.Sprintf("%s x %s", matchedQualityProperties[1], matchedQualityProperties[2])
	} else {
		re = regexp.MustCompile("data-(width|height)='([0-9]+)")
		matchedQualityProperties := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedQualityProperties) != 2 {
			return static.Data{}, errors.New("[Rule34] quality not found for post " + URL)
		}

		postQuality = fmt.Sprintf("%s x %s", matchedQualityProperties[1][2], matchedQualityProperties[0][2])
	}

	return static.Data{
		Site:  site,
		Title: title,
		Type:  postType,
		Streams: map[string]static.Stream{
			"0": static.Stream{
				URLs: []static.URL{
					static.URL{
						URL: postSrcURL,
						Ext: utils.GetLastItemString(strings.Split(postSrcURL, ".")),
					},
				},
				Quality: postQuality,
				Size:    size,
			},
		},
		Url: URL,
	}, nil
}
