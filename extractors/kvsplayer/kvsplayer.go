package kvsplayer

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
)

/*

	Special Thanks to rigstot https://github.com/ytdl-patched/ytdl-patched/commits?author=rigstot for the python based template https://github.com/ytdl-patched/ytdl-patched/commit/a318f59d14792d25b2206c3f50181e03e8716db7
	Further documentation on the kvsplayer can be found here: https://www.kernel-scripts.com/en/documentation/player/
*/

var reHasKVSPlayer = regexp.MustCompile(`<script [^>]*?src="https://.+?/kt_player\.js\?v=(?P<ver>(?P<maj_ver>\d+)(\.\d+)+)".*?>`)
var reFlashVars = regexp.MustCompile(`var\s+(?:flashvars|[\w]{11})\s*=[\s\S]*?};`)
var reFlashVarsValues = regexp.MustCompile(`(\w+): ['"](.*?)['"],`)
var reTitle = regexp.MustCompile(`<link href="https?://[^"]+/(.+?)/?" rel="canonical"\s*/?>`)

type extractor struct{}

// New returns a kvsplayer extractor.
func New() static.Extractor {
	return &extractor{}
}

// Extract from URL
func (e *extractor) Extract(URL string) ([]*static.Data, error) {

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	data, err := ExtractFromHTML(&htmlString)
	if err != nil {
		return nil, err
	}

	data[0].URL = URL
	return data, nil
}

// ExtractFromHTML as usable utility for other scrapers
func ExtractFromHTML(htmlString *string) ([]*static.Data, error) {

	matchedKVSPlayer := reHasKVSPlayer.FindAllStringSubmatch(*htmlString, -1) // 1=Vers 2=Major Version 3=Sub Version
	if len(matchedKVSPlayer) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	switch matchedKVSPlayer[0][2] {
	case "4", "5", "8", "9", "11", "12", "13", "15":
		break
	default:
		fmt.Printf("Untested major version (%s) in player engine--Download may fail.", matchedKVSPlayer[0][2])
	}

	flashvars, err := parseFlashVars(htmlString)
	if err != nil {
		return nil, err
	}

	matchedTitle := reTitle.FindStringSubmatch(*htmlString)
	if len(matchedTitle) < 1 {
		return nil, errors.New("no title found in 'cronical' URL link")
	}

	ext := "mp4"
	if val, ok := flashvars["postfix"]; ok {
		splitVal := strings.Split(val, ".") // post.mp4 to mp4
		if len(splitVal) > 1 {
			ext = splitVal[1]
		}
	}

	quality := ""

	streams := map[string]*static.Stream{}
	dataLen := 0 //number of possible streams
	var flashvarsVideoURL = []string{"video_url", "video_alt_url", "video_alt_url2", "video_alt_url3", "video_alt_url4"}
	for _, key := range flashvarsVideoURL {
		if _, ok := flashvars[key]; ok {
			dataLen += 1
		}
	}
	for i, key := range flashvarsVideoURL {
		if _, ok := flashvars[key]; !ok {
			continue
		}
		if !strings.Contains(flashvars[key], "/get_file") {
			continue
		}
		rURL, err := getRealURL(flashvars[key], flashvars["license_code"])
		if err != nil {
			return nil, err
		}
		if val, ok := flashvars[key+"_text"]; ok {
			quality = val
		}

		size, _ := request.Size(rURL, rURL)

		streams[fmt.Sprint(dataLen-i-1)] = &static.Stream{
			Type: static.DataTypeVideo,
			URLs: []*static.URL{
				{
					URL: rURL,
					Ext: ext,
				},
			},
			Quality: quality,
			Size:    size,
		}
	}

	return []*static.Data{
		{
			Site:    "https://kvsplayer.com/",
			Title:   matchedTitle[1],
			Type:    static.DataTypeVideo,
			Streams: streams,
		},
	}, nil
}

func getRealURL(videoURL, licenseCode string) (string, error) {
	if !strings.HasPrefix(videoURL, "function/0/") {
		return videoURL, nil //not obfuscated
	}

	splitURL := strings.Split(videoURL, "?")
	lenSplitURL := len(splitURL)

	URLPath := splitURL[0]
	URLQuery := ""
	if lenSplitURL > 1 {
		URLQuery = splitURL[1]
	}
	URLParts := strings.Split(URLPath, "/")[2:]
	license := getLicenseToken(licenseCode)
	newMagic := URLParts[5][:32]

	for o := len(newMagic) - 1; o > -1; o-- {
		new := ""
		lNum := 0
		nAsInt := 0
		for _, n := range license[o:] {
			nAsInt, _ = strconv.Atoi(string(n)) //conv to single digit string : warn OK
			lNum += nAsInt
		}

		l := (o + lNum) % 32

		for i := 0; i < len(newMagic); i++ {
			if i == o {
				new += string(newMagic[l])
			} else if i == l {
				new += string(newMagic[o])
			} else {
				new += string(newMagic[i])
			}
		}
		newMagic = new
	}

	URLParts[5] = newMagic + URLParts[5][32:]
	return fmt.Sprintf("%s?%s", strings.Join(URLParts, "/"), URLQuery), nil //add ending "?" to ignore the rnd param
}

func getLicenseToken(license string) string {
	modlicense := strings.ReplaceAll(strings.ReplaceAll(license, "$", ""), "0", "1")

	center := len(modlicense) / 2
	fronthalf, err := strconv.Atoi(modlicense[:center+1])
	if err != nil {
		return ""
	}
	backhalf, err := strconv.Atoi(modlicense[center:])
	if err != nil {
		return ""
	}

	modlicense = fmt.Sprint(4 * int(math.Abs(float64(fronthalf-backhalf))))
	retVal := ""
	val1 := 0
	val2 := 0
	for o := 0; o < center+1; o++ {
		for i := 1; i < 5; i++ {
			val1, _ = strconv.Atoi(string(license[o+i]))
			val2, _ = strconv.Atoi(string(modlicense[o]))
			retVal += fmt.Sprint((val1 + val2) % 10)
		}
	}
	return retVal
}

func parseFlashVars(htmlString *string) (map[string]string, error) {
	matchedHtmlFlashvars := reFlashVars.FindAllString(*htmlString, -1)
	slices.SortFunc(matchedHtmlFlashvars, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})
	slices.Reverse(matchedHtmlFlashvars)
	idxFlashVars := slices.IndexFunc(matchedHtmlFlashvars, func(matchedString string) bool {
		return strings.HasPrefix(matchedString, "var flashvars")
	})
	// if one of the matches starts with flashvars -> take that one
	// if said match is already at idx 0 you don't need to reassign
	// if idx == -1 move on normally
	if idxFlashVars > 0 {
		matchedHtmlFlashvars[0] = matchedHtmlFlashvars[idxFlashVars]
	}
	htmlFlashvars := matchedHtmlFlashvars[0]
	if htmlFlashvars == "" {
		return nil, static.ErrDataSourceParseFailed
	}

	matchedFlashVarsValues := reFlashVarsValues.FindAllStringSubmatch(htmlFlashvars, -1) //1=key 2=val
	if len(matchedFlashVarsValues) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}

	flashvars := map[string]string{}
	for _, val := range matchedFlashVarsValues {
		flashvars[val[1]] = val[2]
	}

	return flashvars, nil
}
