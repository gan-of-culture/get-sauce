package imgboard

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

var siteURL string
var mass = false

// ParseURL of input
func ParseURL(url string) []string {

	re := regexp.MustCompile("(?:show/|&id=)[0-9]*")
	if re.MatchString(url) {
		return []string{url}
	}

	re = regexp.MustCompile(`(?:s=list|post\?|page=[0-9]+)`)
	if !re.MatchString(url) {
		return []string{}
	}

	pageParam := ""
	if strings.Contains(url, "index.php?") {
		pageParam = "&pid=%d"
	}
	if strings.Contains(url, "post?") {
		pageParam = "&page=%d"
	}

	re = regexp.MustCompile(`(.+(?:pid=|page=))[0-9]+([^\s]+)?`) //1=basequeryurl 2=parameters after the page parameter
	matchedBaseQueryURL := re.FindStringSubmatch(url)
	baseQueryURL := ""
	switch len(matchedBaseQueryURL) {
	case 0, 1:
		baseQueryURL = fmt.Sprintf("%s%s", url, pageParam)
	case 2:
		baseQueryURL = fmt.Sprintf("%s%s", matchedBaseQueryURL[1], "%d")
	case 3:
		baseQueryURL = fmt.Sprintf("%s%s%s", matchedBaseQueryURL[1], "%d", matchedBaseQueryURL[2])
	}

	rePost := regexp.MustCompile(`(?:index.php\?page=post(?:(?:&)|(?:&amp;))s=view(?:(?:&)|(?:&amp;))id=[0-9]*)|"/post/show/[^"]*`)
	reDirectLinks := regexp.MustCompile(`directlink largeimg"\s*href="([^"]*)`)
	found := 0
	mass = false
	urls := []string{}
	for i := 0; ; {
		htmlString, err := request.Get(fmt.Sprintf(baseQueryURL, i))
		if err != nil {
			break
		}

		matchedDirectLinks := reDirectLinks.FindAllStringSubmatch(htmlString, -1)
		if len(matchedDirectLinks) > 0 {
			for _, l := range matchedDirectLinks {
				if found >= config.Amount && config.Amount > 0 {
					return urls
				}
				urls = append(urls, l[1])
				found++
			}
			mass = true
		}

		if !mass {
			matchedPosts := rePost.FindAllString(htmlString, -1)
			if len(matchedPosts) == 0 {
				return urls
			}

			for _, p := range matchedPosts {
				if found >= config.Amount && config.Amount > 0 {
					return urls
				}
				p = strings.TrimLeft(p, `"/`)
				p = strings.ReplaceAll(p, "&amp;", "&")

				urls = append(urls, fmt.Sprintf("%s/%s", siteURL, p))
				found++
			}
		}
		if config.Amount == 0 {
			return urls
		}
		if pageParam == "&pid=%d" {
			i += 42
			continue
		}
		i++
	}

	return urls
}

// Extract post data
func Extract(url string) ([]static.Data, error) {
	re := regexp.MustCompile("https://[^/]*")
	siteURL = re.FindString(url)

	urls := ParseURL(url)
	if len(urls) == 0 {
		return nil, fmt.Errorf("Can't find a post for %s", url)
	}

	var data []static.Data
	if mass {
		for _, u := range urls {
			d, err := extractDataFromDirectLink(u)
			if err != nil {
				return nil, err
			}
			data = append(data, d)
		}
	} else {
		for _, u := range urls {
			d, err := extractData(u)
			if err != nil {
				return nil, err
			}
			data = append(data, d)
		}
	}

	return data, nil
}

func extractData(url string) (static.Data, error) {

	re := regexp.MustCompile("http?s://(?:www.)?([^.]*)")
	siteName := re.FindStringSubmatch(siteURL)[1]

	postHTML, err := request.Get(url)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`<a.+href="([^"]+\.([^"?]+)).+>\s*(?:Original|Download PNG)`) //1=url 2=ext
	matchedPostURL := re.FindStringSubmatch(postHTML)
	if len(matchedPostURL) != 3 {
		re = regexp.MustCompile(`<a.+href="([^"]+\.([^"?]+)).+>\s*(?:Original|View larger|Download PNG)`)
		matchedPostURL = re.FindStringSubmatch(postHTML)
		if len(matchedPostURL) != 3 {
			return static.Data{}, err
		}
	}

	if !strings.HasPrefix(matchedPostURL[1], "https") {
		matchedPostURL[1] = fmt.Sprintf("%s%s", "https:", matchedPostURL[1]) //tbib.org/ direct img link has no https:
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

	return static.Data{
		Site:  siteURL,
		Title: fmt.Sprintf("%s_%s", siteName, strings.ReplaceAll(id, "Id: ", "")),
		Type:  utils.GetMediaType(matchedPostURL[2]),
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

func extractDataFromDirectLink(url string) (static.Data, error) {
	re := regexp.MustCompile(`https://[^/]*/[^/]*/([^/]*)/[^.\s]*\.[^\.\s]*\.(\w{3,4})`) //1=title //2=ext
	matchedURL := re.FindStringSubmatch(url)
	if len(matchedURL) != 3 {
		return static.Data{}, fmt.Errorf("[IMGBoard] direct download can't match URL %s", url)
	}

	return static.Data{
		Site:  siteURL,
		Title: matchedURL[1],
		Type:  utils.GetMediaType(matchedURL[2]),
		Streams: map[string]static.Stream{
			"0": {
				URLs: []static.URL{
					{
						URL: matchedURL[1],
						Ext: matchedURL[2],
					},
				},
				Quality: "best",
				Size:    0,
			},
		},
		Err: nil,
		Url: url,
	}, nil
}
