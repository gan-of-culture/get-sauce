package hentainexus

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const domain = "hentainexus.com"
const site = "https://" + domain

var reReader = regexp.MustCompile(`initReader\("([^"]+)`)

type gallery []struct {
	Image    string `json:"image"`
	Label    string `json:"label"`
	URLLabel string `json:"url_label"`
	Type     string `json:"type"`
}

type extractor struct{}

// Extract implements static.Extractor.
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

// New returns a hentainexus extractor.
func New() static.Extractor {
	return &extractor{}
}

func parseURL(URL string) []string {
	urlPrefix, err := url.JoinPath(site, "view")
	if err != nil {
		return nil
	}
	if strings.HasPrefix(URL, urlPrefix) {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	var out []string
	re := regexp.MustCompile(`<a href="(/view/\d+)">`)
	for _, matchURLPart := range re.FindAllStringSubmatch(htmlString, -1) {
		out = append(out, site+matchURLPart[1])
	}

	return out
}

func extractData(URL string) (*static.Data, error) {
	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	title := utils.GetMeta(&htmlString, "og:title")

	firstPageURL, err := url.JoinPath(strings.Replace(URL, "view", "read", 1), "001")
	if err != nil {
		return nil, err
	}
	htmlString, err = request.Get(firstPageURL)
	if err != nil {
		return nil, err
	}

	base64String := utils.GetLastItemString(reReader.FindStringSubmatch(htmlString))
	jsonString, err := decryptJson(base64String)
	if err != nil {
		return nil, err
	}

	var gallery gallery
	err = json.Unmarshal([]byte(*jsonString), &gallery)
	if err != nil {
		return nil, err
	}
	var URLs []*static.URL
	wantedPages := utils.NeedDownloadList(len(gallery))
	for _, pIdx := range wantedPages {
		page := gallery[pIdx]
		URLs = append(URLs, &static.URL{
			URL: page.Image,
			Ext: utils.GetFileExt(page.Image),
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
		URL: URL,
	}, nil
}

func decryptJson(base64String string) (*string, error) {
	base64Decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	var origData []rune
	for _, byte := range base64Decoded {
		origData = append(origData, rune(byte))
	}

	domainLength := len(domain)
	loopBound := min(domainLength, 64)
	for i := range loopBound {
		origData[i] = rune(byte(origData[i]) ^ byte(domain[i]))
	}

	// a loop creates the first 16 prime numbers using sieve (upto 256)
	primes := []uint{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}
	// Compute hash for prime index
	var hashVal uint
	//for _, r := range origData {
	//hashVal ^= uint(byte(r))
	for i := range 64 {
		hashVal ^= uint(byte(origData[i]))
		for range 8 {
			if hashVal&1 != 0 {
				hashVal = (hashVal >> 1) ^ 12
				continue
			}
			hashVal >>= 1
		}
	}
	hashVal &= 7
	selectedPrime := primes[hashVal]

	// KSA: Initialize and key schedule S-box
	sbox := make([]byte, 256)
	for i := range sbox {
		sbox[i] = byte(i)
	}
	var j uint
	for i := range 256 {
		j = (j + uint(sbox[i]) + uint(origData[i%64])) % 256
		sbox[i], sbox[j] = sbox[j], sbox[i]
	}

	// PRGA-like decryption loop
	var (
		output string
		i      uint
		jVal   uint
		l      uint
		k      uint
		l_idx  uint
	)
	for ; l_idx+64 < uint(len(origData)); l_idx++ {
		// Update i
		i = (i + selectedPrime) % 256

		// Update j
		tempJ := (jVal + uint(sbox[i])) % 256
		sTemp := uint(sbox[tempJ])
		jVal = (l + sTemp) % 256

		// Update l
		l = (l + i + uint(sbox[i])) % 256

		// Swap in sbox
		sbox[i], sbox[jVal] = sbox[jVal], sbox[i]

		// Compute keystream k (nested)
		inner1 := (k + l) % 256
		s_inner1 := sbox[inner1]
		inner2 := (i + uint(s_inner1)) % 256
		s_inner2 := sbox[inner2]
		inner3 := (jVal + uint(s_inner2)) % 256
		k = uint(sbox[inner3])

		decypted := byte(origData[l_idx+64]) ^ byte(k)
		output += string(rune(decypted))
	}

	return &output, nil
}
