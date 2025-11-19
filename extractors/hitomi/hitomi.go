package hitomi

import (
	"cmp"
	"encoding/binary"
	"encoding/json"
	"fmt"
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
	Date              string          `json:"date"`
	Files             []img           `json:"files"`
	JapaneseTitle     string          `json:"japanese_title"`
	Type              string          `json:"type"`
	ID                json.RawMessage `json:"id"`
	Tags              []tag           `json:"tags"`
	LanguageLocalName string          `json:"language_localname"`
	Language          string          `json:"language"`
	Title             string          `json:"title"`
}

const site = "https://hitomi.la/"
const domain2 = "gold-usergeneratedcontent.net"
const nozomi = "https://ltn.gold-usergeneratedcontent.net/" // is domain for gg.js
const readerURL = "https://hitomi.la/reader/"
const nozomiExt = "nozomi"
const galleriesPerPage = 25
const ggURL = nozomi + "gg.js"

var reSubdomainPart = regexp.MustCompile(`\/[0-9a-f]{61}([0-9a-f]{2})([0-9a-f])`)
var reURLFromURL = regexp.MustCompile(`\/\/..?\.(?:gold-usergeneratedcontent\.net|hitomi\.la)\/`)
var rePathFromHash = regexp.MustCompile(`(..)(.)$`)

var ggValues []*int
var b string
var ggMatchedValue int
var ggNonMatchedValue int

type extractor struct{}

// New returns a hitomi extractor
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	err := initGGValues()
	if err != nil {
		return nil, err
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
	if ok, _ := regexp.MatchString(fmt.Sprintf("%s(?:manga|doujinshi|cg|gamecg|imageset)/", site), URL); ok {
		re := regexp.MustCompile(`(\d*).html#*\d*$`)
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

	// from galleryblock.js func fetchnozomi
	startByte := (pageNumber - 1) * galleriesPerPage * 4
	endByte := startByte + galleriesPerPage*4 - 1
	htmlData, err := request.GetAsBytesWithHeaders(nozomiURL, map[string]string{
		"Range": fmt.Sprintf("bytes=%d-%d", startByte, endByte),
	})
	if err != nil {
		return nil
	}

	URLs := []string{}
	from := 0
	for i := 4; i <= int(len(htmlData)); i += 4 {
		URLs = append(URLs, fmt.Sprintf("%sgalleries/%d.js", nozomi, binary.BigEndian.Uint32(htmlData[from:i])))
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

	galleryData := gallery{}
	err = json.Unmarshal([]byte(jsString[jsonStart:]), &galleryData)
	if err != nil {
		return nil, err
	}

	dir := ""
	u := ""
	imgFile := img{}
	URLs := []*static.URL{}
	pages := utils.NeedDownloadList(len(galleryData.Files))
	for _, pageIdx := range pages {
		dir = ""
		imgFile = galleryData.Files[pageIdx]
		if imgFile.HasWebp == 1 || imgFile.HasAVIF == 1 {
			dir = "avif"
		}
		u = urlFromURL(urlFromHash(imgFile, dir), dir)
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
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: fmt.Sprintf("%s%s.html", readerURL, string(galleryData.ID)),
	}, nil
}

func subdomainFromURL(URL, base, dir string) string {
	retval := ""
	if base == "" {
		switch dir {
		case "webp":
			retval = "w"
		case "avif":
			retval = "a"
		}
	}

	b := 16

	var m = reSubdomainPart.FindStringSubmatch(URL)
	if len(m) == 0 {
		return retval
	}

	g, err := strconv.ParseInt(m[2]+m[1], b, 64)
	if err == nil {
		if base != "" {
			retval = fmt.Sprintf("%c", 97+inGGValues(int(g))) + base
		} else {
			retval = retval + fmt.Sprint((1 + inGGValues(int(g))))
		}
	}

	return retval
}

func urlFromURL(URL, dir string) string {
	return reURLFromURL.ReplaceAllString(URL, fmt.Sprintf("//%s.%s/", subdomainFromURL(URL, "", dir), domain2))
}

func urlFromHash(imgFile img, dir string) string {
	ext := cmp.Or(dir, strings.Split(imgFile.Name, ".")[1])
	if dir == "webp" || dir == "avif" {
		dir = ""
	} else {
		dir += "/"
	}

	p, err := url.JoinPath(fmt.Sprintf("https://a.%s", domain2), dir, fullPathFromHash(imgFile.Hash)+"."+ext)
	if err != nil {
		return ""
	}
	return p
}

func fullPathFromHash(hash string) string {
	m := rePathFromHash.FindStringSubmatch(hash)

	dec, _ := strconv.ParseInt(m[2]+m[1], 16, 64)
	return fmt.Sprintf("%s%d/%s", b, dec, hash)
}

func initGGValues() error {
	jsStr, _ := request.GetWithHeaders(ggURL, map[string]string{"Referer": site})

	b = regexp.MustCompile(`\d+/`).FindString(jsStr)

	matchedLimitValues := regexp.MustCompile(`\d;`).FindAllString(jsStr, -1)
	if len(matchedLimitValues) < 2 {
		return fmt.Errorf("no limit values found in: %s", ggURL)
	}

	var err error
	ggNonMatchedValue, err = strconv.Atoi(strings.Trim(matchedLimitValues[0], ";"))
	if err != nil {
		return err
	}
	ggMatchedValue, _ = strconv.Atoi(strings.Trim(matchedLimitValues[1], ";"))

	re := regexp.MustCompile(`case (\d+)`)
	matchedCases := re.FindAllStringSubmatch(jsStr, -1)
	ggValues = make([]*int, len(matchedCases))
	for idx, num := range matchedCases {
		n, _ := strconv.Atoi(num[1])

		ggValues[idx] = &n
	}

	return nil
}

func inGGValues(num int) int {
	if ggValues == nil {
		return ggNonMatchedValue
	}
	for _, value := range ggValues {
		if *value == num {
			return ggMatchedValue
		}
	}
	return ggNonMatchedValue
}
