/**
	Unfortunatly tsumino's "auth system" is quite sharp making mass extraction impossilble
	if you still want to mass extract doujins use nhentai or one of the other extractors instead
**/
package tsumino

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

type tag struct {
	Type    string `json:"Type"`
	Text    string `json:"Text"`
	Exclude string `json:"Exclude"`
}

type searchParams struct {
	PageNumber    string                 `json:"PageNumber"`
	Text          string                 `json:"Text"`
	List          string                 `json:"List"`
	Length        string                 `json:"Length"`
	MinimumRating string                 `json:"MinimumRating"`
	Sort          string                 `json:"Sort"`
	Include       map[string]interface{} `json:"Include"`
	Tags          []tag                  `json:"Tags"`
}

type entry struct {
	ID             uint32  `json:"id"`
	Title          string  `json:"title"`
	Rating         float64 `json:"rating"`
	Duration       uint64  `json:"duration"`
	CancelPosition float64 `json:"collectionPosition"`
	EntryType      string  `json:"entryType"`
}

type data struct {
	Entry       entry   `json:"entry"`
	Impression  string  `json:"impression"`
	HistoryPage float64 `json:"historyPage"`
}

type searchResult struct {
	PageNumber uint32 `json:"pageNumber"`
	PageCount  uint32 `json:"pageCount"`
	Data       []data `json:"data"`
}

const site = "https://www.tsumino.com/"
const baseSearchURL = "https://www.tsumino.com/Search/Operate/?type="
const galleryMainPage = "https://www.tsumino.com/Read/Index/"
const cdn = "https://content.tsumino.com/"

func ParseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"entry/") {
		return []string{URL}
	}

	paramsMap, err := decodeURL(URL)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	params := url.Values{}
	for k, v := range paramsMap {
		params.Add(k, v)
	}

	mediaType := strings.Split(strings.TrimPrefix(URL, site), "#")[0] //books or videos
	req, err := http.NewRequest(http.MethodPost, baseSearchURL+mediaType, strings.NewReader(params.Encode()))
	if err != nil {
		return []string{}
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	searchRes := searchResult{}
	err = json.Unmarshal(buffer, &searchRes)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	URLs := []string{}
	for _, v := range searchRes.Data {
		URLs = append(URLs, "https://www.tsumino.com/entry/"+fmt.Sprint(v.Entry.ID))
	}
	//https://www.tsumino.com/Search/Operate/?type=Book
	//https://www.tsumino.com/Read/Index/55285?page=1
	return URLs
}

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, fmt.Errorf("[Tsumino] No scrapable URL found for %s", URL)
	}

	data := []static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}

	title := html.UnescapeString(utils.GetMeta(htmlString, "og:title"))
	id := utils.GetLastItemString(strings.Split(URL, "/"))

	if !strings.Contains(htmlString, "<a href=\"/Read/Index/") {
		headers, err := request.Headers(fmt.Sprintf("%svideos/%s/video.mp4", cdn, id), site)
		if err != nil {
			return static.Data{}, err
		}

		s := headers.Get("Content-Range")
		if s == "" {
			return static.Data{}, errors.New("content-range is not present")
		}
		size, err := strconv.ParseInt(strings.Split(s, "/")[1], 10, 64)
		if err != nil {
			return static.Data{}, err
		}

		return static.Data{
			Site:  site,
			Title: title,
			Type:  headers.Get("content-type"),
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						{
							URL: fmt.Sprintf("%svideos/%s/video.mp4", cdn, id),
							Ext: "mp4",
						},
					},
					Quality: "best",
					Size:    size,
				},
			},
			Url: URL,
		}, nil
	}

	htmlString, err = request.Get(galleryMainPage + id + "?page=1")
	if err != nil {
		return static.Data{}, err
	}

	re := regexp.MustCompile(`"[^"]*\?key[^"]*"`)
	templateURL, _ := strconv.Unquote(re.FindString(htmlString))
	templateURL = html.UnescapeString(templateURL)

	re = regexp.MustCompile(`of (\d*)</h1>`)
	matchedNumOfPages := re.FindStringSubmatch(htmlString)
	if len(matchedNumOfPages) < 2 {
		return static.Data{}, fmt.Errorf("[Tsumino] No numOfPages found for %s", URL)
	}

	numOfPages, _ := strconv.Atoi(matchedNumOfPages[1])

	pages := utils.NeedDownloadList(numOfPages)
	URLs := []static.URL{}
	for _, v := range pages {
		URLs = append(URLs, static.URL{
			URL: strings.Replace(templateURL, "[PAGE]", fmt.Sprint(v), -1),
			Ext: "jpg",
		})
	}

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "image",
		Streams: map[string]static.Stream{
			"0": {
				URLs:    URLs,
				Quality: "best",
			},
		},
	}, nil
}

