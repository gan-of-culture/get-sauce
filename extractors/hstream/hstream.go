package hstream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	mpegdash "github.com/gan-of-culture/get-sauce/parsers/mpeg_dash"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"golang.org/x/exp/slices"
)

const site = "https://hstream.moe/"
const api = "https://hstream.moe/player/api"
const fileProvider = "https://str.h-dl.xyz"

type APIResponse struct {
	Title      string `json:"title"`
	Poster     string `json:"poster"`
	Legacy     int    `json:"legacy"`
	Resolution string `json:"resolution"`
	StreamURL  string `json:"stream_url"`
}

type APIPayload struct {
	EpisodeID string `json:"episode_id"`
}

var reEpisodeID = regexp.MustCompile(`e_id" type="hidden" value="([^"]*)`)
var reCaptionSource = regexp.MustCompile(`https://.+/\d+/[\w.]+/[\w./]+\.ass`)

type extractor struct{}

// New returns a hstream extractor
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)

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
	if ok, _ := regexp.MatchString(site+`hentai/[\w\-]+/\d+`, URL); ok {
		return []string{URL}
	}

	if ok, _ := regexp.MatchString(site+`hentai/[\w\-]+/?`, URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	htmlString = strings.Split(htmlString, `class="bixbox"`)[0]

	re := regexp.MustCompile(`hentai/[\w\-]+/\d+`)
	out := []string{}
	for _, episode := range re.FindAllString(htmlString, -1) {
		out = append(out, site+episode)
	}
	return utils.RemoveAdjDuplicates(out)
}

func extractData(URL string) (*static.Data, error) {
	resp, err := request.Request(http.MethodGet, URL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	htmlString := string(body)
	cookies := resp.Cookies()
	xsrf := cookies[slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "XSRF-TOKEN"
	})]
	hstreamSession := cookies[slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == "hstream_session"
	})]

	if strings.Contains(htmlString, "<title>DDOS-GUARD</title>") {
		time.Sleep(300 * time.Millisecond)
		htmlString, _ = request.Get(URL)
	}

	matchedEpisodeID := reEpisodeID.FindStringSubmatch(htmlString)
	if len(matchedEpisodeID) < 2 {
		return nil, errors.New("cannot find e_id for")
	}

	payload := APIPayload{EpisodeID: strings.TrimSpace(matchedEpisodeID[1])}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	xsrfValueEscaped, err := url.PathUnescape(xsrf.Value)
	if err != nil {
		return nil, err
	}

	jsonString, err := request.PostAsBytesWithHeaders(api, map[string]string{
		"Content-Length":   fmt.Sprint(len(payloadBytes)),
		"Content-Type":     "application/json",
		"Cookie":           fmt.Sprintf("%s=%s; %s=%s", xsrf.Name, xsrf.Value, hstreamSession.Name, hstreamSession.Value),
		"Referer":          URL,
		"X-Requested-With": "XMLHttpRequest",
		"X-Xsrf-Token":     xsrfValueEscaped,
	}, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}

	res := APIResponse{}
	err = json.Unmarshal(jsonString, &res)
	if err != nil {
		return nil, err
	}

	videoSourceBaseURL := fmt.Sprintf("%s/%s", fileProvider, res.StreamURL)
	videoSources := []string{
		videoSourceBaseURL + "/x264.720p.mp4",
	}

	if res.Resolution == "1080p" {
		videoSources = append(videoSources, videoSourceBaseURL+"/av1.1080p.webm")
	}

	if res.Resolution == "4k" {
		videoSources = append(videoSources, videoSourceBaseURL+"/av1.1080p.webm")
		videoSources = append(videoSources, videoSourceBaseURL+"/av1.2160p.webm")
	}

	if res.Legacy == 0 {
		videoSources = append(videoSources, []string{
			videoSourceBaseURL + "/720/manifest.mpd",
			videoSourceBaseURL + "/1080/manifest.mpd",
			videoSourceBaseURL + "/2160/manifest.mpd",
		}...)
	}

	videoSources = reverse(videoSources)

	streams := make(map[string]*static.Stream)
	counter := 0
	// only keep one audio stream
	foundAudioStream := false
	for _, sourceURL := range videoSources {
		if !strings.HasSuffix(sourceURL, ".mpd") {
			continue
		}

		streamsTmp, err := mpegdash.ExtractDASHManifest(sourceURL, map[string]string{"Referer": site})
		if err != nil {
			return nil, err
		}
		for _, streamTmp := range streamsTmp {
			if streamTmp.Type == static.DataTypeAudio {
				if foundAudioStream {
					continue
				} else {
					foundAudioStream = true
				}
			}

			streams[fmt.Sprint(counter)] = streamTmp
			counter += 1
		}
	}

	// skip direct file downloads if newer mpd is supplied
	if len(streams) > 0 {
		videoSources = []string{}
	}

	for i, sourceURL := range videoSources {
		size, err := request.Size(sourceURL, site)
		if err != nil {
			return nil, err
		}

		streams[fmt.Sprint(i)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: sourceURL,
					Ext: utils.GetFileExt(sourceURL),
				},
			},
			Quality: utils.GetLastItemString(strings.Split(sourceURL, "/")),
			Size:    size,
		}
	}

	captionURL := videoSourceBaseURL + "/eng.ass"

	return &static.Data{
		Site:    site,
		Title:   strings.TrimSpace(utils.GetSectionHeadingElement(&htmlString, 1, 0)),
		Type:    static.DataTypeVideo,
		Streams: streams,
		Captions: []*static.Caption{
			{
				URL: static.URL{
					URL: captionURL,
					Ext: utils.GetFileExt(captionURL),
				},
				Language: "English",
			},
		},
		URL: URL,
	}, nil

}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
