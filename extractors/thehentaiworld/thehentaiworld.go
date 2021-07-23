package thehentaiworld

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type tag struct {
	ID uint `json:"id"`
}

type post struct {
	ID   uint   `json:"id"`
	Link string `json:"Link"`
}

type title struct {
	Rendered string `json:"rendered"`
}

type audio struct {
	DataFormat    string `json:"dataformat"`
	Codec         string `json:"codec"`
	SampleRate    uint   `json:"sample_rate"`
	Channels      uint   `json:"channels"`
	BitsPerSample uint   `json:"bits_per_sample"`
	Lossless      bool   `json:"lossless"`
	ChannelMode   string `json:"channelmode"`
}

func (a *audio) string() string {
	out := ""
	if a.DataFormat != "" {
		out += " format: " + a.DataFormat + ","
	}
	if a.Codec != "" {
		out += " codec: " + a.Codec + ","
	}
	if a.SampleRate != 0 {
		out += " sample rate: " + fmt.Sprint(a.SampleRate) + ","
	}
	if a.Channels != 0 {
		out += " channels: " + fmt.Sprint(a.Channels) + ","
	}
	if a.BitsPerSample != 0 {
		out += " bits per sample: " + fmt.Sprint(a.BitsPerSample) + ","
	}
	out += " lossless: " + fmt.Sprint(a.Lossless) + ","
	if a.ChannelMode != "" {
		out += " channel mode: " + fmt.Sprint(a.ChannelMode) + ","
	}
	return strings.Trim(out, " ,")
}

type size struct {
	Width     uint   `json:"width"`
	Height    uint   `json:"height"`
	MimeType  string `json:"mime_type"`
	SourceURL string `json:"source_url"`
}

type mediaDetails struct {
	FileSize uint            `json:"filesize"`
	Width    uint            `json:"width"`
	Height   uint            `json:"height"`
	Sizes    map[string]size `json:"sizes"`
	Audio    audio           `json:"audio"`
}

type media struct {
	Slug         string       `json:"slug"`
	Title        title        `json:"title"`
	MediaType    string       `json:"media_type"`
	MimeType     string       `json:"mime_type"`
	MediaDetails mediaDetails `json:"media_details"`
	ID           uint         `json:"id"`
	SourceURL    string       `json:"source_url"`
}

const site = "https://thehentaiworld.com/"
const postPerPage = "24"
const postsAPI = "https://thehentaiworld.com/wp-json/wp/v2/posts?"
const tagsAPI = "https://thehentaiworld.com/wp-json/wp/v2/tags?slug="
const mediaAPI = "https://thehentaiworld.com/wp-json/wp/v2/media?parent="

// https://thehentaiworld.com/wp-json/wp/v2/categories
var rePost *regexp.Regexp = regexp.MustCompile(`https://thehentaiworld.com/(?:3d-cgi-hentai-images|gif-animated-hentai-images|hentai-cosplay-images|hentai-doujinshi|flash-hentai|hentai-images|videos)/([^/]+)`)
var rePage *regexp.Regexp = regexp.MustCompile(`https://thehentaiworld.com/(?:tag/([^/]+)/)?(?:page/(\d+)/)?(\?s=[^&\n]+)?`)

// downloading large amount of content with -a might take a while
// the api call is quite slow

type extractor struct{}

// New returns a thehentaiworld extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	postIDs := parseURL(URL)
	if len(postIDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, pID := range postIDs {
		d, err := extractData(pID)
		if err != nil {
			return nil, utils.Wrap(err, pID)
		}
		data = append(data, d...)
	}

	return data, nil
}