func decodeURL(URL string) (map[string]string, error) {
	// no parsable URL parameters
	if !strings.ContainsAny(URL, "?~(") {
		return nil, fmt.Errorf("[Tsumino] no URL params found in URL %s", URL)
	}

	if !strings.ContainsAny(URL, "~(") {
		//default search
		return nil, nil
	}

	re := regexp.MustCompile(`#(.+)#`)
	matchedParamsString := re.FindStringSubmatch(URL)
	if len(matchedParamsString) < 2 {
		return nil, fmt.Errorf("[Tsumino] no URL params found in URL %s", URL)
	}

	paramsString := matchedParamsString[1]
	// ----- begin of json(fication) -----
	// 1. Transform () to {}
	re = regexp.MustCompile(`([a-zA-Z])~\(`)
	paramsString = re.ReplaceAllString(paramsString, `$1:{`)

	// 1.a transform the rest of the ~( that don't preceed with a character
	paramsString = strings.ReplaceAll(paramsString, "~(", "{")

	//1.b transform )~ to },
	paramsString = strings.ReplaceAll(paramsString, ")~", "},")

	//1.c transform ) to }
	paramsString = strings.ReplaceAll(paramsString, ")", "}")

	// 2. Make key value pairs
	// 2.a \w~\w to \w:\w
	paramsString = strings.ReplaceAll(paramsString, "'", "")
	re = regexp.MustCompile(`(\w+)~(['\w]+)`)
	paramsString = re.ReplaceAllString(paramsString, `$1:$2,`)

	// 2.b transform ,~ to ,
	paramsString = strings.ReplaceAll(paramsString, ",~", ",")

	// 2.c transform ,} to }
	paramsString = strings.ReplaceAll(paramsString, ",}", "}")

	// 2.d transform {{ and }} to [{ }]
	paramsString = strings.ReplaceAll(paramsString, "{{", "[{")
	paramsString = strings.ReplaceAll(paramsString, "}}", "}]")

	// 3. Add quotes
	// 3.a \w+ to "$&"
	re = regexp.MustCompile(`(\w+)`)
	paramsString = re.ReplaceAllString(paramsString, `"$1"`)

	// 3.b transform ~ to "" (empty value)
	paramsString = strings.ReplaceAll(paramsString, "~", "")

	// to json
	searchP := searchParams{}
	//var inInterface map[string]interface{}
	err := json.Unmarshal([]byte(paramsString), &searchP)
	if err != nil {
		return nil, err
	}
	// ----- end of json(fication) -----

	out := map[string]string{}
	out["PageNumber"] = searchP.PageNumber
	if searchP.Text != "" {
		out["Text"] = searchP.Text
	}
	if searchP.List != "" {
		out["List"] = searchP.List
	}
	if searchP.Length != "" {
		out["Length"] = searchP.Length
	}
	if searchP.MinimumRating != "" {
		out["MinimumRating"] = searchP.MinimumRating
	}
	out["Sort"] = searchP.Sort
	for i, v := range searchP.Tags {
		out[fmt.Sprintf("Tags[%d][%s]", i, "Type")] = v.Type
		out[fmt.Sprintf("Tags[%d][%s]", i, "Text")] = v.Text
		out[fmt.Sprintf("Tags[%d][%s]", i, "Exclude")] = v.Exclude
	}

	return out, nil
}
