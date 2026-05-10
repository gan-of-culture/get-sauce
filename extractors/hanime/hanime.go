package hanime

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/gan-of-culture/get-sauce/parsers/hls"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
)

type videoData struct {
	Layout string `json:"layout"`
	Data   []struct {
	} `json:"data"`
	Fetch struct {
	} `json:"fetch"`
	Error          string `json:"error"`
	ServerRendered string `json:"serverRendered"`
	RoutePath      string `json:"routePath"`
	Config         struct {
		LandingEvent struct {
			Periods []struct {
				Start  time.Time `json:"start"`
				End    time.Time `json:"end"`
				ImgSrc string    `json:"img_src"`
				URL    string    `json:"url"`
			} `json:"periods"`
		} `json:"LANDING_EVENT"`
		CdnBaseURL   string `json:"CDN_BASE_URL"`
		EnvJSONURL   string `json:"ENV_JSON_URL"`
		SearchHvsURL string `json:"SEARCH_HVS_URL"`
		App          struct {
			BasePath   string `json:"basePath"`
			AssetsPath string `json:"assetsPath"`
			CdnURL     string `json:"cdnURL"`
		} `json:"_app"`
	} `json:"config"`
	State struct {
		Num3                         string `json:"3"`
		ScrollY                      string `json:"scrollY"`
		CsrfToken                    string `json:"csrf_token"`
		CsrfTokenLastFetchedTimeUnix string `json:"csrf_token_last_fetched_time_unix"`
		Version                      string `json:"version"`
		IsNewVersion                 string `json:"is_new_version"`
		CountryCode                  string `json:"country_code"`
		PageName                     string `json:"page_name"`
		UserAgent                    string `json:"user_agent"`
		IP                           string `json:"ip"`
		Referrer                     string `json:"referrer"`
		Geo                          string `json:"geo"`
		IsDev                        string `json:"is_dev"`
		IsWasmSupported              string `json:"is_wasm_supported"`
		IsMounted                    string `json:"is_mounted"`
		IsLoading                    string `json:"is_loading"`
		IsImageProcessing            string `json:"is_image_processing"`
		IsSearching                  string `json:"is_searching"`
		BrowserWidth                 string `json:"browser_width"`
		BrowserHeight                string `json:"browser_height"`
		SystemMsg                    string `json:"system_msg"`
		Data                         struct {
			Video struct {
				Num1056469103 string `json:"1056469103"`
				PlayerBaseURL string `json:"player_base_url"`
				HentaiVideo   struct {
					ID              string    `json:"id"`
					IsVisible       string    `json:"is_visible"`
					Name            string    `json:"name"`
					Slug            string    `json:"slug"`
					CreatedAt       time.Time `json:"created_at"`
					ReleasedAt      time.Time `json:"released_at"`
					Description     string    `json:"description"`
					Views           string    `json:"views"`
					Interests       string    `json:"interests"`
					PosterURL       string    `json:"poster_url"`
					CoverURL        string    `json:"cover_url"`
					IsHardSubtitled string    `json:"is_hard_subtitled"`
					Brand           string    `json:"brand"`
					DurationInMs    string    `json:"duration_in_ms"`
					IsCensored      string    `json:"is_censored"`
					Rating          string    `json:"rating"`
					Likes           string    `json:"likes"`
					Dislikes        string    `json:"dislikes"`
					Downloads       string    `json:"downloads"`
					MonthlyRank     string    `json:"monthly_rank"`
					BrandID         string    `json:"brand_id"`
					IsBannedIn      string    `json:"is_banned_in"`
					PreviewURL      string    `json:"preview_url"`
					PrimaryColor    string    `json:"primary_color"`
					CreatedAtUnix   string    `json:"created_at_unix"`
					ReleasedAtUnix  string    `json:"released_at_unix"`
					HentaiTags      []struct {
						ID   string `json:"id"`
						Text string `json:"text"`
					} `json:"hentai_tags"`
					Titles []struct {
						Lang  string `json:"lang"`
						Kind  string `json:"kind"`
						Title string `json:"title"`
					} `json:"titles"`
				} `json:"hentai_video"`
				VideosManifest struct {
					Servers []struct {
						ID          string `json:"id"`
						AsiaRating  string `json:"asia_rating"`
						EuRating    string `json:"eu_rating"`
						NaRating    string `json:"na_rating"`
						IsPermanent string `json:"is_permanent"`
						Name        string `json:"name"`
						Sequence    string `json:"sequence"`
						Slug        string `json:"slug"`
						Streams     []struct {
							ServerID           string `json:"server_id"`
							ID                 int    `json:"id"`
							Width              int    `json:"width"`
							Height             string `json:"height"`
							Compatibility      string `json:"compatibility"`
							DurationInMs       string `json:"duration_in_ms"`
							Extension          string `json:"extension"`
							Extra2             string `json:"extra2"`
							Filename           string `json:"filename"`
							FilesizeMbs        int    `json:"filesize_mbs"`
							HvID               string `json:"hv_id"`
							IsDownloadable     string `json:"is_downloadable"`
							IsGuestAllowed     string `json:"is_guest_allowed"`
							IsMemberAllowed    string `json:"is_member_allowed"`
							IsPremiumAllowed   string `json:"is_premium_allowed"`
							Kind               string `json:"kind"`
							MimeType           string `json:"mime_type"`
							ServerSequence     string `json:"server_sequence"`
							Slug               string `json:"slug"`
							URL                string `json:"url"`
							VideoStreamGroupID string `json:"video_stream_group_id"`
						} `json:"streams"`
					} `json:"servers"`
				} `json:"videos_manifest"`
			} `json:"video"`
		} `json:"data"`
	} `json:"state"`
}

