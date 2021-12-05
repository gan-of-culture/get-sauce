package request

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

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

// LogRedirects to sanitize "Location" URLs
type LogRedirects struct {
	Transport http.RoundTripper
}

//RoundTrip implementation
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
		LocationURL, _ := resp.Location()
		if !strings.ContainsAny(LocationURL.String(), " ") {
			return
		}
		resp.Header.Set("Location", strings.ReplaceAll(LocationURL.String(), " ", "%20"))
	}
	return
}

//DefaultClient to use in the scraper
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: LogRedirects{&http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:       true,
				PreferServerCipherSuites: false,
				CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
			},
			IdleConnTimeout: 5 * time.Second,
			//DisableKeepAlives:   true,
		}},
		Timeout: 10 * time.Minute,
	}
}

//Request http
func Request(method string, URL string, headers map[string]string, body io.Reader) (*http.Response, error) {

	client := DefaultClient()

	req, err := http.NewRequest(method, URL, body)
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
		req.Header.Set("Referer", URL)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get content as string
func Get(URL string) (string, error) {
	body, err := GetAsBytes(URL)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetAsBytes content as string
func GetAsBytes(URL string) ([]byte, error) {
	resp, err := Request(http.MethodGet, URL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// PostAsBytes content as string
func PostAsBytes(URL string) ([]byte, error) {
	resp, err := Request(http.MethodPost, URL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// GetWithHeaders content as string
func GetWithHeaders(URL string, headers map[string]string) (string, error) {
	body, err := GetAsBytesWithHeaders(URL, headers)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetAsBytesWithHeaders content as string
func GetAsBytesWithHeaders(URL string, headers map[string]string) ([]byte, error) {
	resp, err := Request(http.MethodGet, URL, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// PostAsBytesWithHeaders content as string
func PostAsBytesWithHeaders(URL string, headers map[string]string) ([]byte, error) {
	resp, err := Request(http.MethodPost, URL, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// Headers return the HTTP Headers of the URL
func Headers(URL, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}
	res, err := Request(http.MethodHead, URL, headers, nil)
	if err == nil {
		return res.Header, nil
	}
	if res != nil && res.StatusCode == 503 {
		time.Sleep(200 * time.Millisecond)
		res, err := Request(http.MethodHead, URL, headers, nil)
		if err == nil {
			return res.Header, nil
		}
	}

	headers["Range"] = "bytes=0-1"
	res, err = Request(http.MethodGet, URL, headers, nil)
	if err != nil {
		return nil, err
	}
	return res.Header, nil
}

// Size get size of the URL
func Size(URL, refer string) (int64, error) {
	// if you are trying to scrape more than one thing
	// sending size request just make it slower thinking of image boards etc.
	if config.Amount != 0 {
		return 0, nil
	}

	headers, err := Headers(URL, refer)
	if err != nil {
		return 0, err
	}

	size, err := GetSizeFromHeaders(&headers)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// GetSizeFromHeaders of http.Response
func GetSizeFromHeaders(headers *http.Header) (int64, error) {
	s := utils.GetLastItemString(strings.Split(headers.Get("Content-Range"), "/"))
	if s == "" {
		s = headers.Get("Content-Length")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if size == 0 {
		return 0, errors.New("Size not found")
	}
	return size, nil
}

// Myjar of client
type Myjar struct {
	Jar map[string][]*http.Cookie
}

// SetCookies of client
func (p *Myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// swap cookie assignment after login
	if u.Host == "forums.e-hentai.org" {
		u.Host = "exhentai.org"
	}
	//fmt.Printf("The URL is : %s\n", u.String())
	//fmt.Printf("The cookie being set is : %s\n", cookies)

	//preserve old cookies and overwrite old ones with new cookies
	isInJar := false
	for k, cookie := range cookies {
		isInJar = false
		for keyOld, cookieInJar := range p.Jar[u.Host] {
			if cookie.Name == cookieInJar.Name && !isInJar {
				isInJar = true
				p.Jar[u.Host][keyOld] = cookies[k]
			}
		}
		if !isInJar {
			p.Jar[u.Host] = append(p.Jar[u.Host], cookies[k])
		}
	}
	//p.Jar[u.Host] = cookies
}

// Cookies of client
func (p *Myjar) Cookies(u *url.URL) []*http.Cookie {
	//fmt.Printf("The URL is : %s\n", u.String())
	//fmt.Printf("Cookie being returned is : %s\n", p.Jar[u.Host])
	return p.Jar[u.Host]
}

// GetM3UMeta segment URLs
func GetM3UMeta(master *string, URL string) ([]*static.URL, []byte, error) {
	re := regexp.MustCompile(`\s[^#]+\s`) // 1=segment URI
	matchedSegmentURLs := re.FindAllString(*master, -1)
	if len(matchedSegmentURLs) == 0 {
		fmt.Println(*master)
		return nil, nil, errors.New("no segements found")
	}

	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, nil, err
	}

	segments := []*static.URL{}
	segmentURI := ""
	for _, v := range matchedSegmentURLs {
		segmentURI = strings.TrimSpace(v)
		if !strings.Contains(segmentURI, "http") {
			segmentURL, err := baseURL.Parse(segmentURI)
			if err != nil {
				return nil, nil, err
			}
			segmentURI = segmentURL.String()
		}
		segments = append(segments, &static.URL{
			URL: segmentURI,
			Ext: utils.GetFileExt(segmentURI),
		})
	}

	re = regexp.MustCompile(`#EXT-X-KEY:METHOD=([^,]*),URI="([^"]*)`) //1=HASH e.g. AES-128 2=KEYURI
	matchedEncryptionMeta := re.FindStringSubmatch(*master)
	if len(matchedEncryptionMeta) != 3 {
		return segments, nil, nil
	}

	keyURL := matchedEncryptionMeta[2]
	if !strings.HasPrefix(matchedEncryptionMeta[2], "http") {
		keyURI, err := baseURL.Parse(matchedEncryptionMeta[2])
		if err != nil {
			return nil, nil, err
		}
		keyURL = keyURI.String()
	}

	res, err := Request(http.MethodGet, keyURL, map[string]string{
		"Referer": URL,
	}, nil)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return segments, buffer, nil
}
