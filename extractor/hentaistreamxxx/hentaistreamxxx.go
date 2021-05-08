package hentaistreamxxx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
	"github.com/grafov/m3u8"
)

type source struct {
	File  string `json:"file"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type videoData struct {
	VideoImage   string   `json:"videoImage"`
	VideoSources []source `json:"videoSources"`
}

type pageData struct {
	HostList    map[string][]string `json:"hostList"`
	VideoURL    string              `json:"videoUrl"`
	VideoServer string              `json:"videoServer"`
	Title       string              `json:"title"`
	VideoData   videoData           `json:"videoData"`
}

const site = "https://hentaistream.xxx"

func ParseURL(URL string) []string {
	if strings.HasPrefix(URL, "https://hentaistream.xxx/watch/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://hentaistream.xxx/videos/category/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`[^"]*xxx/watch[^"]*`)
	matchedURLs := re.FindAllString(htmlString, -1)

	URLs := []string{}
	for i := range matchedURLs {
		if i%2 == 0 {
			URLs = append(URLs, fmt.Sprintf("%s%s", site, strings.TrimPrefix(matchedURLs[i], "/")))
		}
	}

	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) < 1 {
		return nil, fmt.Errorf("[Hanime] No matching URL found.")
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
	title := utils.GetMeta(htmlString, "og:title")

	re := regexp.MustCompile(`[^"]*php\?data[^"]*`)
	playerURL := re.FindString(htmlString)

	htmlString, err = request.Get(playerURL)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`({"[\s\S]*?), `)
	jsonString := re.FindStringSubmatch(htmlString)
	if len(jsonString) < 2 {
		return static.Data{}, fmt.Errorf("[HentaistreamXXX] JSON string no found")
	}

	pData := pageData{}
	err = json.Unmarshal([]byte(jsonString[1]), &pData)
	if err != nil {
		log.Println(URL)
		log.Println(jsonString)
		return static.Data{}, err
	}

	ext := "ts"
	re = regexp.MustCompile(`\.([\d\w]*)$`)
	matchedExt := re.FindStringSubmatch(pData.Title)
	if len(matchedExt) >= 2 {
		ext = matchedExt[1]
	}

	m3u8MasterURL := strings.Replace(pData.VideoData.VideoSources[0].File, pData.VideoServer, pData.HostList[pData.VideoServer][0], -1)
	m3u8MasterURL = strings.Replace(m3u8MasterURL, "hls", "down", -1)

	baseCDNURL := m3u8MasterURL[:len(m3u8MasterURL)-10] //remove master.txt

	m3u8Master, err := request.Request(http.MethodGet, m3u8MasterURL, nil)
	if err != nil {
		return static.Data{}, err
	}

	p := m3u8.NewMasterPlaylist()
	err = p.DecodeFrom(m3u8Master.Body, false)
	if err != nil {
		fmt.Println(err)
	}

	streams := map[string]static.Stream{}
	for i, stream := range p.Variants {
		streams[fmt.Sprintf("%d", i)] = static.Stream{
			URLs: []static.URL{
				{
					URL: fmt.Sprintf("%s%s", baseCDNURL, stream.URI),
					Ext: ext,
				},
			},
			Quality: stream.Resolution,
			Size:    int64(stream.Bandwidth),
			Info:    pData.Title,
		}
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "application/x-mpegurl",
		Streams: streams,
		Err:     nil,
		Url:     URL,
	}, nil
}
