package htdoujin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

/*
I have noticed there there are some doujin sites with the same site design just different branding
because they linked to some of the htstreaming sites I called this extractor htdoujin this might change in the future
*/

type siteConfig struct {
	CDNPrefix       string
	ReaderURLPrefix string
}

var sites map[string]siteConfig = map[string]siteConfig{
	"comicporn.xxx": {
		CDNPrefix:       "m5",
		ReaderURLPrefix: "view",
	},
	"imhentai.xxx": {
		CDNPrefix:       "m5",
		ReaderURLPrefix: "view",
	},
	"hentaiera.com": {
		CDNPrefix:       "m1",
		ReaderURLPrefix: "view",
	},
	"hentaifox.com": {
		CDNPrefix:       "i",
		ReaderURLPrefix: "g",
	},
	"hentairox.com": {
		CDNPrefix:       "m5",
		ReaderURLPrefix: "view",
	},
}

var site string
var cdn string
var readerURLPrefix string

var reUID *regexp.Regexp = regexp.MustCompile(`/gallery/(\d+)/`)
var reTitle *regexp.Regexp = regexp.MustCompile(`<title>(.+)</title>`)
var reJSONData *regexp.Regexp = regexp.MustCompile(`'{[^']+`)
var reImgDir *regexp.Regexp = regexp.MustCompile(`image_dir" value="([^"]*)`)
var reGalleryID *regexp.Regexp = regexp.MustCompile(`gallery_id" value="([^"]*)`)

type extractor struct{}

// New returns a htdoujin extractor.
func New() static.Extractor {
	return &extractor{}
}

func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	if _, ok := sites[u.Host]; !ok {
		return nil, errors.New("site not configured for htdoujin extractor")
	}
	site = "https://" + u.Host + "/"
	cdn = fmt.Sprintf("https://%s.%s/", sites[u.Host].CDNPrefix, u.Host)
	readerURLPrefix = sites[u.Host].ReaderURLPrefix

	IDs := parseURL(URL)
	if len(IDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, id := range IDs {
		d, err := extractData(id)
		if err != nil {
			return nil, utils.Wrap(err, id)
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string) []string {
	if strings.HasPrefix(URL, site+"gallery/") {
		return []string{URL[len(site+"gallery/") : len(URL)-1]}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	IDs := []string{}
	for _, v := range reUID.FindAllStringSubmatch(htmlString, -1) {
		IDs = append(IDs, v[1])
	}

	return utils.RemoveAdjDuplicates(IDs)
}

func extractData(ID string) (*static.Data, error) {
	htmlString, err := request.Get(fmt.Sprintf("%s%s/%s/1/", site, readerURLPrefix, ID))
	if err != nil {
		return &static.Data{}, err
	}

	title := strings.Split(reTitle.FindStringSubmatch(htmlString)[1], " - Page 1 - ")[0]

	jsonString := strings.Trim(reJSONData.FindString(htmlString), "'")
	//fmt.Println(jsonString)

	gData := map[string]string{}
	err = json.Unmarshal([]byte(jsonString), &gData)
	if err != nil {
		return &static.Data{}, err
	}

	imageDir := reImgDir.FindStringSubmatch(htmlString)
	if len(imageDir) < 1 {
		return &static.Data{}, errors.New("cannot find image_dir for")
	}

	gID := reGalleryID.FindStringSubmatch(htmlString)
	if len(gID) < 1 {
		return &static.Data{}, errors.New("cannot find gallery_id for")
	}

	pages := utils.NeedDownloadList(len(gData))

	URLs := []*static.URL{}
	for _, i := range pages {
		params := strings.Split(gData[fmt.Sprint(i)], ",") //type, width, height
		switch params[0] {
		case "j":
			params[0] = "jpg"
		case "p":
			params[0] = "png"
		case "b":
			params[0] = "bmp"
		case "g":
			params[0] = "gif"
		}
		URLs = append(URLs, &static.URL{
			URL: fmt.Sprintf("%s%s/%s/%d.%s", cdn, imageDir[1], gID[1], i, params[0]),
			Ext: params[0],
		})
	}

	return &static.Data{
		Site:  site,
		Title: title,
		Type:  "image",
		Streams: map[string]*static.Stream{
			"0": {
				URLs: URLs,
			},
		},
		URL: fmt.Sprintf("%sgallery/%s/", site, ID),
	}, nil
}
