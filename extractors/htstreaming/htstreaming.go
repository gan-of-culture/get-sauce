package htstreaming

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	Size  string `json:"size"`
}

type videoData struct {
	VideoImage    string   `json:"videoImage"`
	VideoSource   string   `json:"videoSource"`
	DownloadLinks []source `json:"downloadLinks"`
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
				log.Println(utils.Wrap(err, u).Error())
				continue
			}
			return nil, utils.Wrap(err, u)
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

	title := strings.Split(utils.GetMeta(&htmlString, "og:title"), " - ")[0]
	title = strings.Split(title, " | ")[0]

	re := regexp.MustCompile(`[^"]*index.php\?data[^"]*`)
	playerURL := re.FindString(htmlString)
	if playerURL == "" {
		return static.Data{}, errors.New("player URL not found %s")
	}

	URLValues := url.Values{}
	URLValues.Add("hash", strings.Split(playerURL, "data=")[1])

	res, err := request.Request(http.MethodPost, playerURL+"&do=getVideo", map[string]string{
		"Referer":          playerURL,
		"x-requested-with": "XMLHttpRequest",
	}, strings.NewReader(URLValues.Encode()))
	if err != nil {
		return static.Data{}, err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return static.Data{}, err
	}

	pData := videoData{}
	err = json.Unmarshal(buffer, &pData)
	if err != nil {
		log.Println(string(buffer))
		log.Println(htmlString)
		return static.Data{}, err
	}

	m3u8MasterRes, err := request.Request(http.MethodGet, pData.VideoSource, map[string]string{
		"referer": playerURL,
		"accept":  "*/*",
	}, nil)
	if err != nil {
		return static.Data{}, err
	}
	defer m3u8MasterRes.Body.Close()

	buffer, err = ioutil.ReadAll(m3u8MasterRes.Body)
	if err != nil {
		return static.Data{}, err
	}

	m3u8Master := string(buffer)

	dummyStreams, err := utils.ParseM3UMaster(&m3u8Master)
	if err != nil {
		return static.Data{}, err
	}

	ext := "ts"
	streams := map[string]*static.Stream{}
	idx := 0
	for _, stream := range dummyStreams {
		masterRes, err := request.Request(http.MethodGet, stream.URLs[0].URL, map[string]string{
			"referer": playerURL,
			"accept":  "*/*",
		}, nil)
		if err != nil {
			return static.Data{}, err
		}
		defer masterRes.Body.Close()

		buffer, err = ioutil.ReadAll(masterRes.Body)
		if err != nil {
			return static.Data{}, err
		}

		master := string(buffer)

		URLs, key, err := request.GetM3UMeta(&master, stream.URLs[0].URL, ext)
		if err != nil {
			return static.Data{}, err
		}

		if strings.Contains(stream.Info, "mp4a") {
			ext = "mp4"
		}

		// len(p.Variants)-i-1 builds stream map in reverse order
		// in order for the best quality stream to be on top
		streams[fmt.Sprintf("%d", len(dummyStreams)-idx-1)] = &static.Stream{
			URLs:    URLs,
			Quality: stream.Quality,
			Size:    stream.Size,
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
