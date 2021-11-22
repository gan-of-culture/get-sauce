package hitomi

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type img struct {
	Hash    string `json:"hash"`
	HasAVIF int    `json:"hasavif"`
	HasWebp int    `json:"haswebp"`
	Width   uint32 `json:"width"`
	Height  uint32 `json:"height"`
	Name    string `json:"name"`
}

type tag struct {
	URL string `json:"url"`
	Tag string `json:"tag"`
}

type gallery struct {
	Date              string `json:"date"`
	Files             []img  `json:"files"`
	JapaneseTitle     string `json:"japanese_title"`
	Type              string `json:"type"`
	ID                string `json:"id"`
	Tags              []tag  `json:"tags"`
	LanguageLocalName string `json:"language_localname"`
	Language          string `json:"language"`
	Title             string `json:"title"`
}

const site = "https://hitomi.la/"
const nozomi = "https://ltn.hitomi.la/"
const nozomiExt = "nozomi"
const galleriesPerPage = 25

var reSubdomainPart = regexp.MustCompile(`\/[0-9a-f]\/([0-9a-f]{2})\/`)
var reURLFromURL = regexp.MustCompile(`\/\/..?\.hitomi\.la\/`)
var rePathFromHash = regexp.MustCompile(`^.*(..)(.)$`)

type extractor struct{}

// New returns a hitomi extractor.
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
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if ok, _ := regexp.MatchString(fmt.Sprintf("%s(?:manga|doujinshi)/", site), URL); ok {
		re := regexp.MustCompile(`(\d*).html$`)
		id := re.FindStringSubmatch(URL)[1]
		return []string{fmt.Sprintf("%sgalleries/%s.js", nozomi, id)}
	}

	u, err := url.Parse(URL)
	if err != nil {
		return []string{}
	}

	if !strings.HasSuffix(u.Path, ".html") {
		return []string{}
	}

	nozomiURL := nozomi + strings.TrimSuffix(u.Path, "html") + nozomiExt

	re := regexp.MustCompile(`page=(\d+)$`)
	pageNumber := 1
	matchedPageNumber := re.FindStringSubmatch(URL)
	if len(matchedPageNumber) >= 2 {
		pageNumber, _ = strconv.Atoi(matchedPageNumber[1])
	}

	startByte := (pageNumber - 1) * galleriesPerPage * 4
	endByte := startByte + galleriesPerPage*4 - 1
	resp, err := request.Request(http.MethodGet, nozomiURL, map[string]string{
		"Range": fmt.Sprintf("bytes=%d-%d", startByte, endByte),
	}, nil)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}
	}

	URLs := []string{}
	from := 0
	for i := 4; i <= int(resp.ContentLength); i += 4 {
		URLs = append(URLs, fmt.Sprintf("%sgalleries/%d.js", nozomi, binary.BigEndian.Uint32(buffer[from:i])))
		from = i
	}
	return URLs
}

func extractData(URL string) (*static.Data, error) {
	jsString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	jsonStart := strings.Index(jsString, "{")
	if err != nil {
		return nil, errors.New("no json string found for")
	}

	galleryData := gallery{}
	err = json.Unmarshal([]byte(jsString[jsonStart:]), &galleryData)
	if err != nil {
		log.Println(URL)
		return nil, err
	}

	base := ""
	u := ""
	imgFile := img{}
	URLs := []*static.URL{}
	pages := utils.NeedDownloadList(len(galleryData.Files))
	for _, pageIdx := range pages {
		base = ""
		imgFile = galleryData.Files[pageIdx-1]
		if imgFile.HasWebp == 1 || imgFile.HasAVIF == 1 {
			base = "a"
		}
		u = urlFromURL(urlFromHash(imgFile), base)
		URLs = append(URLs, &static.URL{
			URL: u,
			Ext: utils.GetLastItemString(strings.Split(u, ".")),
		})
	}

	return &static.Data{
		Site:  site,
		Title: galleryData.Title,
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: URLs,
			},
		},
		URL: URL,
	}, nil
}

func subdomainFromGalleryid(g, numOfFrontends int64) string {
	o := g % numOfFrontends
	return fmt.Sprintf("%c", 97+o)
}

func subdomainFromURL(URL, base string) string {
	retval := "b"
	if base != "" {
		retval = base
	}

	numberOfFrontends := 3
	b := 16

	var m = reSubdomainPart.FindStringSubmatch(URL)
	if len(m) < 2 {
		return "a"
	}

	g, err := strconv.ParseInt(m[1], b, 64)
	if err == nil {
		// check these values if it doesn't work anymore
		if g < 0x80 {
			numberOfFrontends = 2
		}
		if g < 0x59 {
			g = 1
		}
		retval = subdomainFromGalleryid(g, int64(numberOfFrontends)) + retval
	}

	return retval
}

func urlFromURL(URL, base string) string {
	return reURLFromURL.ReplaceAllString(URL, fmt.Sprintf("//%s.hitomi.la/", subdomainFromURL(URL, base)))
}

func urlFromHash(imgFile img) string {
	dir := "images"
	if imgFile.HasWebp == 1 {
		dir = "webp"
	}
	if imgFile.HasAVIF == 1 {
		dir = "avif"
	}
	ext := dir
	if ext == "images" {
		ext = strings.Split(imgFile.Name, ".")[1]
	}

	return "https://a.hitomi.la/" + dir + "/" + fullPathFromHash(imgFile.Hash) + "." + ext
}

func fullPathFromHash(hash string) string {
	if len(hash) < 3 {
		return hash
	}
	return rePathFromHash.ReplaceAllString(hash, "$2/$1/"+hash)
}
