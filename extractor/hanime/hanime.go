package hanime

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type stream struct {
	ID                 float64 `json:"id,omitempty"`
	ServerID           float64 `json:"server_id,omitempty"`
	Slug               string  `json:"slug,omitempty"`
	Kind               string  `json:"kind,omitempty"`
	Extension          string  `json:"extension,omitempty"`
	MimeType           string  `json:"mime_type"`
	Width              float64 `json:"width,omitempty"`
	Height             string  `json:"height,omitempty"`
	DurationInMs       float64 `json:"duration_in_ms"`
	FileSizeInMBs      float64 `json:"filesize_mbs"`
	Filename           string  `json:"filename,omitempty"`
	URL                string  `json:"url,omitempty"`
	IsGuestAllowed     bool    `json:"is_guest_allowed,omitempty"`
	IsMemeberAllowed   bool    `json:"is_member_allowed,omitempty"`
	IsPremiumAllowed   bool    `json:"is_premium_allowed,omitempty"`
	IsDownloadable     bool    `json:"is_downloadable,omitempty"`
	Compatibility      string  `json:"compatibility,omitempty"`
	HvID               float64 `json:"hv_id,omitempty"`
	HostID             float64 `json:"host_id,omitempty"`
	SubDomain          string  `json:"sub_domain,omitempty"`
	ServerSequence     float64 `json:"server_sequence,omitempty"`
	VideoStreamGroupID string  `json:"video_stream_group_id,omitempty"`
}

type server struct {
	ID          float64  `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Slug        string   `json:"slug,omitempty"`
	NARating    float64  `json:"na_rating,omitempty"`
	EURating    float64  `json:"eu_rating,omitempty"`
	AsiaRating  float64  `json:"asia_rating,omitempty"`
	Sequence    float64  `json:"sequence,omitempty"`
	IsPermanent bool     `json:"is_permanent,omitempty"`
	Streams     []stream `json:"streams,omitempty"`
}

type videosManifest struct {
	Servers []server `json:"servers,omitempty"`
}

type hentaiVideo struct {
	Name string `json:"name"`
}

type video struct {
	HentaiVideo    hentaiVideo    `json:"hentai_video"`
	VideosManifest videosManifest `json:"videos_manifest"`
}

type data struct {
	Video video `json:"video"`
}

type state struct {
	Data data `json:"data"`
}

type pageData struct {
	State state `json:"state"`
}

const site = "https://hanime.tv/"

func ParseURL(URL string) []string {
	if strings.HasPrefix(URL, "https://hanime.tv/videos/hentai/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://hanime.tv/browse/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`[^"]*videos/hentai[^"]*`)
	matchedURLs := re.FindAllString(htmlString, -1)

	URLs := []string{}
	for _, matchedURL := range matchedURLs {
		URLs = append(URLs, fmt.Sprintf("%s%s", site, strings.TrimPrefix(matchedURL, "/")))
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

	re := regexp.MustCompile(`NUXT__=([\s\S]*?);<`)
	matchedJSONString := re.FindStringSubmatch(htmlString)
	if len(matchedJSONString) < 2 {
		return static.Data{}, fmt.Errorf("[Hanime] JSON string no found")
	}

	jsonString := strings.Trim(matchedJSONString[1], ",")

	pData := pageData{}
	err = json.Unmarshal([]byte(jsonString), &pData)
	if err != nil {
		fmt.Println(URL)
		fmt.Println(jsonString)
		return static.Data{}, err
	}

	streams := map[string]static.Stream{}
	for _, stream := range pData.State.Data.Video.VideosManifest.Servers[0].Streams {
		if stream.URL == "" {
			continue
		}

		streams[fmt.Sprintf("%d", len(streams))] = static.Stream{
			URLs: []static.URL{
				{
					URL: stream.URL,
					Ext: "ts",
				},
			},
			Quality: fmt.Sprintf("%v x %s", stream.Width, stream.Height),
			Size:    utils.CalcSizeInByte(stream.FileSizeInMBs, "MB"),
			Info:    stream.Filename,
		}
	}

	return static.Data{
		Site:    site,
		Title:   pData.State.Data.Video.HentaiVideo.Name,
		Type:    "application/x-mpegurl",
		Streams: streams,
		Err:     nil,
		Url:     URL,
	}, nil
}