type manifest struct {
	VideosManifest struct {
		Servers []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			NaRating    int    `json:"na_rating"`
			EuRating    int    `json:"eu_rating"`
			AsiaRating  int    `json:"asia_rating"`
			Sequence    int    `json:"sequence"`
			IsPermanent bool   `json:"is_permanent"`
			Streams     []struct {
				ID                 int    `json:"id"`
				ServerID           int    `json:"server_id"`
				Slug               string `json:"slug"`
				Kind               string `json:"kind"`
				Extension          string `json:"extension"`
				MimeType           string `json:"mime_type"`
				Width              int    `json:"width"`
				Height             string `json:"height"`
				DurationInMs       int    `json:"duration_in_ms"`
				FilesizeMbs        int    `json:"filesize_mbs"`
				Filename           string `json:"filename"`
				URL                string `json:"url"`
				IsGuestAllowed     bool   `json:"is_guest_allowed"`
				IsMemberAllowed    bool   `json:"is_member_allowed"`
				IsPremiumAllowed   bool   `json:"is_premium_allowed"`
				IsDownloadable     bool   `json:"is_downloadable"`
				Compatibility      string `json:"compatibility"`
				HvID               int    `json:"hv_id"`
				ServerSequence     int    `json:"server_sequence"`
				VideoStreamGroupID string `json:"video_stream_group_id"`
				Extra2             any    `json:"extra2"`
			} `json:"streams"`
		} `json:"servers"`
	} `json:"videos_manifest"`
	Errors []string `json:"errors"`
}

const site = "https://hanime.tv/"
const video_manifest_api = "https://cached.freeanimehentai.net/api/v8/guest/videos/%s/manifest"

// returns 1=function parameters (for substitution) 2=JS struct (JSON info) 3=values passed to function
var reNuxtState = regexp.MustCompile(`function\((?<PARAMS>[^\)]+)\){return (?<JSON>{"?layout:"?[\s\S].*?)}\((?<VALUES>[^)]+)`)

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

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	nuxtMatch := reNuxtState.FindStringSubmatch(htmlString)
	for idx, g := range reNuxtState.SubexpNames() {
		if idx == 0 {
			continue
		}
		if nuxtMatch[idx] == "" {
			return nil, fmt.Errorf("regex match group: '%s' returned empty", g)
		}
	}

	jsonString := nuxtMatch[2]

	reKey := regexp.MustCompile(`([{\[,])([a-zA-Z_][\w_]*):`)
	jsonString = reKey.ReplaceAllString(jsonString, `$1"$2":`)

	reValue := regexp.MustCompile(`:([a-zA-Z_$\.][\w_]*)`)
	jsonString = reValue.ReplaceAllString(jsonString, `:"$1"`)

	funcParams := strings.Split(nuxtMatch[1], ",")
	substitutes := strings.Split(nuxtMatch[3], ",")

	replacements, err := zipStringSlices(funcParams, substitutes)
	if err != nil {
		return nil, err
	}
	for i := range replacements {
		replacements[i] = fmt.Sprintf(`"%s"`, strings.Trim(replacements[i], `"`))
	}

	replacer := strings.NewReplacer(replacements...)
	jsonString = replacer.Replace(jsonString)

	vData := videoData{}
	err = json.Unmarshal([]byte(jsonString), &vData)
	if err != nil {
		log.Println(jsonString)
		return nil, err
	}

	manifest_URL := fmt.Sprintf(video_manifest_api, vData.State.Data.Video.HentaiVideo.ID)
	manifest_res, err := request.GetAsBytesWithHeaders(manifest_URL, map[string]string{
		"Accept":  "application/json",
		"Referer": site,
		"Origin":  strings.TrimSuffix(site, "/"),
		"x-time":  fmt.Sprint(time.Now().Unix()),
	})
	if err != nil {
		return nil, err
	}

	manifest := manifest{}
	err = json.Unmarshal(manifest_res, &manifest)
	if err != nil {
		return nil, err
	}
	if len(manifest.Errors) > 0 {
		fmt.Println(manifest_URL)
		fmt.Println(string(manifest_res))
		return nil, errors.WithStack(errors.New(strings.Join(manifest.Errors, ";")))
	}

	// remove first entry if it's the 1080p stream since it only works if you are logged in
	if manifest.VideosManifest.Servers[0].Streams[0].Height == "1080" {
		manifest.VideosManifest.Servers[0].Streams = slices.Delete(manifest.VideosManifest.Servers[0].Streams, 0, 1)
	}

	streams := map[string]*static.Stream{}
	for idx, streamData := range manifest.VideosManifest.Servers[0].Streams {
		mediaStr, err := request.Get(streamData.URL)
		if err != nil {
			return nil, err
		}

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

// zipStringSlices merge two slices with alternating elements. a is first. Note: replace with iter when stable
func zipStringSlices(a, b []string) ([]string, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("zipping is only supported for two slices of the same length")
	}
	zipped := make([]string, len(a)*2)

	for i := range len(a) {
		zipped = append(zipped, a[i], b[i])
	}

	return zipped, nil
}
