package ninehentai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/jsurl"
)

//9hentai

type resultSearch struct {
	TotalCount int       `json:"total_count"`
	Results    []gallery `json:"results"`
}

type resultBook struct {
	Results gallery `json:"results"`
}

type resultTag struct {
	Results tag `json:"results"`
}

type tag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type uint   `json:"type"`
}

type items struct {
	Excluded []tag `json:"excluded"`
	Included []tag `json:"included"`
}

type searchTag struct {
	Items items         `json:"items"`
	Tags  []interface{} `json:"tags"`
	Text  string        `json:"text"`
	Type  uint          `json:"type"`
}

type pages struct {
	Range []uint `json:"range"`
}

type search struct {
	Page  uint      `json:"page"`
	Pages pages     `json:"pages"`
	Sort  uint      `json:"sort"`
	Tag   searchTag `json:"tag"`
	Text  string    `json:"text"`
}

type searchReq struct {
	Search search `json:"search"`
}

type gallery struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	AltTitle    string `json:"alt_title"`
	TotalPage   uint   `json:"total_page"`
	ImageServer string `json:"image_server"`
}

const site = "https://9hentai.to/"
const apiURLGetTagByID = "https://9hentai.to/api/getTagByID"
const apiURLGetBookByID = "https://9hentai.to/api/getBookByID"
const apiURLGetBook = "https://9hentai.to/api/getBook"

func ParseURL(URL string) ([]gallery, error) {
	re := regexp.MustCompile(`/([gt])/(\d+)/`) //1=indicator g=gallery t=tag etc=searchQuery?
	matchedURLParams := re.FindStringSubmatch(URL)
	if len(matchedURLParams) < 2 {
		return nil, fmt.Errorf("[9Hentai] URL parameters cannot be parsed: %s", URL)
	}

	switch matchedURLParams[1] {
	case "g":
		g, err := getBookByID(matchedURLParams[2])
		if err != nil {
			return nil, err
		}
		return []gallery{g}, nil

	case "t":
		// one tag is supplied
		if !strings.Contains(URL, "#") {
			t, err := getTagByID(matchedURLParams[2])
			if err != nil {
				return nil, err
			}

			s := searchReq{
				search{
					Tag: searchTag{
						Items: items{
							Included: []tag{t},
						},
						Type: 1,
					},
					Pages: pages{
						Range: []uint{0, 2000},
					},
				},
			}

			rGalleries, err := getBook(s)
			if err != nil {
				return nil, err
			}

			return rGalleries, nil
		}

		urlParams := strings.TrimSuffix(strings.Split(URL, "#")[1], "#")

		searchFromURLParams := search{}
		err := jsurl.Parse(urlParams, &searchFromURLParams)
		if err != nil {
			return nil, err
		}

		rGalleries, err := getBook(searchReq{Search: searchFromURLParams})
		if err != nil {
			return nil, err
		}

		return rGalleries, nil

	default:
		return nil, fmt.Errorf("[9Hentai] URL indicator cannot be parsed: %s. Expected t or g got %s", URL, matchedURLParams[1])
	}

}

func getTagByID(ID string) (tag, error) {
	res, err := request.Request(http.MethodPost, apiURLGetTagByID, map[string]string{
		"content-type": "application/json",
	}, strings.NewReader(fmt.Sprintf("{\"id\": %s}", ID)))
	if err != nil {
		return tag{}, err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return tag{}, err
	}

	rTag := resultTag{}
	err = json.Unmarshal(buffer, &rTag)
	if err != nil {
		return tag{}, err
	}

	return rTag.Results, nil
}

func getBookByID(ID string) (gallery, error) {
	res, err := request.Request(http.MethodPost, apiURLGetBookByID, map[string]string{
		"content-type": "application/json",
	}, strings.NewReader(fmt.Sprintf("{\"id\": %s}", ID)))
	if err != nil {
		return gallery{}, err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return gallery{}, err
	}

	rBook := resultBook{}
	err = json.Unmarshal(buffer, &rBook)
	if err != nil {
		return gallery{}, err
	}

	return rBook.Results, nil
}

func getBook(s searchReq) ([]gallery, error) {
	buffer, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	res, err := request.Request(http.MethodPost, apiURLGetBook, map[string]string{
		"content-type": "application/json",
	}, bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buffer, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rSearch := resultSearch{}
	err = json.Unmarshal(buffer, &rSearch)
	if err != nil {
		return nil, err
	}

	return rSearch.Results, nil
}

func Extract(URL string) ([]static.Data, error) {
	galleries, err := ParseURL(URL)
	if err != nil {
		return nil, fmt.Errorf("[9Hentai] No scrapable gallery id found for: %s", URL)
	}

	data := []static.Data{}
	for _, g := range galleries {
		d, err := extractData(g)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(g gallery) (static.Data, error) {
	URLs := []static.URL{}
	for i := 1; i < int(g.TotalPage); i++ {
		URLs = append(URLs, static.URL{
			URL: fmt.Sprintf("%s%d/%d.jpg", g.ImageServer, g.ID, i),
			Ext: "jpg",
		})
	}

	return static.Data{
		Site:  site,
		Title: g.Title,
		Type:  "image",
		Streams: map[string]static.Stream{
			"0": {
				URLs:    URLs,
				Quality: "best",
				Info:    fmt.Sprint(g.TotalPage),
			},
		},
		Url: fmt.Sprintf("%sg/%d", site, g.ID),
	}, nil
}