func parseURL(URL string) []string {
	slug := rePost.FindStringSubmatch(URL)
	if len(slug) == 2 {
		htmlString, err := request.Get(postsAPI + "slug=" + slug[1])
		if err != nil {
			return nil
		}

		posts := []post{}
		err = json.Unmarshal([]byte(htmlString), &posts)
		if err != nil {
			return nil
		}

		if len(posts) < 1 {
			return nil
		}

		return []string{fmt.Sprint(posts[0].ID)}
	}

	matchedPage := rePage.FindStringSubmatch(URL)
	if len(matchedPage) < 2 {
		matchedPage = []string{"0", "1"}
	}
	currentPage, _ := strconv.ParseInt(matchedPage[2], 10, 0)

	if currentPage < 1 {
		currentPage = 1
	}

	tmpURL := postsAPI + "per_page=" + postPerPage + "&page=%d"

	if len(matchedPage) == 4 {
		search := strings.Split(matchedPage[3], "=")
		if len(search) == 2 {
			tmpURL += "&search=" + search[1]
		}
		if matchedPage[1] != "" {
			tmpURL += "&tags=" + get_tagID_from_slug(matchedPage[1])
		}
	}

	out := []string{}
	count := 0
	for i := int(currentPage); ; {
		htmlString, err := request.Get(fmt.Sprintf(tmpURL, i))
		if err != nil {
			return nil
		}
		if config.Amount > 0 {
			fmt.Println(count)
		}

		posts := []post{}
		json.Unmarshal([]byte(htmlString), &posts)

		for _, v := range posts {
			out = append(out, fmt.Sprint(v.ID))
		}
		count += len(posts)
		i += 1
		if config.Amount == 0 || count >= config.Amount || len(posts) == 0 {
			break
		}
	}

	if config.Amount > 0 && len(out) > config.Amount {
		out = out[:config.Amount]
	}

	return out
}

func extractData(pID string) ([]*static.Data, error) {
	//set per_page value to max (100) I have seen posts that contain 20+ images
	mediaJSON, err := request.Get(mediaAPI + pID + "&per_page=100")
	if err != nil {
		return nil, err
	}

	mS := []media{}
	err = json.Unmarshal([]byte(mediaJSON), &mS)
	if err != nil {
		return nil, err
	}

	// for vids you get thumbnail and video - overwrite with the video instead of the thumb
	if len(mS) > 1 {
		for i, v := range mS {
			if strings.Contains(v.Slug, "-thumbnail") {
				mS = remove(mS, i)
				break
			}
		}
	}

	data := []*static.Data{}
	for _, m := range mS {
		if len(m.MediaDetails.Sizes) == 0 {
			m.MediaDetails.Sizes = map[string]size{
				"full": {
					Width:     m.MediaDetails.Width,
					Height:    m.MediaDetails.Height,
					MimeType:  m.MimeType,
					SourceURL: m.SourceURL,
				},
			}
		}

		sizes := genSortedStreams(m.MediaDetails.Sizes)

		streams := map[string]*static.Stream{}
		for i, s := range sizes {
			streams[fmt.Sprint(i)] = &static.Stream{
				URLs: []*static.URL{
					{
						URL: s.SourceURL,
						Ext: utils.GetLastItemString(strings.Split(s.MimeType, "/")),
					},
				},
				Quality: fmt.Sprintf("%dp; %d x %d", s.Height, s.Width, s.Height),
				Size:    int64(m.MediaDetails.FileSize),
			}
			if m.MediaDetails.Audio.Codec != "" {
				streams[fmt.Sprint(i)].Info = m.MediaDetails.Audio.string()
			}
		}
		// need to append m.ID because neither title or Post + title is a unique file name
		data = append(data, &static.Data{
			Site:    site,
			Title:   fmt.Sprintf("%d – %s", m.ID, html.UnescapeString(m.Title.Rendered)),
			Type:    static.DataType(strings.Split(m.MimeType, "/")[0]),
			Streams: streams,
			Url:     "https://thehentaiworld.com/?p=" + fmt.Sprint(m.ID),
		})
	}

	return data, nil
}

func genSortedStreams(sizes map[string]size) []size {
	sortedSizes := make([]size, 0, len(sizes))
	for _, v := range sizes {
		sortedSizes = append(sortedSizes, v)
	}
	if len(sortedSizes) < 1 {
		return nil
	}
	sort.Slice(sortedSizes, func(i, j int) bool {
		return sortedSizes[i].Height*sortedSizes[i].Width > sortedSizes[j].Height*sortedSizes[j].Width
	})

	return sortedSizes
}

func remove(s []media, i int) []media {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func get_tagID_from_slug(slug string) string {
	jsonStr, err := request.Get(tagsAPI + slug)
	if err != nil {
		return ""
	}

	tags := []tag{}
	err = json.Unmarshal([]byte(jsonStr), &tags)
	if err != nil {
		return ""
	}

	return fmt.Sprint(tags[0].ID)
}
