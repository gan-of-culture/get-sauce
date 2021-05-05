package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/schollz/progressbar/v2"
)

// Download data
func Download(data static.Data) error {

	var wg sync.WaitGroup

	if config.SelectStream == "" {
		config.SelectStream = "0"
	}

	// select stream to download
	var stream static.Stream
	var ok bool
	if stream, ok = data.Streams[config.SelectStream]; !ok {
		return fmt.Errorf("Stream %s not found", config.SelectStream)
	}

	appendEnum := false
	if len(stream.URLs) > 1 {
		appendEnum = true
	}

	var saveErr error
	var URLTitle string
	for idx, URL := range stream.URLs {
		if appendEnum {
			URLTitle = fmt.Sprintf("%s_%d", data.Title, idx)
		} else {
			URLTitle = data.Title
		}

		//sanitize filename here
		URLTitle = strings.ReplaceAll(URLTitle, "|", "_")

		wg.Add(1)
		go func(URL static.URL, title string) {
			defer wg.Done()
			err := save(URL, title, config.FakeHeaders)
			if err != nil {
				saveErr = err
			}
		}(URL, URLTitle)
		if saveErr != nil {
			return saveErr
		}
	}
	wg.Wait()

	return nil
}

func save(url static.URL, fileName string, headers map[string]string) error {
	if config.OutputName != "" {
		fileName = config.OutputName
	}

	var filePath string
	if config.OutputPath != "" {
		filePath = config.OutputPath
	}

	file, err := os.Create(filePath + fileName + "." + url.Ext)
	if err != nil {
		return err
	}

	_, err = writeFile(url.URL, file, headers)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(url string, file *os.File, headers map[string]string) (int64, error) {
	res, err := request.Request(http.MethodGet, url, headers)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("Url: %s, Status: %s, Size: %d", url, res.Status, res.ContentLength)
	if res.Status != "200 OK" {
		time.Sleep(1 * time.Second)
		res, err = request.Request(http.MethodGet, url, headers)
	}
	defer res.Body.Close()

	bar := progressbar.NewOptions(
		int(res.ContentLength),
		progressbar.OptionSetBytes(int(res.ContentLength)),
		progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s ...", file.Name())),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
		//progressbar.OptionShowCount(),
	)
	writer := io.MultiWriter(file, bar)

	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, fmt.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}
