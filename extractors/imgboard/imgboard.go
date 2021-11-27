package imgboard

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

var reSiteName = regexp.MustCompile(`http?s://(?:www.)?([^.]*)`)
var rePostURL = regexp.MustCompile(`<a.+href="([^"]+\.([^"?]+)).+>\s*(?:Original|Download PNG)`)                //1=url 2=ext
var rePostBackup = regexp.MustCompile(`<a.+href="([^"]+\.([^"?]+)).+>\s*(?:Original|View larger|Download PNG)`) //1=url 2=ext
var reID = regexp.MustCompile(`Id: [^<]*`)
var reSize = regexp.MustCompile(`Size: [^<]*`)
var reDirectLink = regexp.MustCompile(`https://[^/]*/[^/]*/([^/]*)/[^.\s]*\.[^\.\s]*\..*(\w{3,4})$`) //1=title //2=ext

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
	var extractDataFunc func(URL string) (*static.Data, error)
	extractDataFunc = extractData
	if mass {
		extractDataFunc = extractDataFromDirectLink
	}

	for _, u := range URLs {
		d, err := extractDataFunc(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
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
	URLs := []string{}
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
					return URLs
				}
				URLs = append(URLs, l[1])
				found++
			}
		}

		if !mass {
			matchedPosts := rePost.FindAllString(htmlString, -1)
			if len(matchedPosts) == 0 {
				return URLs
			}

			for _, p := range matchedPosts {
				if found >= config.Amount && config.Amount > 0 {
					return URLs
				}
				p = strings.TrimLeft(p, `"/`)
				p = strings.ReplaceAll(p, "&amp;", "&")

				URLs = append(URLs, fmt.Sprintf("%s/%s", siteURL, p))
				found++
			}
		}
		if config.Amount == 0 {
			return URLs
		}
		if pageParam == "&pid=%d" {
			i += 42
			continue
		}
		i++
	}

	return URLs
}

func extractData(URL string) (*static.Data, error) {

	siteName := reSiteName.FindStringSubmatch(siteURL)[1]

	postHTML, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	matchedPostURL := rePostURL.FindStringSubmatch(postHTML)
	if len(matchedPostURL) != 3 {
		matchedPostURL = rePostBackup.FindStringSubmatch(postHTML)
		if len(matchedPostURL) != 3 {
			return nil, static.ErrDataSourceParseFailed
		}
	}

	if !strings.HasPrefix(matchedPostURL[1], "https") {
		matchedPostURL[1] = fmt.Sprintf("%s%s", "https:", matchedPostURL[1]) //tbib.org/ direct img link has no https:
	}

	var size int64
	if config.Amount == 0 {
		size, err = request.Size(matchedPostURL[1], URL)
		if err != nil {
			return nil, errors.New("no image size not found")
		}
	}

	id := reID.FindString(postHTML)

	quality := reSize.FindString(postHTML)
	quality = strings.ReplaceAll(quality, "Size: ", "")

	return &static.Data{
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
		URL: URL,
	}, nil

}

func extractDataFromDirectLink(URL string) (*static.Data, error) {
	matchedURL := reDirectLink.FindStringSubmatch(URL)
	if len(matchedURL) != 3 {
		return nil, errors.New("direct download can't match URL")
	}

	return &static.Data{
		Site:  siteURL,
		Title: matchedURL[1],
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: []*static.URL{
					{
						URL: URL,
						Ext: matchedURL[2],
					},
				},
				Size: 0,
			},
		},
		URL: URL,
	}, nil
}
