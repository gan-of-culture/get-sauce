package imgboard

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

var siteURL string
var mass bool

type extractor struct{}

// New returns a booru imgboard extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract post data
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	re := regexp.MustCompile(`https://[^/]*`)
	siteURL = re.FindString(URL)
	mass = false

	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	var data []*static.Data
	var extractDataFunc func(url string) (static.Data, error)
	extractDataFunc = extractData
	if mass {
		extractDataFunc = extractDataFromDirectLink
	}

	for _, u := range URLs {
		d, err := extractDataFunc(u)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
	}

	return data, nil
}

// parseURL of input
func parseURL(URL string) []string {

	re := regexp.MustCompile(`(?:show/|&id=)[0-9]*`)
	if re.MatchString(URL) {
		return []string{URL}
	}

	re = regexp.MustCompile(`(?:s=list|post\?|page=[0-9]+)`)
	if !re.MatchString(URL) {
		return []string{}
	}

	pageParam := ""
	if strings.Contains(URL, "index.php?") {
		pageParam = "&pid=%d"
	}
	if strings.Contains(URL, "post?") {
		pageParam = "&page=%d"
	}

	re = regexp.MustCompile(`(.+(?:pid=|page=))([0-9]+)([^\s]+)?`) //1=basequeryurl 2=current page 3=parameters after the page parameter
	matchedBaseQueryURL := re.FindStringSubmatch(URL)
	baseQueryURL := ""
	switch len(matchedBaseQueryURL) {
	case 0, 1:
		baseQueryURL = fmt.Sprintf("%s%s", URL, pageParam)
	case 2:
		baseQueryURL = fmt.Sprintf("%s%s", matchedBaseQueryURL[1], "%d")
	case 4:
		baseQueryURL = fmt.Sprintf("%s%s%s", matchedBaseQueryURL[1], "%d", matchedBaseQueryURL[3])
	}

	rePost := regexp.MustCompile(`(?:index.php\?page=post(?:(?:&)|(?:&amp;))s=view(?:(?:&)|(?:&amp;))id=[0-9]*)|"/post/show/[^"]*`)
	reDirectLinks := regexp.MustCompile(`directlink largeimg"\s*href="([^"]*)`)
	found := 0
	urls := []string{}
	pID := 0

	// if the url contains a specific page number and there is no amount set
	// scrape only this page
	if config.Amount == 0 && len(matchedBaseQueryURL) >= 3 {
		pID, _ = strconv.Atoi(matchedBaseQueryURL[2])
	}

	for i := pID; ; {
		htmlString, err := request.Get(fmt.Sprintf(baseQueryURL, i))
		if err != nil {
			break
		}

		matchedDirectLinks := reDirectLinks.FindAllStringSubmatch(htmlString, -1)
		if len(matchedDirectLinks) > 0 {
			mass = true
			for _, l := range matchedDirectLinks {
				if found >= config.Amount && config.Amount > 0 {
					return urls
				}
				urls = append(urls, l[1])
				found++
			}
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

func extractData(url string) (static.Data, error) {

	re := regexp.MustCompile(`http?s://(?:www.)?([^.]*)`)
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
			return static.Data{}, static.ErrDataSourceParseFailed
		}
	}

	if !strings.HasPrefix(matchedPostURL[1], "https") {
		matchedPostURL[1] = fmt.Sprintf("%s%s", "https:", matchedPostURL[1]) //tbib.org/ direct img link has no https:
	}

	var size int64
	if config.Amount == 0 {
		size, err = request.Size(matchedPostURL[1], url)
		if err != nil {
			return static.Data{}, errors.New("no image size not found")
		}
	}

	re = regexp.MustCompile(`Id: [^<]*`)
	id := re.FindString(postHTML)

	re = regexp.MustCompile(`Size: [^<]*`)
	quality := re.FindString(postHTML)
	quality = strings.ReplaceAll(quality, "Size: ", "")

	return static.Data{
		Site:  siteURL,
		Title: fmt.Sprintf("%s_%s", siteName, strings.ReplaceAll(id, "Id: ", "")),
		Type:  utils.GetMediaType(matchedPostURL[2]),
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					{
						URL: matchedPostURL[1],
						Ext: matchedPostURL[2],
					},
				},
				Quality: quality,
				Size:    size,
			},
		},
		Url: url,
	}, nil

}

func extractDataFromDirectLink(url string) (static.Data, error) {
	re := regexp.MustCompile(`https://[^/]*/[^/]*/([^/]*)/[^.\s]*\.[^\.\s]*\..*(\w{3,4})$`) //1=title //2=ext
	matchedURL := re.FindStringSubmatch(url)
	if len(matchedURL) != 3 {
		return static.Data{}, fmt.Errorf("direct download can't match URL %s", url)
	}

	return static.Data{
		Site:  siteURL,
		Title: matchedURL[1],
		Type:  utils.GetMediaType(matchedURL[2]),
		Streams: map[string]*static.Stream{
			"0": {
				URLs: []*static.URL{
					{
						URL: url,
						Ext: matchedURL[2],
					},
				},
				Size: 0,
			},
		},
		Url: url,
	}, nil
}
