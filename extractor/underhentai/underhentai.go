package underhentai

/*import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

const site = "https://underhentai.net"

// Extract to extract data form underthentai.net
func Extract(URL string) ([]static.Data, error) {
	URLs, err := ParseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []static.Data{}
	for _, parsedURL := range URLs {
		extractedData, err := extractData(parsedURL)
		if err != nil {
			return nil, err
		}
		data = append(data, extractedData...)
	}

	return data, nil
}

// ParseURL extractable URL
func ParseURL(URL string) ([]string, error) {

	re := regexp.MustCompile("(/releases-)|(/tag/)|(/index/)|(/uncensored/)|(/top/)")
	matches := re.FindStringSubmatch(URL)
	// if not a match -> probably URL of some hentai
	if len(matches) == 0 {
		return []string{URL}, nil
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	re = regexp.MustCompile("<[\t\n\f\r ]*div class=\"article[^>]*>[\t\n\f\r ]<a href=\"([^\"]+)\".*?<[\t\n\f\r ]*\/[\t\n\f\r ]*div>")
	matchedTitles := re.FindAllStringSubmatch(htmlString, -1)
	if len(matchedTitles) == 0 {
		return nil, errors.New("[Underhentai] No content found")
	}

	content := []string{}
	for _, title := range matchedTitles {
		content = append(content, title[1])
	}

	return content, nil
}

func extractData(URL string) ([]static.Data, error) {
	re := regexp.MustCompile("net/(.+)/")
	matches := re.FindStringSubmatch(URL)
	if len(matches) == 0 {
		return nil, errors.New("[Underhentai] can't parse URL of content")
	}
	title := matches[1]

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	var episode string
	var extension string
	var size int64
	var hasSubtitles bool
	var isCensored bool
	var torrentURL string

	currentEp := "01"

	streams := map[string]static.Stream{}
	data := []static.Data{}

	re = regexp.MustCompile("<td class=\"c([0-9])\">(.*?)</td>")
	matchedData := re.FindAllStringSubmatch(htmlString, -1)
	for _, dataElement := range matchedData {
		switch dataElement[1] {
		case "1":
			episode = dataElement[2]
			if episode != currentEp {
				data = append(data, static.Data{
					Site:    site,
					Title:   fmt.Sprintf("%s episode %s", title, currentEp),
					Type:    "video",
					Streams: streams,
					Url:     URL,
				})
				streams = map[string]static.Stream{}
			}
			currentEp = episode
		case "2":
			extension = dataElement[2]
		case "3":
			size, _ = strconv.ParseInt(dataElement[2], 10, 64)
		case "4":
			if strings.Contains(dataElement[2], "/img/xnone.png") {
				hasSubtitles = true
				break
			}
			hasSubtitles = false
		case "5":
			// not implemented - not needed yet
		case "6":
			if dataElement[2] == "Yes" {
				isCensored = true
				break
			}
			isCensored = false
		case "7":
			ep := strconv.Itoa(len(streams))
			re = regexp.MustCompile(fmt.Sprintf("sv=nya&id=([0-9]*)&ep=%s", ep))
			torrentURLSuffix := re.FindStringSubmatch(htmlString)
			if len(torrentURLSuffix) == 0 {
				// stream with no bittorrent
				streams[fmt.Sprintf("%d", len(streams))] = static.Stream{}
				continue
			}

			html, err := request.Get("https://www.underhentai.net/out/?" + torrentURLSuffix[0])
			if err != nil {
				log.Println(errors.New("[Underhentai] no .torrent found " + episode))
			}

			re = regexp.MustCompile("url=\"(https://.+.torrent)")
			matchedTorrentURL := re.FindStringSubmatch(html)
			if len(matchedTorrentURL) != 0 {
				torrentURL = matchedTorrentURL[1]
			}

			streams[fmt.Sprintf("%d", len(streams))] = static.Stream{
				URLs: []static.URL{
					{
						URL: torrentURL,
						Ext: extension,
					},
				},
				Quality: "?",
				Size:    size,
				Info:    fmt.Sprintf("hasSubtitles: %t isCensored: %t", hasSubtitles, isCensored),
			}
		case "8":
			/*re = regexp.MustCompile("href=\"(.*?)\"")
			watchURL := re.FindStringSubmatch(dataElement[2])
			if len(watchURL) != 2 {
				// stream with no bittorrent
				streams[fmt.Sprintf("%d", len(streams))] = static.Stream{}
				continue
			}

		}
	}

	data = append(data, static.Data{
		Site:    site,
		Title:   fmt.Sprintf("%s episode %s", title, episode),
		Type:    "video",
		Streams: streams,
		Url:     URL,
	})

	return data, nil

}*/
