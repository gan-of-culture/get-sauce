package exhentai

/*import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://ex-hentai.org/"
const loginFormURL = "https://forums.e-hentai.org/index.php?act=Login&CODE=01"

var headerCookies map[string]string

// LogRedirects to sanitize "Location" URLs
type LogRedirects struct {
	Transport http.RoundTripper
}

//RoundTrip implementaion
func (l LogRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		fmt.Println(resp.Cookies())
	}
	return
}

func login() (map[string]string, error) {

	jar := &request.Myjar{}
	jar.Jar = make(map[string][]*http.Cookie)

	client := &http.Client{
		Transport: LogRedirects{&http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		}},
		Timeout: 15 * time.Minute,
		Jar:     jar,
	}

	headers := map[string]string{
		"cache-control": "no-cache",
		"content-type":  "application/x-www-form-urlencoded",
		"referer":       "https://forums.e-hentai.org/index.php?",
	}

	data := fmt.Sprintf("{ 'CookieDate': '1', 'b': 'd', 'bt': '1-1', 'UserName': '%s', 'PassWord': '%s', 'ipb_login_submit': 'Login!' }", config.Username, config.UserPassword)

	req, err := http.NewRequest(http.MethodPost, loginFormURL, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client.Jar.SetCookies(req.URL, []*http.Cookie{
		{
			Name:  "ipb_coppa",
			Value: "0",
		}, {
			Name:  "ipb_anonlogin",
			Value: "-1",
		}, {
			Name:  "ipb_member_id",
			Value: "0",
		}, {
			Name:  "ipb_pass_hash",
			Value: "0",
		}})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp.Cookies())
	fmt.Println(req.Cookies())
	fmt.Println(resp.Request.Cookies())
	fmt.Println(client.Jar.Cookies(req.URL))

	return nil, nil
}

// Extract data
func Extract(URL string) ([]static.Data, error) {
	headerCookies, _ = login()
	if len(headerCookies) == 0 {
		return nil, errors.New("[ExHentai] can't retrieve login cookies")
	}

	URLs := ParseURL(URL)
	if len(URLs) == 0 {
		return nil, errors.New("[ExHentai] no vaild URL found")
	}

	data := []static.Data{}
	for _, URL := range URLs {
		rData, err := extractData(URL)
		if err != nil {
			return nil, err
		}
		data = append(data, rData...)
	}
	return data, nil
}

// Get content as string
func Get(url string) (string, error) {
	resp, err := request.Request(http.MethodGet, url, headerCookies)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return "", err
		}
	}

	return string(body), nil
}

// ParseURL to gallery URL
func ParseURL(URL string) []string {
	if strings.Contains(URL, "https://exhentai.org/g/") {
		return []string{URL}
	}

	htmlString, err := Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile("https://exhentai.org/g/[^\"]+")
	galleries := re.FindAllStringSubmatch(htmlString, -1)
	if len(galleries) == 0 {
		return []string{}
	}

	out := []string{}

	for _, gallery := range galleries {
		out = append(out, gallery[0])
	}
	return out
}

func extractData(URL string) ([]static.Data, error) {
	htmlString, err := Get(URL)
	if err != nil {
		return nil, errors.New("[ExHentai] invaild URL")
	}

	if strings.Contains(htmlString, "<h1>Content Warning</h1>") {
		if config.RestrictContent {
			return []static.Data{
				{Err: errors.New("[Exhentai] Restricted content")},
			}, nil
		}
		return extractData(URL + "?nw=session")
	}

	re := regexp.MustCompile("([0-9]+) pages")
	htmlNumberOfPages := re.FindStringSubmatch(htmlString)
	if len(htmlNumberOfPages) != 2 {
		return nil, errors.New("[ExHentai] error while trying to access the gallery images")
	}
	numberOfPages, err := strconv.Atoi(htmlNumberOfPages[1])
	if err != nil {
		return nil, errors.New("[ExHentai] couldn't get number of pages")
	}

	re = regexp.MustCompile("https://exhentai.org/s[^\"]+-[0-9]+")
	matchedImgURLs := re.FindAllStringSubmatch(htmlString, -1)
	imgURLs := []string{}
	for _, imgURL := range matchedImgURLs {
		imgURLs = append(imgURLs, imgURL[0])
	}

	for page := 1; len(imgURLs) < numberOfPages; page++ {
		htmlString, err := Get(fmt.Sprintf("%s?p=%d", URL, page))
		if err != nil {
			return nil, errors.New("[ExHentai] unvaild page URL")
		}
		imgURLs = append(imgURLs, re.FindStringSubmatch(htmlString)...)
	}

	data := []static.Data{}
	for idx, URL := range imgURLs {
		htmlString, err := Get(URL)
		if err != nil {
			return nil, errors.New("[ExHentai] unvaild image URL")
		}

		re := regexp.MustCompile("<h1>([^<]+)")
		matchedTitle := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedTitle) == 0 {
			return nil, errors.New("[ExHentai] unvaild image title")
		}

		re = regexp.MustCompile("<div>[^.]+([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Z]{2})")
		matchedFileInfo := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("[ExHentai] unvaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		re = regexp.MustCompile("https://exhentai.org/fullimg[^\"]+")
		srcURL := re.FindStringSubmatch(htmlString)
		if len(srcURL) != 1 {

			// sometimes the "full image url is not provided"
			re = regexp.MustCompile("<img id=\"img\" src=\"([^\"]+)")
			matchedSrcURL := re.FindAllStringSubmatch(htmlString, -1)
			if len(matchedSrcURL) != 1 {
				return nil, errors.New("[ExHentai] unvaild image src")
			}
			srcURL = []string{matchedSrcURL[0][1]}
		}

		// size will be empty if err occurs
		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		data = append(data, static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s - %d", matchedTitle[0][1], idx+1),
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						{
							URL: srcURL[0],
							Ext: fileInfo[1],
						},
					},
					Quality: fileInfo[2],
					// ex						735       KB 	== 735000Bytes
					Size: utils.CalcSizeInByte(fSize, fileInfo[4]),
				},
			},
			Url: URL,
		})

	}

	return data, nil
}*/
