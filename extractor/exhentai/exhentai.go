package exhentai

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://exhentai.org/"
const loginFormURL = "https://forums.e-hentai.org/index.php?act=Login&CODE=01"

type extractor struct {
	client *http.Client
}

// Login your user
func (ex *extractor) Login() error {

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br",
		"Cache-Control":             "max-age=0",
		"Connection":                "keep-alive",
		"Content-Type":              "application/x-www-form-urlencoded",
		"Host":                      "forums.e-hentai.org",
		"Origin":                    "https://forums.e-hentai.org",
		"Referer":                   "https://forums.e-hentai.org/index.php?act=Login&CODE=00?",
		"Upgrade-Insecure-Requests": "1",
	}

	//data := fmt.Sprintf("{ 'referer': 'https://forums.e-hentai.org/index.php?act=Login&CODE=00', 'b': '', 'bt': '', 'UserName': '%s', 'PassWord': '%s', 'CookieDate': '1'}", config.Username, config.UserPassword)
	//data := "referer=https%3A%2F%2Fforums.e-hentai.org%2Findex.php%3Fact%3DLogin%26CODE%3D00&b=&bt=&UserName=config.UserName&PassWord=config.UserPassword&CookieDate=1"
	params := url.Values{}
	params.Add("CookieDate", "1")
	params.Add("PassWord", config.UserPassword)
	params.Add("UserName", config.Username)
	params.Add("bt", "")
	params.Add("b", "")
	params.Add("referer", "https://forums.e-hentai.org/index.php?act=Login")

	req, err := http.NewRequest(http.MethodPost, loginFormURL, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	_, err = ex.client.Do(req)
	if err != nil {
		return err
	}

	u, _ := url.Parse(site)
	for _, cookie := range ex.client.Jar.Cookies(u) {
		if cookie.Name == "ipb_member_id" {
			return nil
		}
	}

	return fmt.Errorf("[Exhentai]No login possible for User: %s and Password: %s", config.Username, config.UserPassword)

}

//Request http
func (ex *extractor) Request(method string, url string, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New("Request can't be created")
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Referer"]; !ok {
		req.Header.Set("Referer", url)
	}

	resp, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

// Get content as string
func (ex *extractor) Get(url string) (string, error) {
	resp, err := ex.Request(http.MethodGet, url, nil)
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
func (ex *extractor) ParseURL(URL string) []string {
	//typical URL
	if ok, _ := regexp.MatchString("https://exhentai.org/[gs]/", URL); ok {
		return []string{URL}
	}

	htmlString, err := ex.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://exhentai.org/g/[^\"\s]+`)
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

func (ex *extractor) extractData(URLs []string) ([]static.Data, error) {

	data := []static.Data{}
	for idx, URL := range URLs {
		htmlString, err := ex.Get(URL)
		if err != nil {
			return nil, err
		}

		re := regexp.MustCompile("<h1>([^<]+)")
		matchedTitle := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedTitle) == 0 {
			return nil, errors.New("[ExHentai] invaild image title")
		}

		re = regexp.MustCompile(`<div>[^.]+\.([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Z]{2})`)
		matchedFileInfo := re.FindAllStringSubmatch(htmlString, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("[ExHentai] invaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		re = regexp.MustCompile("https://exhentai.org/fullimg[^\"]+")
		srcURL := re.FindStringSubmatch(htmlString)
		if len(srcURL) != 1 {

			// sometimes the "full image url is not provided"
			re = regexp.MustCompile("<img id=\"img\" src=\"([^\"]+)")
			matchedSrcURL := re.FindAllStringSubmatch(htmlString, -1)
			if len(matchedSrcURL) != 1 {
				return nil, errors.New("[ExHentai] invaild image src")
			}
			srcURL = []string{matchedSrcURL[0][1]}
		}

		// size will be empty if err occurs
		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		//get direct image full size download link by resolving the redirect
		//this http request will stop at the redirect and send back the location it's getting redirected to
		//so we don't receive the image data in this step -> it's a lot faster
		//check the New() function to see how the redirect is intercepted
		resp, err := ex.client.Get(strings.ReplaceAll(srcURL[0], "&amp;", "&"))
		if err != nil {
			switch resp.StatusCode {
			case http.StatusOK, http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
				if u, _ := resp.Location(); u.String() == "" {
					return nil, errors.New("[Exhentai]Error 509 - Bandwidth Exceeded. Check https://ehwiki.org/wiki/Technical_Issues#509")
				}
				l, _ := resp.Location()
				srcURL[0] = l.String()
			default:
				return nil, err
			}
		}

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
}

// New instance of extractor
func New() extractor {
	jar := &request.Myjar{}
	jar.Jar = make(map[string][]*http.Cookie)

	ex := extractor{client: &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 15 * time.Minute,
		Jar:     jar,
	}}
	ex.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if !strings.HasSuffix(req.URL.Host, ".hath.network") {
			return nil
		}
		return errors.New("Redirect")
	}

	return ex
}

// Extract data
func Extract(URL string) ([]static.Data, error) {
	ex := New()

	err := ex.Login()
	if err != nil {
		return nil, err
	}

	URLs := ex.ParseURL(URL)
	if len(URLs) == 0 {
		return nil, errors.New("[ExHentai] no vaild URL found")
	}

	//unpack galleries
	imgURLs := []string{}
	for _, URL := range URLs {
		htmlString, err := ex.Get(URL)
		if err != nil {
			return nil, errors.New("[ExHentai] invaild URL")
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
		pages := utils.NeedDownloadList(numberOfPages)

		re = regexp.MustCompile(`[^"]*/s/[^"\s]*`)
		matchedImgURLs := re.FindAllString(htmlString, -1)

		// with this only necessary pages of gallery are scraped
		// for example you have a gallery with 150 sites, but you only
		// want -p "1-10" there is no need to scrape the other sites
		numberOfPages = pages[len(pages)-1]
		// if gallery has more than 40 images -> walk other pages for links aswell
		for page := 1; len(matchedImgURLs) < numberOfPages; page++ {
			htmlString, err := ex.Get(fmt.Sprintf("%s?p=%d", URL, page))
			if err != nil {
				return nil, errors.New("[ExHentai] invaild page URL")
			}
			matchedImgURLs = append(matchedImgURLs, re.FindAllString(htmlString, -1)...)
		}

		for _, page := range pages {
			imgURLs = append(imgURLs, matchedImgURLs[page-1])
		}
	}

	data, err := ex.extractData(imgURLs)
	if err != nil {
		return nil, err
	}

	return data, nil
}
