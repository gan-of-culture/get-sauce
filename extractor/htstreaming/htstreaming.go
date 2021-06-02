package htstreaming

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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

var site string

func ParseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`episode-\d+[/_\-]`, URL); ok {
		return []string{URL}
	}

	//check if it's an overview/series page maybe
	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https[^"\s]*?episode-\d*[/_]?[^"]*`)
	if strings.HasPrefix(URL, "https://hentaihaven.red") {
		re = regexp.MustCompile(`[^"]*red/hentai[^"]*`) //this sites URLs are built diff
	}
	matchedURLs := re.FindAllString(htmlString, -1)
	if strings.HasPrefix(URL, "https://hentaihaven.red") {
		//remove the five popular hentai on the side bar
		matchedURLs = matchedURLs[:len(matchedURLs)-5]
	}

	return removeAdjDuplicates(matchedURLs)
}

func Extract(URL string) ([]static.Data, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	site = baseURL.Host

	URLs := ParseURL(URL)
	if len(URLs) < 1 {
		log.Println(URL)
		return nil, fmt.Errorf("[%s] No matching URL found", site)
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := ExtractData(u)
		if err != nil {
			if strings.Contains(err.Error(), "Video not found") || strings.Contains(err.Error(), "PlayerURL not found") {
				log.Println(err.Error())
				continue
			}
			log.Println(u)
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

// ExtractData for a single episode that is hosted by the htstreaming network
func ExtractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		log.Println(htmlString)
		return static.Data{}, err
	}

	title := utils.GetMeta(&htmlString, "og:title")

	re := regexp.MustCompile(`[^"]*index.php\?data[^"]*`)
	playerURL := re.FindString(htmlString)
	if playerURL == "" {
		return static.Data{}, fmt.Errorf("[%s] PlayerURL not found %s", site, URL)
	}

	htmlString, err = request.Get(playerURL)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`({"[\s\S]*?), false`)
	jsonString := re.FindStringSubmatch(htmlString)
	if len(jsonString) < 2 {
		if strings.Contains(htmlString, "Video not found") {
			return static.Data{}, fmt.Errorf("[%s] Video not found %s", site, URL)
		}
		fmt.Println(htmlString)
		return static.Data{}, fmt.Errorf("[%s] JSON string no found %s", site, URL)
	}

	pData := pageData{}
	err = json.Unmarshal([]byte(jsonString[1]), &pData)
	if err != nil {
		log.Println(jsonString)
		log.Println(htmlString)
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
		// len(p.Variants)-i-1 builds stream map in reverse order
		// in order for the best quality stream to be on top
		streams[fmt.Sprintf("%d", len(p.Variants)-i-1)] = static.Stream{
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
		Url:     URL,
	}, nil
}

func removeAdjDuplicates(slice []string) []string {
	out := []string{}
	var last string
	for _, s := range slice {
		if s != last {
			out = append(out, s)
		}
		last = s
	}

	return out
}
