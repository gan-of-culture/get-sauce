package booruproject

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

// ParseURL of input
func ParseURL(url string) []string {
	if strings.Contains(url, "s=view") {
		return []string{url}
	}

	if !strings.Contains(url, "s=list") {
		return []string{}
	}

	re := regexp.MustCompile("https://[^/]*")
	siteURL := re.FindString(url)

	re = regexp.MustCompile("(.*)pid=[0-9]*")
	matchedBaseQueryURL := re.FindStringSubmatch(url)
	baseQueryURL := url
	if len(matchedBaseQueryURL) == 2 {
		baseQueryURL = matchedBaseQueryURL[1]
	}

	rePost := regexp.MustCompile("index.php\\?page=post(?:(?:&)|(?:&amp;))s=view(?:(?:&)|(?:&amp;))id=[0-9]*")
	found := 0
	urls := []string{}
	for i := 0; ; i += 42 {
		htmlString, err := request.Get(fmt.Sprintf("%s&pid=%d", baseQueryURL, i))
		if err != nil {
			break
		}

		matchedPosts := rePost.FindAllString(htmlString, -1)

		for _, p := range matchedPosts {
			if found >= config.Amount && config.Amount > 0 {
				return urls
			}
			urls = append(urls, fmt.Sprintf("%s/%s", siteURL, strings.ReplaceAll(p, "&amp;", "&")))
			found++
		}
		if config.Amount == 0 {
			return urls
		}
	}

	return urls
}

// Extract post data
func Extract(url string) ([]static.Data, error) {
	urls := ParseURL(url)
	if len(urls) == 0 {
		return nil, fmt.Errorf("Can't find a post for %s", url)
	}

	re := regexp.MustCompile("https://[^/]*")
	siteURL := re.FindString(url)

	var data []static.Data
	for _, u := range urls {
		d, err := extractData(u, siteURL)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(url string, site string) (static.Data, error) {

	re := regexp.MustCompile("http?s://(?:www.)?([^.]*)")
	siteName := re.FindStringSubmatch(site)[1]

	postHTML, err := request.Get(url)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile("<a href=\"(https.*/?/images[^\"]*\\.([^\"?]*)(?:[^\"])*?)\"[\\s\\S]*Original") //1=url 2=ext
	matchedPostURL := re.FindStringSubmatch(postHTML)
	if len(matchedPostURL) != 3 {
		re := regexp.MustCompile("(https.*/?/images[^.]*\\.([^\"?]*)(?:[^\"])*?)[\\s\\S]*?id=\"image")
		matchedPostURL = re.FindStringSubmatch(postHTML)
		if len(matchedPostURL) != 3 {
			return static.Data{}, err
		}

	}

	var size int64
	if config.Amount == 0 {
		size, err = request.Size(matchedPostURL[1], url)
		if err != nil {
			return static.Data{}, fmt.Errorf("[%s]No image size not found", siteName)
		}
	}

	re = regexp.MustCompile("Id: [^<]*")
	id := re.FindString(postHTML)

	re = regexp.MustCompile("Size: [^<]*")
	quality := re.FindString(postHTML)
	quality = strings.ReplaceAll(quality, "Size: ", "")

	dataType := ""
	switch matchedPostURL[2] {
	case "jpg", "jpeg", "png", "gif", "webp":
		dataType = fmt.Sprintf("%s/%s", "image", matchedPostURL[2])
	case "webm", "mp4", "mkv", "m4a":
		dataType = fmt.Sprintf("%s/%s", "video", matchedPostURL[2])
	default:
		dataType = fmt.Sprintf("%s/%s", "unknown", matchedPostURL[2])
	}

	return static.Data{
		Site:  site,
		Title: fmt.Sprintf("%s_%s", siteName, strings.ReplaceAll(id, "Id: ", "")),
		Type:  dataType,
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: matchedPostURL[1],
						Ext: matchedPostURL[2],
					},
				},
				Quality: quality,
				Size:    size,
			},
		},
		Err: nil,
		Url: url,
	}, nil

}
