package request

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

//Request http
func Request(method string, url string, headers map[string]string) (*http.Response, error) {

	transport := &http.Transport{
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
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
		return "", err
	}

	return string(body), nil
}

// Headers return the HTTP Headers of the url
func Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}
	res, err := Request(http.MethodGet, url, headers)
	if err != nil {
		return nil, err
	}
	return res.Header, nil
}

// Size get size of the url
func Size(url, refer string) (int64, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return 0, err
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
