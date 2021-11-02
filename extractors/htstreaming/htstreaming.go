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

	"github.com/gan-of-culture/get-sauce/downloader"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
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

	re := regexp.MustCompile(`https[^"\s]*?episode-\d*/*"`)
	matchedURLs := re.FindAllString(htmlString, -1)

	out := []string{}
	for _, u := range matchedURLs {
		out = append(out, strings.Trim(u, `"`))
	}

	return utils.RemoveAdjDuplicates(out)
}

// ExtractData for a single episode that is hosted by the htstreaming network
func ExtractData(URL string) (static.Data, error) {
	title := ""

	re := regexp.MustCompile(`[^"]*index.php\?data[^"]*`)
	playerURL := re.FindString(URL)
	if playerURL == "" {

		htmlString, err := request.Get(URL)
		if err != nil {
			log.Println(htmlString)
			return static.Data{}, err
		}

		htmlString = strings.ReplaceAll(htmlString, `\`, ``)

		re = regexp.MustCompile(`[^"]*index.php\?data[^"]*`)
		playerURL = re.FindString(htmlString)
		if playerURL == "" {
			re = regexp.MustCompile(`https://htstreaming.com/video/([^"]*)`)
			hash := utils.GetLastItemString(re.FindStringSubmatch(htmlString))
			if hash != "" {
				playerURL = "https://htstreaming.com/player/index.php?data=" + hash
			}
		}
		if playerURL == "" {
			return static.Data{}, errors.New("player URL not found")
		}

		title = strings.Split(utils.GetMeta(&htmlString, "og:title"), " - ")[0]
		title = strings.Split(title, " | ")[0]

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
		return static.Data{}, err
	}

	m3u8Master, err := request.GetWithHeaders(pData.VideoSource, map[string]string{
		"referer": playerURL,
		"accept":  "*/*",
	})
	if err != nil {
		return static.Data{}, err
	}

	dummyStreams, err := utils.ParseM3UMaster(&m3u8Master)
	if err != nil {
		return static.Data{}, err
	}

	sortedStreams := downloader.GenSortedStreams(dummyStreams)

	ext := "ts"
	streams := map[string]*static.Stream{}
	for idx, stream := range sortedStreams {
		master, err := request.GetWithHeaders(stream.URLs[0].URL, map[string]string{
			"referer": playerURL,
			"accept":  "*/*",
		})
		if err != nil {
			return static.Data{}, err
		}

		URLs, key, err := request.GetM3UMeta(&master, stream.URLs[0].URL, ext)
		if err != nil {
			return static.Data{}, err
		}

		if strings.Contains(stream.Info, "mp4a") {
			ext = "mp4"
		}

		// len(p.Variants)-i-1 builds stream map in reverse order
		// in order for the best quality stream to be on top
		streams[fmt.Sprint(len(sortedStreams)-idx-1)] = &static.Stream{
			URLs:    URLs,
			Quality: stream.Quality,
			Size:    stream.Size,
			Ext:     ext,
			Key:     key,
		}
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: streams,
		Url:     URL,
	}, nil
}
