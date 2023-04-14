package nhentai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type gallery struct {
	ID      json.RawMessage `json:"id"`
	MediaID string          `json:"media_id"`
	Title   struct {
		English  string `json:"english"`
		Japanese string `json:"japanese"`
		Pretty   string `json:"pretty"`
	} `json:"title"`
	Images struct {
		Pages []struct {
			T string `json:"t"` //p=png j=jpg
			W int    `json:"w"`
			H int    `json:"h"`
		} `json:"pages"`
		Cover struct {
			T string `json:"t"`
			W int    `json:"w"`
			H int    `json:"h"`
		} `json:"cover"`
		Thumbnail struct {
			T string `json:"t"`
			W int    `json:"w"`
			H int    `json:"h"`
		} `json:"thumbnail"`
	} `json:"images"`
	Scanlator  string `json:"scanlator"`
	UploadDate int    `json:"upload_date"`
	Tags       []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Count int    `json:"count"`
	} `json:"tags"`
	NumPages     int `json:"num_pages"`
	NumFavorites int `json:"num_favorites"`
}

const site = "https://nhentai.net"
const cdn = "https://i.nhentai.net/galleries/"
const api = "https://nhentai.net/api/gallery/"

type extractor struct{}

// New returns a nhentai extractor
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	IDs, page := parseURL(URL)
	if len(IDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, id := range IDs {
		d, err := extractData(id, page)
		if err != nil {
			return nil, utils.Wrap(err, id)
		}
		data = append(data, d)
	}

	return data, nil
}

// parseURL data
func parseURL(URL string) ([]string, string) {
	re := regexp.MustCompile(`//[^/]*/([^/]*)`)
	matchedURLPart := re.FindStringSubmatch(URL)
	if len(matchedURLPart) < 1 {
		return nil, ""
	}

	switch matchedURLPart[1] {
	case "g":
		// if there are two "int" values it means the exact page was supplied
		var page string

		re = regexp.MustCompile(`[\d]+`)
		urlNumbers := re.FindAllString(URL, -1)
		if len(urlNumbers) <= 0 {
			return nil, ""
		}
		if len(urlNumbers) > 1 {
			page = urlNumbers[1]
		}

		return []string{urlNumbers[0]}, page
	case "search", "tag", "artist", "characters", "parodies", "groups":
		htmlString, err := request.Get(URL)
		if err != nil {
			return nil, ""
		}

		re = regexp.MustCompile(`/g/(\d*)/`)
		matchedIDs := re.FindAllStringSubmatch(htmlString, -1)

		ids := []string{}
		for _, id := range matchedIDs {
			ids = append(ids, id[1])
		}
		return ids, ""
	}
	return nil, ""
}

func extractData(id string, page string) (*static.Data, error) {
	URL := api + id

	apiRes, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	if utils.GetH1(&apiRes, 0) == "429 Too Many Requests" {
		time.Sleep(250 * time.Millisecond)
		apiRes, _ = request.Get(URL)
	}

	gData := &gallery{}
	err = json.Unmarshal([]byte(apiRes), &gData)
	if err != nil {
		return nil, err
	}

	pages := utils.NeedDownloadList(gData.NumPages)
	if page != "" {
		pageNo, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		pages = []int{pageNo}
	}

	URLs := []*static.URL{}
	for _, p := range pages {
		ext := gData.Images.Pages[p-1].T
		switch ext {
		case "j":
			ext = "jpg"
		case "p":
			ext = "png"
		default:
			ext = "gif"
		}
		URLs = append(URLs, &static.URL{
			URL: fmt.Sprintf("%s%s/%d.%s", cdn, gData.MediaID, p, ext),
			Ext: ext,
		})
	}

	return &static.Data{
		Site:  site,
		Title: gData.Title.Pretty,
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: URL,
	}, nil
}
