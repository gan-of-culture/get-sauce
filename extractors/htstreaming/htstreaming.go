package htstreaming

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

var reTitle = regexp.MustCompile(`"title":"[^"]+Episode \d+`)
var reVideoURL = regexp.MustCompile(`https://htstreaming.com/video/([^"]*)`)
var rePlayerURL = regexp.MustCompile(`[^"]*index.php\?data[^"]*`)

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
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`episode-\d+[/_\-]*`, URL); ok {
		return []string{URL}
	}

	//check if it's an overview/series page maybe
	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https[^"\s]*?episode-\d*(?:/*|[-\w]*)"`)
	matchedURLs := re.FindAllString(htmlString, -1)

	out := []string{}
	for _, u := range matchedURLs {
		out = append(out, strings.Trim(u, `"`))
	}

	return utils.RemoveAdjDuplicates(out)
}

// ExtractData for a single episode that is hosted by the htstreaming network
func ExtractData(URL string) (*static.Data, error) {

	playerURL, err := getPlayerURL(&URL)
	if err != nil {
		return nil, err
	}

	URLValues := url.Values{}
	URLValues.Add("hash", strings.Split(playerURL, "data=")[1])

	res, err := request.Request(http.MethodPost, playerURL+"&do=getVideo", map[string]string{
		"Referer":          playerURL,
		"x-requested-with": "XMLHttpRequest",
	}, strings.NewReader(URLValues.Encode()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buffer, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	pData := videoData{}
	err = json.Unmarshal(buffer, &pData)
	if err != nil {
		log.Println(string(buffer))
		return nil, err
	}

	streams, err := request.ExtractHLS(pData.VideoSource, map[string]string{
		"Referer": playerURL,
		"Accept":  "*/*",
	})
	if err != nil {
		return nil, err
	}

	ext := "ts"
	for _, stream := range streams {
		if strings.Contains(stream.Info, "mp4a") {
			ext = "mp4"
		}

		stream.Ext = ext
	}

	htmlString, err := request.Get(playerURL)
	if err != nil {
		return nil, err
	}

	matchedSubtitleParams := reSubtitleParams.FindStringSubmatch(htmlString) //1=jsTemplate 2=a 3=c 4=keywords
	if len(matchedSubtitleParams) < 5 {
		return nil, static.ErrDataSourceParseFailed
	}

	a, err := strconv.Atoi(matchedSubtitleParams[2])
	if err != nil {
		return nil, err
	}
	c, err := strconv.Atoi(matchedSubtitleParams[3])
	if err != nil {
		return nil, err
	}

	jsParams := parseFirePlayerParams(matchedSubtitleParams[1], a, c, strings.Split(matchedSubtitleParams[4], "|"))
	title := utils.GetLastItemString(strings.Split(reTitle.FindString(jsParams), `"`))
	title = strings.TrimSuffix(title, ".mkv")
	title = strings.TrimSuffix(title, ".mp4")

	return &static.Data{
		Site:     site,
		Title:    title,
		Type:     static.DataTypeVideo,
		Streams:  streams,
		Captions: parseCaptions(jsParams),
		URL:      URL,
	}, nil
}

// parsePlayerURL parses either a full URL
// like "https://htstreaming.com/player/index.php?data=" + hash or by parsing the hash and then
// returning a full URL.
func parsePlayerURL(target *string) string {
	playerURL := rePlayerURL.FindString(*target)
	if playerURL == "" {
		hash := utils.GetLastItemString(reVideoURL.FindStringSubmatch(*target))
		if hash != "" {
			playerURL = "https://htstreaming.com/player/index.php?data=" + hash
		}
	}
	return playerURL
}

// getPlayerURL returns a full and vaild htstreaming playerURL
// like "https://htstreaming.com/player/index.php?data=" + hash
func getPlayerURL(URL *string) (string, error) {
	var playerURL string
	if playerURL = parsePlayerURL(URL); playerURL != "" {
		return playerURL, nil
	}

	htmlString, err := request.Get(*URL)
	if err != nil {
		log.Println(htmlString)
		return "", err
	}

	htmlString = strings.ReplaceAll(htmlString, `\`, ``)

	if playerURL = parsePlayerURL(&htmlString); playerURL == "" {
		return "", errors.New("player URL not found")
	}

	return playerURL, nil
}
