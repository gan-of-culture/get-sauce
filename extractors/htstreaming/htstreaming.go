package htstreaming

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
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

type extractor struct{}

// New returns a htstreaming extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	site = baseURL.Host

	URLs := parseURL(URL)
	if len(URLs) == 0 {
		log.Println(URL)
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := ExtractData(u)
		if err != nil {
			if strings.Contains(err.Error(), "video not found") || strings.Contains(err.Error(), "player URL not found") {
				log.Println(err.Error())
				continue
			}
			log.Println(u)
			return nil, err
		}
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
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

	return utils.RemoveAdjDuplicates(matchedURLs)
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
		return static.Data{}, errors.New("player URL not found %s")
	}

	htmlString, err = request.Get(playerURL)
	if err != nil {
		return static.Data{}, err
	}

	re = regexp.MustCompile(`({"[\s\S]*?), false`)
	jsonString := re.FindStringSubmatch(htmlString)
	if len(jsonString) < 2 {
		if strings.Contains(htmlString, "video not found") {
			return static.Data{}, errors.New("video not found %s")
		}
		fmt.Println(htmlString)
		return static.Data{}, errors.New("JSON string no found %s")
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

	// the meta tag contains bloat
	if site == "hentai.pro" {
		title = strings.TrimSuffix(pData.Title, "."+ext)
	}

	m3u8MasterURL := strings.Replace(pData.VideoData.VideoSources[0].File, pData.VideoServer, pData.HostList[pData.VideoServer][0], -1)
	m3u8MasterURL = strings.Replace(m3u8MasterURL, "hls", "down", -1)

	baseCDNURL := m3u8MasterURL[:len(m3u8MasterURL)-10] //remove master.txt

	m3u8Master, err := request.Get(m3u8MasterURL)
	if err != nil {
		return static.Data{}, err
	}

	dummyStreams, err := utils.ParseM3UMaster(&m3u8Master)
	if err != nil {
		return static.Data{}, err
	}

	streams := map[string]*static.Stream{}
	idx := 0
	for _, stream := range dummyStreams {
		streamURL := fmt.Sprintf("%s%s", baseCDNURL, stream.URLs[0].URL)

		master, err := request.Get(streamURL)
		if err != nil {
			return static.Data{}, err
		}

		URLs, key, err := request.GetM3UMeta(&master, streamURL, ext)
		if err != nil {
			return static.Data{}, err
		}

		// len(p.Variants)-i-1 builds stream map in reverse order
		// in order for the best quality stream to be on top
		streams[fmt.Sprintf("%d", len(dummyStreams)-idx-1)] = &static.Stream{
			URLs:    URLs,
			Quality: stream.Quality,
			Size:    stream.Size,
			Info:    pData.Title,
			Ext:     ext,
			Key:     key,
		}
		idx += 1
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: streams,
		Url:     URL,
	}, nil
}
