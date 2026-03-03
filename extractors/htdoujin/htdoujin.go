package htdoujin

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

/*
I have noticed there there are some doujin sites with the same site design just different branding
because they linked to some of the htstreaming sites I called this extractor htdoujin this might change in the future
*/

type CDNDetermenationType string

const (
	// Simple means there is only one CDN prefix
	Simple CDNDetermenationType = "simple"
	// ServerID is used to get the correct CDN prefix
	ServerID CDNDetermenationType = "server_id"
	// uID is used to get the correct CDN prefix
	UID CDNDetermenationType = "u_id"
	// Unknown CDN prefix identifier
	Unknown CDNDetermenationType = "unknown"
)

type siteConfig struct {
	BaseURL string
	CDNDetermenationType
	CDNPrefix           string
	CDNPrefixLevels     []int
	CDNPrefixSrcURLPart string
	ImageExt            string
	GalleryPrefix       string
	ReaderURLPrefix     string
}

const defaultGalleryPrefix = "gallery"

var sites map[string]siteConfig = map[string]siteConfig{
	"asmhentai.com": {
		CDNPrefix:       "images",
		GalleryPrefix:   "g",
		ImageExt:        "jpg",
		ReaderURLPrefix: "gallery",
	},
	"comicporn.xxx": {
		ReaderURLPrefix: "view",
	},
	"hentaienvy.com": {
		ReaderURLPrefix: "g",
	},
	"hentaiera.com": {
		ReaderURLPrefix: "view",
	},
	"hentaifox.com": {
		CDNPrefix:           "i",
		CDNPrefixSrcURLPart: "i",
		ReaderURLPrefix:     "g",
	},
	"hentairox.com": {
		ReaderURLPrefix: "view",
	},
	"hentaizap.com": {
		ReaderURLPrefix: "g",
	},
	"imhentai.xxx": {
		ReaderURLPrefix: "view",
	},
}

var extensionMap = map[string]string{
	"j": "jpg",
	"p": "png",
	"b": "bmp",
	"g": "gif",
	"w": "webp",
}

var reMainJsPath *regexp.Regexp = regexp.MustCompile(`js/main[_\.]?\w*\.js`)
var reUIDLevels *regexp.Regexp = regexp.MustCompile(`u_id\s*>\s*(\d+)`)
var reTitle *regexp.Regexp = regexp.MustCompile(`<title>(.+)</title>`)
var reJSONData *regexp.Regexp = regexp.MustCompile(`'{[^']+`)
var reImgDir *regexp.Regexp = regexp.MustCompile(`image_dir" value="([^"]*)`)
var reGalleryID *regexp.Regexp = regexp.MustCompile(`gallery_id" value="([^"]*)`)
var reUID *regexp.Regexp = regexp.MustCompile(`u_id" value="([^"]*)`)
var reServerID *regexp.Regexp = regexp.MustCompile(`server_id" value="([^"]*)`)
var reServerIDLevels *regexp.Regexp = regexp.MustCompile(`server_id\s*==\s*(\d+)`)
var rePages = regexp.MustCompile(`pages" value="([^"]*)`)

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

	siteCfg, ok := sites[u.Host]
	if !ok {
		return nil, errors.New("site not configured for htdoujin extractor")
	}
	siteCfg.GalleryPrefix = cmp.Or(siteCfg.GalleryPrefix, defaultGalleryPrefix)
	siteCfg.BaseURL = cmp.Or(siteCfg.BaseURL, fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if siteCfg.CDNDetermenationType == "" {
		siteCfg.CDNPrefixLevels, siteCfg.CDNDetermenationType, err = parseCDNPrefixLevels(siteCfg)
		if err != nil {
			return nil, err
		}
	}
	// set prased/default values for next calls of Extractor
	sites[u.Host] = siteCfg

	IDs := parseURL(URL, siteCfg)
	if len(IDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, id := range IDs {
		d, err := extractData(id, siteCfg)
		if err != nil {
			return nil, utils.Wrap(err, id)
		}
		data = append(data, d)
	}

	return data, nil
}

func parseURL(URL string, siteCfg siteConfig) []string {
	galleryPrefixURL, _ := url.JoinPath(siteCfg.BaseURL, siteCfg.GalleryPrefix)
	if urlPart, ok := strings.CutPrefix(URL, galleryPrefixURL); ok {
		return strings.Split(strings.TrimPrefix(urlPart, "/"), "/")[:1]
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	var reGID *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`/%s/(\d+)/`, siteCfg.GalleryPrefix))
	IDs := []string{}
	for _, v := range reGID.FindAllStringSubmatch(htmlString, -1) {
		IDs = append(IDs, utils.GetLastItemString(v))
	}

	return utils.RemoveAdjDuplicates(IDs)
}

