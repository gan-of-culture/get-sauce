package request

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/klauspost/compress/zstd"
	"github.com/pkg/errors"
)

// LogRedirects to sanitize "Location" URLs
type LogRedirects struct {
	Transport http.RoundTripper
}

// RoundTrip implementation
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

// DefaultClient for HTTP requests
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: LogRedirects{&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521, tls.X25519},
			},
			IdleConnTimeout: 5 * time.Second,
			//DisableKeepAlives: true,
		}},
		Timeout: time.Duration(config.Timeout) * time.Minute,
	}
}

// Firefox117Client tries to impersonate firefox117 https://github.com/lwthiker/curl-impersonate/blob/main/firefox/curl_ff117
func Firefox117Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_AES_128_GCM_SHA256, tls.TLS_CHACHA20_POLY1305_SHA256, tls.TLS_AES_256_GCM_SHA384, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA, tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_RSA_WITH_AES_128_CBC_SHA, tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			},
			DisableCompression:    false,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

// Request HTTP
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

	if config.UserHeaders != "" {
		values := strings.Split(config.UserHeaders, "\n")
		for _, val := range values {
			keyValue := strings.Split(val, ":")
			if len(keyValue) < 2 {
				continue
			}
			req.Header.Set(keyValue[0], strings.Join(keyValue[1:], ":"))
		}
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

// GetAsBytes content as bytes
func GetAsBytes(URL string) ([]byte, error) {
	resp, err := Request(http.MethodGet, URL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decompressedBody, err := DecompressHttpResponse(resp.Body, resp.Header.Get("Content-Encoding"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer decompressedBody.Close()

	body, err := io.ReadAll(decompressedBody)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}
	return body, nil
}

// GetAsBytesWithClient content as bytes using a custom client
func GetAsBytesWithClient(client *http.Client, URL string, referer string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for key, val := range config.FakeHeadersFirefox117 {
		req.Header.Set(key, val)
	}

	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	decompressedBody, err := DecompressHttpResponse(res.Body, res.Header.Get("Content-Encoding"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer decompressedBody.Close()

	body, err := io.ReadAll(decompressedBody)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// PostAsBytes content as bytes
func PostAsBytes(URL string) ([]byte, error) {
	resp, err := Request(http.MethodPost, URL, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

// GetAsBytesWithHeaders content as bytes
func GetAsBytesWithHeaders(URL string, headers map[string]string) ([]byte, error) {
	resp, err := Request(http.MethodGet, URL, headers, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decompressedBody, err := DecompressHttpResponse(resp.Body, resp.Header.Get("Content-Encoding"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer decompressedBody.Close()

	body, err := io.ReadAll(decompressedBody)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil
}

// PostAsBytesWithHeaders content as bytes
func PostAsBytesWithHeaders(URL string, headers map[string]string, body io.Reader) ([]byte, error) {
	resp, err := Request(http.MethodPost, URL, headers, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return resBody, nil
}

// GetWithCookies content as string
func GetWithCookies(URL string, jar *Myjar) (string, error) {
	client := DefaultClient()

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return "", errors.New("Request can't be created")
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	req.Header.Set("Referer", URL)

	for _, cookie := range jar.Cookies(req.URL) {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return "", err
		}
	}

	jar.SetCookies(req.URL, resp.Cookies())

	return string(body), nil
}

// Headers of a HTTP response
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

// Size of a HTTP content in response
func Size(URL, refer string) (int64, error) {
	// if you are trying to scrape more than one thing
	// sending size request just make it slower
	if config.Amount != 0 {
		return 0, nil
	}

	headers, err := Headers(URL, refer)
	if err != nil {
		return 0, nil
	}

	size, err := GetSizeFromHeaders(&headers)
	if err == nil {
		return size, nil
	}

	res, err := Request(http.MethodGet, URL, map[string]string{
		"Referer": refer,
		"Range":   "bytes=0-1",
	}, nil)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	size, err = GetSizeFromHeaders(&res.Header)
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

// DecompressHttpResponse to read it's contents
func DecompressHttpResponse(body io.ReadCloser, contentEncoding string) (io.ReadCloser, error) {
	var reader io.ReadCloser
	var err error
	switch contentEncoding {
	case "br":
		reader = io.NopCloser(brotli.NewReader(body))
	case "gzip":
		reader, err = gzip.NewReader(body)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	case "deflate":
		reader = flate.NewReader(body)
	case "zstd":
		//TODO: replace with impl of standard lib when change is landed (https://github.com/golang/go/issues/62513)
		d, err := zstd.NewReader(body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		reader = io.NopCloser(d)
	default:
		reader = body
	}

	return reader, nil
}

// Myjar of client
type Myjar struct {
	Jar map[string][]*http.Cookie
}

func (p *Myjar) New() {
	p.Jar = make(map[string][]*http.Cookie)
}

// SetCookies of client
func (p *Myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
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
