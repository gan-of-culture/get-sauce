package nhentai

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type tag struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

type page struct {
	T string `json:"t"` //j=jpg p=png
	W int    `json:"w"`
	H int    `json:"h"`
}

type images struct {
	Pages      []page `json:"pages"`
	Cover      page   `json:"cover"`
	Thumbnails page   `json:"thumbnail"`
}

type gallery struct {
	ID           json.RawMessage   `json:"id"` //can be string or int
	MediaID      string            `json:"media_id"`
	Title        map[string]string `json:"title"`
	Images       images            `json:"images"`
	Scanlator    string            `json:"scanlator"`
	UploadDate   int               `json:"upload_date"`
	Tags         []tag             `json:"tags"`
	NumPages     int               `json:"num_pages"`
	NumFavorites int               `json:"num_favorites"`
}

const site = "https://nhentai.net"
const cdn = "https://i.nhentai.net/galleries/"

type extractor struct{}

// New returns a nhentai extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	ids, page := parseURL(URL)
	if len(ids) < 1 {
		return nil, fmt.Errorf("This is not a vaild URL %s", URL)
	}
	data := []*static.Data{}
	for _, id := range ids {
		d, err := extractData(id, page)
		if err != nil {
			return nil, err
		}
		data = append(data, &d)
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

func extractData(id string, page string) (static.Data, error) {
	URL := fmt.Sprintf("https://nhentai.net/g/%s/", id)
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	if utils.GetH1(&htmlString, 0) == "429 Too Many Requests" {
		time.Sleep(250 * time.Millisecond)
		htmlString, _ = request.Get(URL)
	}

	re := regexp.MustCompile(`("{\\u0022[\s\S]*?}")`)
	matchedJsonString := re.FindStringSubmatch(htmlString)
	if len(matchedJsonString) < 2 {
		fmt.Println(htmlString)
		return static.Data{}, fmt.Errorf("[NHentai] invalid JSON %s", URL)
	}
	jsonString, _ := strconv.Unquote(matchedJsonString[1])

	gData := &gallery{}
	err = json.Unmarshal([]byte(jsonString), &gData)
	if err != nil {
		log.Println(jsonString)
		log.Println(URL)
		return static.Data{}, err
	}

	pages := utils.NeedDownloadList(int(gData.NumPages))
	if page != "" {
		pageNo, err := strconv.Atoi(page)
		if err != nil {
			return static.Data{}, err
		}
		pages = []int{pageNo}
	}

	URLs := []static.URL{}
	for _, p := range pages {
		ext := "jpg"
		if gData.Images.Pages[p-1].T == "p" {
			ext = "png"
		}
		URLs = append(URLs, static.URL{
			URL: fmt.Sprintf("%s%s/%d.%s", cdn, gData.MediaID, p, ext),
			Ext: ext,
		})
	}

	title, ok := gData.Title["pretty"]
	if !ok {
		return static.Data{}, fmt.Errorf("[NHentai] Cannot find title for %s", URL)
	}

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "image",
		Streams: map[string]static.Stream{
			"0": {
				URLs:    URLs,
				Quality: "best",
				Size:    0,
				Info:    fmt.Sprintf("Has %d pages", gData.NumPages),
			},
		},
		Url: URL,
	}, nil
}
