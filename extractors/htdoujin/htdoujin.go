package htdoujin

import (
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

type siteConfig struct {
	CDNPrefixSrcURLPart string
	ReaderURLPrefix     string
}

var sites map[string]siteConfig = map[string]siteConfig{
	"comicporn.xxx": {
		ReaderURLPrefix: "view",
	},
	"imhentai.xxx": {
		ReaderURLPrefix: "view",
	},
	"hentaienvy.com": {
		ReaderURLPrefix: "g",
	},
	"hentaiera.com": {
		ReaderURLPrefix: "view",
	},
	"hentaifox.com": {
		CDNPrefixSrcURLPart: "i",
		ReaderURLPrefix:     "g",
	},
	"hentairox.com": {
		ReaderURLPrefix: "view",
	},
	"hentaizap.com": {
		ReaderURLPrefix: "g",
	},
}

var host string
var site string
var cdn string
var cdnDetermenationID CDNDetermenationID
var cdnPrefixLevels []int
var readerURLPrefix string

var reMainJsPath *regexp.Regexp = regexp.MustCompile(`js/main_\w+\.js`)
var reGID *regexp.Regexp = regexp.MustCompile(`/gallery/(\d+)/`)
var reUIDLevels *regexp.Regexp = regexp.MustCompile(`u_id\s*>\s*(\d+)`)
var reTitle *regexp.Regexp = regexp.MustCompile(`<title>(.+)</title>`)
var reJSONData *regexp.Regexp = regexp.MustCompile(`'{[^']+`)
var reImgDir *regexp.Regexp = regexp.MustCompile(`image_dir" value="([^"]*)`)
var reGalleryID *regexp.Regexp = regexp.MustCompile(`gallery_id" value="([^"]*)`)
var reUID *regexp.Regexp = regexp.MustCompile(`u_id" value="([^"]*)`)
var reServerID *regexp.Regexp = regexp.MustCompile(`server_id" value="([^"]*)`)
var reServerIDLevels *regexp.Regexp = regexp.MustCompile(`server_id\s*==\s*(\d+)`)

type CDNDetermenationID string

const (
	// ServerID is used to get the correct CDN prefix
	ServerID CDNDetermenationID = "server_id"
	// uID is used to get the correct CDN prefix
	UID CDNDetermenationID = "u_id"
	// Unknown CDN prefix identifier
	Unknown CDNDetermenationID = "unknown"
)

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

	host = u.Host

	if _, ok := sites[host]; !ok {
		return nil, errors.New("site not configured for htdoujin extractor")
	}
	site = "https://" + host + "/"

	IDs := parseURL(URL)
	if len(IDs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	readerURLPrefix = sites[host].ReaderURLPrefix
	cdnPrefixLevels, cdnDetermenationID, err = parseCDNPrefixLevels()
	if err != nil {
		return nil, err
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
	for _, v := range reGID.FindAllStringSubmatch(htmlString, -1) {
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

	prefixSelectionID := utils.GetLastItemString(reUID.FindStringSubmatch(htmlString))
	if prefixSelectionID == "" {
		prefixSelectionID = utils.GetLastItemString(reServerID.FindStringSubmatch(htmlString))
	}

	CDNPrefix, err := getCDNPrefix(prefixSelectionID, cdnDetermenationID)
	if err != nil {
		return nil, err
	}

	cdn = fmt.Sprintf("https://%s.%s/", CDNPrefix, host)

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
		Type:  static.DataTypeImage,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeImage,
				URLs: URLs,
			},
		},
		URL: fmt.Sprintf("%sgallery/%s/", site, ID),
	}, nil
}

func getCDNPrefix(prefixSelectionID string, cdnDetId CDNDetermenationID) (string, error) {
	if host == "hentaifox.com" {
		return "i", nil
	}

	IDAsNumber, err := strconv.Atoi(prefixSelectionID)
	if err != nil {
		return "", err
	}

	for i := len(cdnPrefixLevels); i >= 1; i-- {
		switch cdnDetId {
		case ServerID:
			if IDAsNumber == cdnPrefixLevels[i-1] {
				return fmt.Sprintf("m%d", i), nil
			}
		default:
			if IDAsNumber > cdnPrefixLevels[i-1] {
				return fmt.Sprintf("m%d", i), nil
			}
		}
	}

	return "", errors.New("no CDN prefix was found. Check if CDNPrefixLevels have been parsed correctly")
}

func parseCDNPrefixLevels() ([]int, CDNDetermenationID, error) {
	htmlString, err := request.Get(site)
	if err != nil {
		return nil, Unknown, err
	}

	mainJsPathPart := reMainJsPath.FindString(htmlString)
	if mainJsPathPart == "" {
		return nil, Unknown, errors.New("no main_*.js file was found linked from the homepage")
	}

	jsString, err := request.Get(site + mainJsPathPart)
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
