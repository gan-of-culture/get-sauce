package jwplayer

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

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type mediaData struct {
	Status bool `json:"status"`
	Data   struct {
		Image   string `json:"image"`
		Mosaic  string `json:"mosaic"`
		Sources []struct {
			Src   string `json:"src"`
			Type  string `json:"type"`
			Label string `json:"label"`
		} `json:"sources"`
	} `json:"data"`
}

const playerLocation = "/wp-content/plugins/player-logic/api.php"

var reJWPlayerURL = regexp.MustCompile(`[^"]+/wp-content/plugins/player-logic/player\.php[^"]+`)
var reMultiPartParams = regexp.MustCompile(`append\('([abc])', ?'([^']*)`) //1=a : some string b : some other string
var reURLBase = regexp.MustCompile(`https://[^/]+/`)

type extractor struct{}

// New returns a jwplayer extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract from URL
func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	URLBase := reURLBase.FindString(URL)
	if URLBase == "" {
		return nil, static.ErrURLParseFailed
	}

	apiURL := URLBase + playerLocation

	// --- Begin of multipart creation
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	vals := [][]string{{"", "action", "zarat_get_data_player_ajax"}}
	vals = append(vals, reMultiPartParams.FindAllStringSubmatch(htmlString, -1)...)

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

	res, err := request.Request(http.MethodPost, apiURL, map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sources := mediaData{}
	err = json.Unmarshal(respBody, &sources)
	if err != nil {
		return nil, err
	}

	if !sources.Status {
		return nil, errors.New("the jwplayer api request for the streams did not return successful for")
	}

	m3u8String, err := request.Get(sources.Data.Sources[0].Src)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(sources.Data.Sources[0].Src)
	if err != nil {
		return nil, err
	}

	streams, err := utils.ParseM3UMaster(&m3u8String)
	if err != nil {
		return nil, err
	}

	out := map[string]*static.Stream{}
	for idx, variant := range streams {
		mediaURL, err := baseURL.Parse(variant.URLs[0].URL)
		if err != nil {
			return nil, err
		}

		mediaStr, err := request.Get(mediaURL.String())
		if err != nil {
			return nil, err
		}

		URLs, key, err := request.GetM3UMeta(&mediaStr, mediaURL.String())
		if err != nil {
			return nil, err
		}

		out[fmt.Sprint(len(streams)-idx-1)] = &static.Stream{
			Type:    static.DataTypeVideo,
			URLs:    URLs,
			Quality: variant.Quality,
			Size:    variant.Size,
			Ext:     "mp4",
			Key:     key,
		}
	}

	return []*static.Data{
		{
			Site:    URLBase,
			Title:   "jwplayer video",
			Type:    static.DataTypeVideo,
			Streams: out,
		},
	}, nil
}

// FindJWPlayerURL in HTML page
func FindJWPlayerURL(htmlString *string) string {
	return reJWPlayerURL.FindString(*htmlString)
}
