package danbooru

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
)

const site = "https://danbooru.donmai.us"
const postURLTemplate = site + "/posts/%s"

// [1] = img original width [2] image original height [3] src URL [4] image name
var reIMGData = regexp.MustCompile(`data-width="([^"]+)"[ ]+data-height="([^"]+)"[\s\S]*?alt="([^"]+)".+src="([^"]+)"`)

// src URL, size, width, height
var reIMGData2 = regexp.MustCompile(`Size: <a href="([^"]+)">([^A-Z]+[^\s]+)[^(]+\(([^)]+)`)

type extractor struct {
	client *http.Client
}

// New returns a danbooru extractor
func New() static.Extractor {
	return newForTesting()
}

func newForTesting() *extractor {
	return &extractor{client: request.Firefox117Client()}
}

// Extract for danbooru pages
func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	postIDs, err := e.parseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []*static.Data{}
	for _, postID := range postIDs {
		contentData, err := e.extractData(postID)
		if err != nil {
			log.Printf(postURLTemplate, postID)
			return nil, err
		}
		data = append(data, contentData)
	}

	return data, nil
}

// parseURL for danbooru pages
func (e *extractor) parseURL(URL string) ([]string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pathSplit := strings.Split(u.Path, "/")
	lenPathSplit := len(pathSplit)
	if lenPathSplit >= 2 && pathSplit[lenPathSplit-2] == "posts" {
		return []string{pathSplit[lenPathSplit-1]}, nil
	}

	htmlString, err := request.GetAsBytesWithClient(e.client, URL, URL)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`data-id="([^"]+)`)
	out := []string{}
	for _, submatchID := range re.FindAllSubmatch(htmlString, -1) {
		out = append(out, string(submatchID[1]))
	}

	return out, nil
}

func (e *extractor) extractData(postID string) (*static.Data, error) {
	postURL := fmt.Sprintf(postURLTemplate, postID)
	htmlBytes, err := request.GetAsBytesWithClient(e.client, postURL, postURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	matchedImgData := reIMGData2.FindStringSubmatch(string(htmlBytes))
	if len(matchedImgData) != 4 {
		log.Println(string(htmlBytes))
		return nil, errors.WithStack(static.ErrDataSourceParseFailed)
	}

	_, fname := path.Split(matchedImgData[1])
	fnameSplit := strings.Split(fname, "_")
	fname = strings.Trim(strings.Join(fnameSplit[:len(fnameSplit)-2], "_"), "_")
	size, unit, ok := strings.Cut(matchedImgData[2], " ")
	var s int64
	if ok {
		sizeF, err := strconv.ParseFloat(size, 10)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		s = utils.CalcSizeInByte(sizeF, unit)
	}

	return &static.Data{
		Site:  site,
		Title: fmt.Sprintf("%s_%s", fname, postID),
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: []*static.URL{
					{
						URL: matchedImgData[1],
						Ext: utils.GetFileExt(matchedImgData[1]),
					},
				},
				Quality: matchedImgData[3],
				Size:    s,
			},
		},
		URL: postID,
	}, nil
}
