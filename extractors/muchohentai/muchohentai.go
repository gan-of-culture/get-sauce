package muchohentai

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type m3u8File struct {
	URL string `json:"file"`
}

const site = "https://muchohentai.com/"

type extractor struct{}

// New returns a muchohentai extractor.
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
			return nil, err
		}
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"aBo4Rk/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`[^"]*aBo4Rk/\d*/">\n`)
	URLs := []string{}
	last := ""
	for _, v := range re.FindAllString(htmlString, -1) {
		if v == last {
			last = v
			continue
		}
		URLs = append(URLs, strings.TrimSuffix(v, "\"/>\n"))
		last = v
	}
	return URLs
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}
	title := utils.GetMeta(&htmlString, "og:title")

	re := regexp.MustCompile(`var servers=\[([^\]]*).*?var server="([^;]*)";var files=\[([^\]]*)`)
	serverInfo := re.FindStringSubmatch(htmlString) //1=servers ('va01','va02','va03','va04') 2=server URL template https://"+servers[choice-1]+"-edge.tmncdn.io" 3=master file URL {"file":"\/wp-content\/uploads\/Soshite_Watashi\/episode_3\/ja.m3u8"}
	if len(serverInfo) < 4 {
		return static.Data{}, fmt.Errorf("cannot extract server info for %s", URL)
	}

	//va01 maybe honeypot? logs ip in response header? doesn't resolve to a segment key?
	//remove it the other mirrors should be fine enough
	servers := strings.Split(strings.ReplaceAll(serverInfo[1], "'", ""), ",")[1:]
	m3u8FileJson := &m3u8File{}
	err = json.Unmarshal([]byte(serverInfo[3]), m3u8FileJson)
	if err != nil {
		return static.Data{}, fmt.Errorf("cannot extract file URL for %s", URL)
	}

	masterURL := strings.Replace(serverInfo[2], "\"+servers[choice-1]+\"", servers[0], 1) + m3u8FileJson.URL
	m3u8String, err := request.Get(masterURL)
	if err != nil {
		return static.Data{}, err
	}

	dummyStreams, err := utils.ParseM3UMaster(&m3u8String)
	if err != nil {
		URLs, key, err := request.GetM3UMeta(&m3u8String, masterURL, "ts")
		if err != nil {
			return static.Data{}, err
		}
		return static.Data{
			Site:  site,
			Title: title,
			Type:  "video",
			Streams: map[string]*static.Stream{
				"0": {
					URLs:    URLs,
					Quality: "best",
					Info:    masterURL,
					Ext:     "ts",
					Key:     key,
				},
			},
			Url: URL,
		}, nil
	}

	streams := map[string]*static.Stream{}
	streamIdx := 0
	vOLD := servers[0]
	for _, v := range servers {
		masterURL = strings.Replace(masterURL, vOLD, v, 1) //servers have mirrored content so we send a request once to reduce traffic
		vOLD = v
		baseURL, err := url.Parse(masterURL)
		if err != nil {
			return static.Data{}, err
		}
		for _, variant := range dummyStreams {
			streamURL, err := baseURL.Parse(variant.URLs[0].URL)
			if err != nil {
				return static.Data{}, err
			}

			master, err := request.Get(streamURL.String())
			if err != nil {
				return static.Data{}, err
			}

			URLs, key, err := request.GetM3UMeta(&master, streamURL.String(), "ts")
			if err != nil {
				return static.Data{}, err
			}

			streams[strconv.Itoa(streamIdx)] = &static.Stream{
				URLs:    URLs,
				Quality: variant.Quality,
				Size:    variant.Size,
				Info:    variant.Info, //contains codec
				Ext:     "ts",
				Key:     key,
			}
			streamIdx++
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
