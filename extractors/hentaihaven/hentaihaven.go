package hentaihaven

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type source struct {
	Src   string `json:"src"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type captions struct {
	Src   string `json:"src"`
	Label string `json:"label"`
}

type data struct {
	Image    string   `json:"image"`
	Mosaic   string   `json:"mosaic"`
	Captions captions `json:"captions"`
	Sources  []source `json:"sources"`
}

type pData struct {
	Status bool `json:"status"`
	Data   data `json:"data"`
}

const site = "https://hentaihaven.xxx/"
const api = "https://hentaihaven.xxx/wp-admin/admin-ajax.php"

type extractor struct{}

// New returns a hentaihaven.xxx extractor.
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
		data = append(data, &d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if ok, _ := regexp.MatchString(`/episode-\d*/?$`, URL); ok {
		return []string{URL}
	}

	if !strings.Contains(URL, "https://hentaihaven.xxx/watch/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}
	slug := strings.Split(URL, "watch/")[1]
	re := regexp.MustCompile(fmt.Sprintf("[^\"]*%sepisode-\\d*", slug))
	return re.FindAllString(htmlString, -1)
}

func extractData(URL string) (static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return static.Data{}, err
	}
	title := strings.Trim(utils.GetH1(&htmlString, -1), " \n\t")

	re := regexp.MustCompile(`[^"]*/player/[^"]*`)
	playerURL := re.FindString(htmlString) // 1=id  2=nonce
	if playerURL == "" {
		return static.Data{}, errors.New("can't locate player URL")
	}

	htmlString, err = request.Get(playerURL)
	if err != nil {
		return static.Data{}, err
	}

	// --- Begin of multipart creation
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	vals := [][]string{{"", "action", "zarat_get_data_player_ajax"}}
	re = regexp.MustCompile(`append\('([abc])', ?'([^']*)`) //1=a : some string b : some other string
	vals = append(vals, re.FindAllStringSubmatch(htmlString, -1)...)

	for _, v := range vals {
		mimeHeader := textproto.MIMEHeader{}
		mimeHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"", v[1]))
		part, _ := writer.CreatePart(mimeHeader)
		part.Write([]byte(v[2]))
	}
	writer.Close()
	// --- End of multipart creation
	// This will create some thing like this
	//------WebKitFormBoundaryDyxVGG0MJMgqpBFh
	//Content-Disposition: form-data; name="action"
	//
	//zarat_get_data_player_ajax
	//------WebKitFormBoundaryDyxVGG0MJMgqpBFh
	//Content-Disposition: form-data; name="a"
	//
	//NaRHayKOyzVTAkNnrg9SLSoYh2BTyYfgWfGO2jWz0NrecL/Vo55dZ8aXX9VztkUcSl8qKRd6GF/8SFfC47WyQEi+Z/Ii4n2FzPzmJwKlefvLxcLZBAJfopxo8M1XfEljw5E9fNOaL/5KMklhF+zwWOvI+lfu0A/hT2Sv5jFPn3k=
	//------WebKitFormBoundaryDyxVGG0MJMgqpBFh
	//Content-Disposition: form-data; name="b"
	//
	//RklZWG9ub0hiWnl5VUR2Y2tSYUpMdz09
	//------WebKitFormBoundaryDyxVGG0MJMgqpBFh--

	res, err := request.Request(http.MethodPost, api, map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}, body)
	if err != nil {
		return static.Data{}, err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return static.Data{}, err
	}

	sources := &pData{}
	//there are 3 weird bytes at the beginning that can't be interpreted so I removed them
	err = json.Unmarshal(respBody[3:], &sources)
	if err != nil {
		return static.Data{}, err
	}
	if !sources.Status {
		return static.Data{}, errors.New("the api request for the streams did not return successful for")
	}

	m3u8String, err := request.Get(sources.Data.Sources[0].Src)
	if err != nil {
		return static.Data{}, err
	}

	baseURL, err := url.Parse(sources.Data.Sources[0].Src)
	if err != nil {
		return static.Data{}, err
	}

	streams, err := utils.ParseM3UMaster(&m3u8String)
	if err != nil {
		return static.Data{}, err
	}

	idx := 0
	out := map[string]*static.Stream{}
	for _, variant := range streams {
		mediaURL, err := baseURL.Parse(variant.URLs[0].URL)
		if err != nil {
			return static.Data{}, err
		}

		mediaStr, err := request.Get(mediaURL.String())
		if err != nil {
			return static.Data{}, err
		}

		URLs, key, err := request.GetM3UMeta(&mediaStr, mediaURL.String(), "ts")
		if err != nil {
			return static.Data{}, err
		}

		out[strconv.Itoa(len(streams)-idx-1)] = &static.Stream{
			URLs:    URLs,
			Quality: variant.Quality,
			Size:    variant.Size,
			Ext:     "ts",
			Key:     key,
		}
		idx += 1
	}

	return static.Data{
		Site:    site,
		Title:   title,
		Type:    "video",
		Streams: out,
		Url:     URL,
	}, nil

}
