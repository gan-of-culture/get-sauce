package underhentai

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
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
	doc := soup.HTMLParse(htmlString)
	articles := doc.FindAll("article", "class", "data-block")
	if len(articles) == 0 {
		return nil, errors.New("[Underhentai] No content found")
	}

	content := []string{}
	for _, article := range articles {
		imgTag := article.Find("a")
		if imgTag.Error != nil {
			continue
		}
		content = append(content, imgTag.Attrs()["href"])
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

	doc := soup.HTMLParse(htmlString)
	tableOfStreams := doc.Find("div", "class", "content-table")
	if tableOfStreams == (soup.Root{}) {
		return nil, errors.New("[Underhentai] no table of streams found")
	}

	var episode string
	var streamError error

	data := []static.Data{}
	streams := make(map[string]static.Stream)

	tableRows := tableOfStreams.Find("tbody").FindAll("tr")
	for rowIdx, row := range tableRows {

		children := row.Children()
		if len(children) > 1 {

			//if tr includes stream info add stream info
			//if tr = images then add data and the streams for the data
			rowIdxString := strconv.Itoa(rowIdx)
			stream := static.Stream{
				URLs: make([]static.URL, 1),
			}

			for idx, child := range children {
				//only used for the switch statement thats why +1
				childIdxString := strconv.Itoa(idx + 1)
				switch fmt.Sprintf("c%s", childIdxString) {
				case "c1":
					episode = child.Text()
				case "c2":
					stream.URLs[0].Ext = child.Text()
				case "c3":
					size, _ := strconv.ParseInt(child.Text(), 10, 64)
					stream.Size = size
				case "c4":
					//idx + 1 to conteract that idx starts at 0
					if math.Mod(float64(rowIdx+1), 2) == 0.0 {
						stream.Info = "Has subtitles"
					}
				case "c5":
					//audio idk if it's important - implement if needed
				case "c6":
					stream.Info = fmt.Sprintf("%s isCensored %s", stream.Info, child.Text())
				case "c7":
					childChildren := child.Children()
					if len(childChildren) == 0 {
						streamError = errors.New("[Underhentai] no torrent found " + episode)
					}
					torrentChild := childChildren[0]
					re := regexp.MustCompile("id=([0-9]+)")
					matchedID := re.FindStringSubmatch(torrentChild.Attrs()["href"])
					if len(matchedID) == 1 {
						streamError = errors.New("[Underhentai] no id found " + episode)
					}
					id := matchedID[1]

					// remove prefixed zeros
					epForURL := strings.TrimPrefix(episode, "0")

					html, err := request.Get(fmt.Sprintf("https://www.underhentai.net/out/?sv=nya&id=%s&ep=%s", id, epForURL))
					if err != nil {
						streamError = errors.New("[Underhentai] no .torrent found " + episode)
					}

					re = regexp.MustCompile("https://.+.torrent")
					stream.URLs[0].URL = re.FindString(html)
				}
			}
			if streamError == nil {
				streams[rowIdxString] = stream
			}

		} else if len(children) == 1 {

			//add new data for each new episode
			data = append(data, static.Data{
				Site:    site,
				Title:   fmt.Sprintf("%s episode %s", title, episode),
				Type:    "video",
				Streams: streams,
				Err:     streamError,
				Url:     URL,
			})

		}

	}

	return data, nil

}
