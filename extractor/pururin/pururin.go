package pururin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://pururin.io/"
const apiGallery = "https://pururin.io/api/gallery"
const cdn = "https://cdn.pururin.io/assets/images/data/%s/%d.jpg"

func ParseURL(URL string) []string {
	if ok, _ := regexp.MatchString(fmt.Sprintf("%sgallery/\\d*/[^\"]*", site), URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(fmt.Sprintf("%sgallery/\\d*/[^\"]*\">", site))

	URLs := []string{}
	for _, v := range re.FindAllString(htmlString, -1) {
		URLs = append(URLs, strings.TrimSuffix(v, ">"))
	}
	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Pururin] No scrapable URL found for %s", URL)
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(fmt.Sprintf("%sgallery/(\\d*)/[^\"]*", site))
	ID := utils.GetLastItemString(re.FindStringSubmatch(URL))

	re = regexp.MustCompile(`(\d+) \( [\d.]+ \w \)`)
	matchedPageInfo := re.FindStringSubmatch(htmlString) //1=numPages
	pages, _ := strconv.Atoi(matchedPageInfo[1])

	URLs := []static.URL{}
	for _, pageNumber := range utils.NeedDownloadList(pages) {
		URLs = append(URLs, static.URL{
			URL: fmt.Sprintf(cdn, ID, pageNumber),
			Ext: "jpg",
		})
	}

	return static.Data{
		Site:  site,
		Title: strings.Split(utils.GetMeta(&htmlString, "og:title"), "/")[0],
		Type:  "image/jpg",
		Streams: map[string]static.Stream{
			"0": {
				URLs:    URLs,
				Quality: "unknown",
				Info:    ID,
			},
		},
		Url: URL,
	}, nil
}