func extractData(ID string, siteCfg siteConfig) (*static.Data, error) {
	readerURL, err := url.JoinPath(siteCfg.BaseURL, siteCfg.ReaderURLPrefix, ID, "1/")
	if err != nil {
		return nil, err
	}
	htmlString, err := request.Get(readerURL)
	if err != nil {
		return nil, err
	}

	title := strings.Split(strings.Split(reTitle.FindStringSubmatch(htmlString)[1], " - Page 1 - ")[0], " Page 1 -")[0]

	jsonString := strings.Trim(reJSONData.FindString(htmlString), "'")

	gData := map[string]string{}
	// ignore error here. AsmHentai has no gData container
	_ = json.Unmarshal([]byte(jsonString), &gData)

	imageDir := reImgDir.FindStringSubmatch(htmlString)
	if len(imageDir) < 1 {
		return nil, errors.New("cannot find image_dir for")
	}

	gID := reGalleryID.FindStringSubmatch(htmlString)
	if len(gID) < 1 {
		return nil, errors.New("cannot find gallery_id for")
	}

	prefixSelectionID := utils.GetLastItemString(reUID.FindStringSubmatch(htmlString))
	if prefixSelectionID == "" {
		prefixSelectionID = utils.GetLastItemString(reServerID.FindStringSubmatch(htmlString))
	}

	CDNPrefix, err := getCDNPrefix(prefixSelectionID, siteCfg)
	if err != nil {
		return nil, err
	}

	cdnURL, err := url.Parse(siteCfg.BaseURL)
	if err != nil {
		return nil, err
	}
	cdnURL.Host = CDNPrefix + "." + cdnURL.Host

	pagesCount := len(gData)
	if pagesCount < 1 {
		matchedPages := rePages.FindStringSubmatch(htmlString)
		if len(imageDir) < 1 {
			return nil, errors.New("cannot find pages for")
		}
		pagesCount, err = strconv.Atoi(utils.GetLastItemString(matchedPages))
		if err != nil {
			return nil, err
		}
	}
	pages := utils.NeedDownloadList(pagesCount)

	var ok bool
	URLs := []*static.URL{}
	for _, i := range pages {
		ext := siteCfg.ImageExt
		if ext == "" {
			ext = strings.Split(gData[fmt.Sprint(i+1)], ",")[0] //type, width, height
			ext, ok = extensionMap[ext]
			if !ok {
				return nil, fmt.Errorf("extension %s cannot be mapped", ext)
			}
		}
		imagePath, err := url.JoinPath(imageDir[1], gID[1], fmt.Sprintf("%d.%s", i+1, ext))
		if err != nil {
			return nil, err
		}
		cdnURL.Path = imagePath

		URLs = append(URLs, &static.URL{
			URL: cdnURL.String(),
			Ext: ext,
		})
	}

	galleryURL, err := url.JoinPath(siteCfg.BaseURL, siteCfg.GalleryPrefix, ID)
	if err != nil {
		return nil, err
	}

	return &static.Data{
		Site:  siteCfg.BaseURL,
		Title: title,
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: galleryURL,
	}, nil
}

func getCDNPrefix(prefixSelectionID string, siteCfg siteConfig) (string, error) {
	if siteCfg.CDNDetermenationType == Simple {
		return siteCfg.CDNPrefix, nil
	}

	IDAsNumber, err := strconv.Atoi(prefixSelectionID)
	if err != nil {
		return "", err
	}

	for i := len(siteCfg.CDNPrefixLevels); i >= 1; i-- {
		switch siteCfg.CDNDetermenationType {
		case ServerID:
			if IDAsNumber == siteCfg.CDNPrefixLevels[i-1] {
				return fmt.Sprintf("m%d", i), nil
			}
		default:
			if IDAsNumber > siteCfg.CDNPrefixLevels[i-1] {
				return fmt.Sprintf("m%d", i), nil
			}
		}
	}

	return "", errors.New("no CDN prefix was found. Check if CDNPrefixLevels have been parsed correctly")
}

func parseCDNPrefixLevels(siteCfg siteConfig) ([]int, CDNDetermenationType, error) {
	if siteCfg.CDNPrefix != "" {
		return nil, Simple, nil
	}

	htmlString, err := request.Get(siteCfg.BaseURL)
	if err != nil {
		return nil, Unknown, err
	}

	mainJsPathPart := reMainJsPath.FindString(htmlString)
	if mainJsPathPart == "" {
		return nil, Unknown, errors.New("no main_*.js file was found linked from the homepage")
	}

	mainJSPath, err := url.JoinPath(siteCfg.BaseURL, mainJsPathPart)
	if err != nil {
		return nil, Unknown, err
	}

	jsString, err := request.Get(mainJSPath)
	if err != nil {
		return nil, Unknown, err
	}

	cdID := UID
	levels := reUIDLevels.FindAllStringSubmatch(jsString, -1)
	if len(levels) == 0 {
		cdID = ServerID
		levels = reServerIDLevels.FindAllStringSubmatch(jsString, -1)
	}

	var out []int
	for _, uID := range levels {
		IDAsNumber, err := strconv.Atoi(uID[1])
		if err != nil {
			return nil, cdID, err
		}
		out = append(out, IDAsNumber)
	}

	return out, cdID, nil
}
