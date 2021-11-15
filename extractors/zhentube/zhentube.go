package zhentube

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type playerData struct {
	HostList struct {
		Num1 []string `json:"1"`
	} `json:"hostList"`
	VideoURL    string      `json:"videoUrl"`
	VideoServer string      `json:"videoServer"`
	VideoDisk   interface{} `json:"videoDisk"`
	VideoPlayer string      `json:"videoPlayer"`
	IsJWPlayer8 bool        `json:"isJWPlayer8"`
	JwPlayerKey string      `json:"jwPlayerKey"`
	JwPlayerURL string      `json:"jwPlayerURL"`
	Logo        struct {
		File     string `json:"file"`
		Link     string `json:"link"`
		Position string `json:"position"`
		Hide     bool   `json:"hide"`
	} `json:"logo"`
	Tracks   []interface{} `json:"tracks"`
	Captions struct {
		FontSize   string `json:"fontSize"`
		Fontfamily string `json:"fontfamily"`
	} `json:"captions"`
	DefaultImage     string `json:"defaultImage"`
	SubtitleManager  bool   `json:"SubtitleManager"`
	Jwplayer8Button1 bool   `json:"jwplayer8button1"`
	Jwplayer8Quality bool   `json:"jwplayer8quality"`
	Title            string `json:"title"`
	Displaytitle     bool   `json:"displaytitle"`
	RememberPosition bool   `json:"rememberPosition"`
	Advertising      struct {
		Client string `json:"client"`
		Tag    string `json:"tag"`
	} `json:"advertising"`
	VideoData struct {
		VideoImage   interface{} `json:"videoImage"`
		VideoSources []struct {
			File  string `json:"file"`
			Label string `json:"label"`
			Type  string `json:"type"`
		} `json:"videoSources"`
	} `json:"videoData"`
}

const site = "https://zhentube.com/"

type extractor struct{}

// New returns a zhentube extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract data from URL
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	reURL := regexp.MustCompile(site + `[^/]+/$`)
	if reURL.MatchString(URL) {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}
	// remove rest of site's html including the sidebar that contains videos we don't want
	htmlString = strings.Split(htmlString, "<aside")[0]

	return regexp.MustCompile(site+`[^"]+episode-[^/]+/"`).FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	reEmbedURL := regexp.MustCompile(`<meta itemprop="embedURL" content="([^"]+)`)
	matchedEmbedURL := reEmbedURL.FindStringSubmatch(htmlString)
	if len(matchedEmbedURL) < 2 {
		return static.Data{}, static.ErrDataSourceParseFailed
	}

	title := utils.GetMeta(&htmlString, "og:title")

	htmlString, err = request.Get(matchedEmbedURL[1])
	if err != nil {
		return static.Data{}, err
	}

	rePlayerData := regexp.MustCompile(`{"hostList[^\n]+`)
	matchedPlayerData := rePlayerData.FindString(htmlString)
	if matchedPlayerData == "" {
		return static.Data{}, utils.Wrap(static.ErrDataSourceParseFailed, matchedEmbedURL[1])
	}
	//remove trailing js ", false);"
	matchedPlayerData = matchedPlayerData[:len(matchedPlayerData)-10]

	playerData := playerData{}
	err = json.Unmarshal([]byte(matchedPlayerData), &playerData)
	if err != nil {
		return static.Data{}, err
	}
	if len(playerData.VideoData.VideoSources) < 1 {
		return static.Data{}, utils.Wrap(static.ErrDataSourceParseFailed, "no videoSources in: "+matchedEmbedURL[1])
	}

	playerData.HostList.Num1[0] = "stream.deepthroatxvideo.com"
	masterURL := strings.Replace(playerData.VideoData.VideoSources[0].File, playerData.VideoServer, playerData.HostList.Num1[0], 1) + "?s=1&d="

	m3u8Master, err := request.GetWithHeaders(masterURL, map[string]string{
		"referer": matchedEmbedURL[1],
		"accept":  "*/*",
	})
	if err != nil {
		return static.Data{}, err
	}

	dummyStreams, err := utils.ParseM3UMaster(&m3u8Master)
	if err != nil {
		return static.Data{}, err
	}

	idx := 0
	streams := map[string]*static.Stream{}
	for _, s := range dummyStreams {
		idx += 1
		m3u8Media, err := request.GetWithHeaders(s.URLs[0].URL, map[string]string{
			"referer": matchedEmbedURL[1],
			"accept":  "*/*",
		})
		if err != nil {
			return static.Data{}, err
		}

		URLs, _, err := request.GetM3UMeta(&m3u8Media, s.URLs[0].URL, "mp4")
		if err != nil {
			return static.Data{}, err
		}

		streams[fmt.Sprint(len(dummyStreams)-idx)] = &static.Stream{
			URLs:    URLs,
			Quality: s.Quality,
			Size:    s.Size,
			Info:    s.Info,
			Ext:     "mp4",
		}
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    static.DataTypeVideo,
		Streams: streams,
		Url:     URL,
	}, nil
}
