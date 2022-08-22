package jwplayer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
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

var reJWPlayerURL = regexp.MustCompile(`[^"]+/wp-content/plugins/player-logic/player\.php[^"]+`)
var reMultiPartParams = regexp.MustCompile(`append\('([abc])', ?([^\)]*)`) //1=a : some string b : some other string
var reAPIURL = regexp.MustCompile("fetch\\(['`\"]([^'`\"]+api\\.php)")
var reVariable = regexp.MustCompile(`\$\{\w+\}`)

const findVarible = `var %s = '([^']+)`

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

	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	site := fmt.Sprintf("https://%s/", u.Host)

	matchedAPIURL := reAPIURL.FindStringSubmatch(htmlString)
	if len(matchedAPIURL) < 2 {
		return nil, static.ErrURLParseFailed
	}

	if variable := reVariable.FindString(htmlString); variable != "" {
		variableValue, err := findVariable(variable, &htmlString)
		if err != nil {
			return nil, err
		}

		matchedAPIURL[1] = strings.ReplaceAll(matchedAPIURL[1], variable, variableValue)
	}

	apiURL := matchedAPIURL[1]

	// --- Begin of multipart creation
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	vals := [][]string{{"", "action", "zarat_get_data_player_ajax"}}
	vals = append(vals, reMultiPartParams.FindAllStringSubmatch(htmlString, -1)...)

	for _, v := range vals {
		mimeHeader := textproto.MIMEHeader{}
		mimeHeader.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"", v[1]))
		part, _ := writer.CreatePart(mimeHeader)

		variableValue, _ := findVariable(v[2], &htmlString)
		if variableValue != "" {
			v[2] = variableValue
		}

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

	respBody, err := io.ReadAll(res.Body)
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

	streams, err := request.ExtractHLS(sources.Data.Sources[0].Src, nil)
	if err != nil {
		return nil, err
	}

	for _, stream := range streams {
		stream.Ext = "mp4"
	}

	return []*static.Data{
		{
			Site:    site,
			Title:   "jwplayer video",
			Type:    static.DataTypeVideo,
			Streams: streams,
		},
	}, nil
}

// FindJWPlayerURL in HTML page
func FindJWPlayerURL(htmlString *string) string {
	return reJWPlayerURL.FindString(*htmlString)
}

func findVariable(variable string, htmlString *string) (string, error) {
	variable = strings.ReplaceAll(variable, "$", "")
	variable = strings.ReplaceAll(variable, "{", "")
	variable = strings.ReplaceAll(variable, "}", "")

	re := regexp.MustCompile(fmt.Sprintf(findVarible, variable))
	matchedVariable := re.FindStringSubmatch(*htmlString)
	if len(matchedVariable) < 1 {
		return "", fmt.Errorf("could not match any for variable '%s'", variable)
	}

	return matchedVariable[1], nil
}
