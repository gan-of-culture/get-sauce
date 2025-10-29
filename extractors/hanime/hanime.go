package hanime

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/parsers/hls"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type stream struct {
	ID          int    `json:"id"`
	Slug        string `json:"slug"`
	Extension   string `json:"extension"`
	MimeType    string `json:"mime_type"`
	Width       int    `json:"width"`
	Height      string `json:"height"`
	FilesizeMbs int    `json:"filesize_mbs"`
	Filename    string `json:"filename"`
	URL         string `json:"url"`
}

type videoData struct {
	Error          any  `json:"error"`
	ServerRendered bool `json:"serverRendered"`
	State          struct {
		Data struct {
			Video struct {
				PlayerBaseURL string `json:"player_base_url"`
				HentaiVideo   struct {
					ID        int    `json:"id"`
					IsVisible bool   `json:"is_visible"`
					Name      string `json:"name"`
					Slug      string `json:"slug"`
				} `json:"hentai_video"`
				VideosManifest struct {
					Servers []struct {
						ID          int      `json:"id"`
						Name        string   `json:"name"`
						Slug        string   `json:"slug"`
						NaRating    int      `json:"na_rating"`
						EuRating    int      `json:"eu_rating"`
						AsiaRating  int      `json:"asia_rating"`
						Sequence    int      `json:"sequence"`
						IsPermanent bool     `json:"is_permanent"`
						Streams     []stream `json:"streams"`
					} `json:"servers"`
				} `json:"videos_manifest"`
			} `json:"video"`
		} `json:"data"`
	} `json:"state"`
}

const site = "https://hanime.tv/"

var reNuxtState = regexp.MustCompile(`{"layout"[^<]+`)

type extractor struct{}

// New returns a hanime.tv extractor.
func New() static.Extractor {
	return &extractor{}
}

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
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if strings.HasPrefix(URL, "https://hanime.tv/videos/hentai/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	re := regexp.MustCompile(`videos/hentai[^"]*`)
	out := []string{}
	for _, URLPart := range re.FindAllString(htmlString, -1) {
		out = append(out, site+URLPart)
	}
	return out
}

func extractData(URL string) (*static.Data, error) {

	htmlData, err := request.GetAsBytes(URL)
	if err != nil {
		return nil, err
	}

	htmlData = reNuxtState.Find(htmlData)

	vData := videoData{}
	err = json.Unmarshal(htmlData[:len(htmlData)-1], &vData)
	if err != nil {
		return nil, err
	}

	// remove first entry if it's the 1080p stream since it only works if you are logged in
	if vData.State.Data.Video.VideosManifest.Servers[0].Streams[0].Height == "1080" {
		vData.State.Data.Video.VideosManifest.Servers[0].Streams = remove(vData.State.Data.Video.VideosManifest.Servers[0].Streams, 0)
	}

	streams := map[string]*static.Stream{}
	for idx, streamData := range vData.State.Data.Video.VideosManifest.Servers[0].Streams {
		mediaStr, err := request.Get(streamData.URL)

		URLs, key, err := hls.ParseMediaStream(&mediaStr, site)
		if err != nil {
			return nil, err
		}

		streams[fmt.Sprint(idx)] = &static.Stream{
			Type:    static.DataTypeVideo,
			URLs:    URLs,
			Quality: fmt.Sprintf("%sp; %d x %s", streamData.Height, streamData.Width, streamData.Height),
			Size:    utils.CalcSizeInByte(float64(streamData.FilesizeMbs), "MB"),
			Key:     key,
			Ext:     "mp4",
		}
	}

	return &static.Data{
		Site:    site,
		Title:   vData.State.Data.Video.HentaiVideo.Name,
		Type:    static.DataTypeVideo,
		Streams: streams,
		URL:     URL,
	}, nil
}

func remove(slice []stream, s int) []stream {
	return append(slice[:s], slice[s+1:]...)
}
