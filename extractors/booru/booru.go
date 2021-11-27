package booru

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
)

const site = "https://booru.io/"
const postURL = "https://booru.io/p/"
const apiDataURL = "https://booru.io/api/legacy/data/"
const apiEntityURL = "https://booru.io/api/legacy/entity/"
const apiQueryURL = "https://booru.io/api/legacy/query/entity?query="

// Entity JSON type
type Entity struct {
	Key         string `json:"key"`
	ContentType string `json:"contentType"`
	Attributes  struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"attributes"`
	Tags       map[string]int    `json:"tags"`
	Transforms map[string]string `json:"transforms"`
}

// EntitySlice JSON type
type EntitySlice struct {
	Data   []Entity `json:"data"`
	Cursor string   `json:"cursor"`
}

var reKey = regexp.MustCompile(`[0-9]+`)

type extractor struct{}

// New returns a booru.io extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract for booru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	query, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	data, err := extractData(query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// parseURL for danbooru pages
func parseURL(URL string) (string, error) {
	if strings.HasPrefix(URL, "https://booru.io/p/") {
		re := regexp.MustCompile(`https://booru\.io/p/(.+)`)
		matchedID := re.FindStringSubmatch(URL)
		if len(matchedID) > 2 {
			return "", static.ErrURLParseFailed
		}
		return fmt.Sprintf("%s%s", apiEntityURL, matchedID[1]), nil
	}

	tags := strings.Split(URL, "https://booru.io/q/")

	return fmt.Sprintf("%s%s", apiQueryURL, tags[1]), nil
}

func extractData(queryURL string) ([]*static.Data, error) {
	jsonData, err := request.GetAsBytes(queryURL)
	if err != nil {
		fmt.Println(queryURL)
		return nil, err
	}

	entitySlice := EntitySlice{}
	//single post
	if !strings.Contains(queryURL, "=") {
		entity := Entity{}
		err := json.Unmarshal(jsonData, &entity)
		if err != nil {
			fmt.Println(queryURL)
			return nil, err
		}
		entitySlice.Data = append(entitySlice.Data, entity)
	}

	if len(entitySlice.Data) == 0 {

		cursor := 0
		for {
			if config.Amount > 0 && config.Amount <= cursor {
				break
			}
			entitySliceTmp := EntitySlice{}
			err = json.Unmarshal(jsonData, &entitySliceTmp)
			if err != nil {
				fmt.Println(queryURL)
				fmt.Println("Cursor", cursor)
				fmt.Println(jsonData)
			}
			if len(entitySliceTmp.Data) == 0 && err == nil {
				break
			}
			entitySlice.Data = append(entitySlice.Data, entitySliceTmp.Data...)
			cursor += 50
			jsonData, err = request.GetAsBytes(fmt.Sprintf("%s&cursor=%d", queryURL, cursor))
			fmt.Printf("%s&cursor=%d", queryURL, cursor)
			if err != nil {
				return nil, err
			}
			time.Sleep(50 * time.Millisecond)
		}
	}

	data := []*static.Data{}
	for _, e := range entitySlice.Data {
		tType, tVal := getBestQualityImg(e.Transforms)
		ext := getFileExt(tType)
		size, _ := request.Size(fmt.Sprintf("%s%s", apiDataURL, tVal), site)

		data = append(data, &static.Data{
			Site:  site,
			Title: e.Key,
			Type:  "image",
			Streams: map[string]*static.Stream{
				"0": {
					Type: static.DataTypeImage,
					URLs: []*static.URL{
						{
							URL: fmt.Sprintf("%s%s", apiDataURL, tVal),
							Ext: ext,
						},
					},
					Quality: fmt.Sprintf("%d x %d", e.Attributes.Width, e.Attributes.Height),
					Size:    size,
				},
			},
			URL: fmt.Sprintf("%s%s", postURL, e.Key),
		})
	}

	return data, nil
}

func getBestQualityImg(transformations map[string]string) (string, string) {
	transformationType := ""
	transformationValue := ""
	currentBest := 0
	for key, val := range transformations {
		resString := reKey.FindString(key)
		resolution, _ := strconv.Atoi(resString)
		if resolution <= 0 {
			continue
		}
		if resolution < currentBest {
			continue
		}
		currentBest = resolution
		transformationType = key
		transformationValue = val
	}
	return transformationType, transformationValue
}

func getFileExt(tranformation string) string {
	transSplit := strings.Split(tranformation, "/")
	if len(transSplit) > 1 {
		if transSplit[1] == "jpeg" {
			return "jpg"
		}
		return transSplit[1]
	}
	return ""
}
