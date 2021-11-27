package damn

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://www.damn.stream/"
const embed = "https://www.damn.stream/embed/"

var reVideoID *regexp.Regexp = regexp.MustCompile(`"/embed/([^"]*)`)
var reSrcMeta *regexp.Regexp = regexp.MustCompile(`<source\s[^=]*="([^"]+)"(?:[\s\S]*?size="([^"]*))?`)
var reVideoMeta *regexp.Regexp = regexp.MustCompile(`[^"]+/v/hentai[^"]+`)

//const cdn = "https://server-one.damn.stream"

type extractor struct{}

// New returns a damn.stream extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) < 1 {
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
	if strings.HasPrefix(URL, "https://www.damn.stream/watch/") {
		return []string{URL}
	}

	if !strings.HasPrefix(URL, "https://www.damn.stream/hentai/") {
		return []string{}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://www.damn.stream/watch/[^"]*`)
	return utils.RemoveAdjDuplicates(re.FindAllString(htmlString, -1))
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := strings.Trim(utils.GetH1(&htmlString, -1), "\n- ")
	videoID := reVideoID.FindStringSubmatch(htmlString)[1]

	htmlString, err = request.Get(fmt.Sprintf("%s%s", embed, videoID))
	if err != nil {
		return nil, err
	}

	srcMeta := reSrcMeta.FindStringSubmatch(htmlString) //1=URL 2=ext
	if len(srcMeta) == 0 {
		srcMeta = append(srcMeta, "")
		srcMeta = append(srcMeta, reVideoMeta.FindAllString(htmlString, -1)...)
	}

	quality := ""
	if len(srcMeta) > 2 {
		quality = srcMeta[2]
	}

	size, _ := request.Size(srcMeta[1], site)

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  "video",
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{
					0: {
						URL: srcMeta[1],
						Ext: utils.GetFileExt(srcMeta[1]),
					},
				},
				Quality: quality + "p",
				Size:    size,
			},
		},
		URL: URL,
	}, err
}
