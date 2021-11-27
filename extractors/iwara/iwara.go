package iwara

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type stream struct {
	Resolution string
	URI        string
	Mime       string
}

const site = "https://ecchi.iwara.tv/"
const videoAPI = "https://ecchi.iwara.tv/api/video/"

var reImgSource *regexp.Regexp = regexp.MustCompile(`([^"]+large/public/photos/[^"]+)"(?: width="([^"]*)[^=]+="([^"]*))`)
var reExt *regexp.Regexp = regexp.MustCompile(`(\w+)\?itok=[a-zA-Z\d]+$`)
var reTitle *regexp.Regexp = regexp.MustCompile(`<title>([^|]+)`)
var reVideoID *regexp.Regexp = regexp.MustCompile(`https://ecchi.iwara.tv/videos/(.+)`)

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
	if ok, _ := regexp.MatchString(site+`(?:videos|images)/`, URL); ok {
		return []string{URL}
	}

	tmpURL := regexp.MustCompile(`page=\d+`).ReplaceAllString(URL, "page=%d")
	if !strings.Contains(tmpURL, "page=%d") {
		tmpURL = URL + "&page=%d"
	}

	out := []string{}
	count := 0
	for i := 0; ; {
		htmlString, err := request.Get(fmt.Sprintf(tmpURL, i))
		if err != nil {
			return nil
		}
		if config.Amount > 0 {
			fmt.Println(count)
		}

		re := regexp.MustCompile(`/(?:videos|images)/[a-zA-Z0-9%=?-]+"`)
		matchedURLs := re.FindAllString(htmlString, -1)

		URLs := []string{}
		for _, matchedURL := range utils.RemoveAdjDuplicates(matchedURLs) {
			URLs = append(URLs, site+strings.Trim(matchedURL, `/"`))
		}
		count += len(URLs)
		i += 1
		out = append(out, URLs...)
		if config.Amount == 0 || count >= config.Amount || len(URLs) == 0 {
			break
		}
	}

	if config.Amount > 0 && len(out) > config.Amount {
		out = out[:config.Amount]
	}

	return out
}

func extractData(URL string) ([]*static.Data, error) {
	resString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetLastItemString(reTitle.FindStringSubmatch(resString))
	title = title[:len(title)-1]

	matchedImages := reImgSource.FindAllStringSubmatch(resString, -1)
	if len(matchedImages) > 0 {
		data := []*static.Data{}
		for i, img := range matchedImages {
			img[1] = "https:" + img[1]

			quality := ""
			if len(img) > 2 {
				quality = fmt.Sprintf("%s x %s", img[2], img[3])
			}

			size, _ := request.Size(img[1], site)

			data = append(data, &static.Data{
				Site:  site,
				Type:  "image",
				Title: fmt.Sprintf("%s_%d", title, i+1),
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: img[1],
								Ext: reExt.FindStringSubmatch(img[1])[1],
							},
						},
						Quality: quality,
						Size:    size,
					},
				},
				URL: URL,
			})
		}
		return data, nil
	}

	videoID := utils.GetLastItemString(reVideoID.FindStringSubmatch(URL))
	if videoID == "" {
		return nil, static.ErrURLParseFailed
	}

	jsonData, err := request.GetAsBytes(videoAPI + videoID)
	if err != nil {
		return nil, err
	}

	vStreams := []stream{}
	err = json.Unmarshal(jsonData, &vStreams)
	if err != nil {
		return nil, err
	}

	streams := map[string]*static.Stream{}
	for i, stream := range vStreams {
		stream.URI = "https:" + stream.URI

		size, _ := request.Size(stream.URI, site)

		streams[fmt.Sprint(i)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: stream.URI,
					Ext: utils.GetLastItemString(strings.Split(stream.Mime, "/")),
				},
			},
			Quality: stream.Resolution,
			Size:    size,
		}
	}

	return []*static.Data{
		{
			Site:    site,
			Title:   title,
			Type:    static.DataTypeVideo,
			Streams: streams,
			URL:     URL,
		},
	}, nil
}
