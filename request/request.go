package request

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

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
		LocationURL, _ := resp.Location()
		if !strings.ContainsAny(LocationURL.String(), " ") {
			return
		}
		resp.Header.Set("Location", strings.ReplaceAll(LocationURL.String(), " ", "%20"))
	}
	return
}

//Request http
func Request(method string, url string, headers map[string]string) (*http.Response, error) {

	client := &http.Client{
		Transport: LogRedirects{&http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			IdleConnTimeout:     5 * time.Second,
			//DisableKeepAlives:   true,
		}},
		Timeout: 15 * time.Minute,
	}

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

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get content as string
func Get(url string) (string, error) {
	resp, err := Request(http.MethodGet, url, nil)
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

// Headers return the HTTP Headers of the url
func Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
		"Range":   "bytes=0-0",
	}
	res, err := Request(http.MethodGet, url, headers)
	if err != nil {
		return nil, err
	}
	return res.Header, nil
}

// Size get size of the url
func Size(url, refer string) (int64, error) {
	if !config.ShowInfo {
		return 0, nil
	}

	resp, err := Request(http.MethodGet, url, map[string]string{
		"Range": "bytes=0-0",
	})
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == 503 {
		time.Sleep(200 * time.Millisecond)
		resp, err = Request(http.MethodGet, url, map[string]string{
			"Range": "bytes=0-0",
		})
	}

	s := resp.Header.Get("Content-Range")
	if s == "" {
		fmt.Println(url)
		return 0, errors.New("content-range is not present")
	}
	size, err := strconv.ParseInt(strings.Split(s, "/")[1], 10, 64)
	if err != nil {
		return 0, err
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
		for k_old, cookieInJar := range p.Jar[u.Host] {
			if cookie.Name == cookieInJar.Name && !isInJar {
				isInJar = true
				p.Jar[u.Host][k_old] = cookies[k]
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
