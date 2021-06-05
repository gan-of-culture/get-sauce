package downloader

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
	"github.com/schollz/progressbar/v2"
)

/*
TODO:
1. Implement concurrency for m3u segement merging
*/

type filePiece struct {
	offset int64
	length int64
}

type downloadInfo struct {
	URL   static.URL
	Title string
}

type Downloader struct {
	data        *static.Data
	stream      string
	client      *http.Client
	filePath    string
	tmpFilePath string
	progressBar *progressbar.ProgressBar
	bar         bool
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func New(stream string, bar bool) *Downloader {
	return &Downloader{
		stream:   stream,
		client:   request.DefaultClient(),
		filePath: config.OutputPath,
		bar:      bar,
	}
}

func (downloader *Downloader) Download(data *static.Data) error {
	downloader.data = data

	// select stream to download
	var stream static.Stream
	var ok bool
	if stream, ok = downloader.data.Streams[downloader.stream]; !ok {
		return fmt.Errorf("Stream %s not found", downloader.stream)
	}

	if downloader.data.Type == "application/x-mpegurl" {
		// ensure a different tmpDir for each download so concurrent processes won't colide
		h := sha1.New()
		h.Write([]byte(downloader.data.Title))
		downloader.tmpFilePath = filepath.Join(downloader.filePath, fmt.Sprintf("%x/", h.Sum(nil)[15:]))
	}

	lenOfUrls := len(stream.URLs)

	appendEnum := false
	if lenOfUrls > 1 {
		appendEnum = true
	}

	var saveErr error
	lock := sync.Mutex{}
	URLchan := make(chan downloadInfo, lenOfUrls)

	workers := config.Workers
	if config.Workers > lenOfUrls {
		workers = lenOfUrls
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for {
				dlInfo, ok := <-URLchan
				if !ok {
					return
				}
				err := downloader.save(dlInfo.URL, dlInfo.Title)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}
			}
		}()
	}

	//get page numbers if -p is supplied to name files correctly
	pageNumbers := utils.NeedDownloadList(lenOfUrls)

	var fileURI string
	for idx, URL := range stream.URLs {
		if appendEnum {
			fileURI = fmt.Sprintf("%s_%d", downloader.data.Title, pageNumbers[idx])
		} else {
			fileURI = downloader.data.Title
		}

		if config.OutputName != "" && lenOfUrls == 1 {
			fileURI = config.OutputName
		}

		//sanitize filename here
		re := regexp.MustCompile(`["&|:?<>/*\\ ]+`)
		fileURI = strings.TrimSpace(re.ReplaceAllString(fileURI, " "))

		//build final file URI
		fileURI = filepath.Join(downloader.filePath, fileURI+"."+URL.Ext)

		URLchan <- downloadInfo{URL, fileURI}
	}
	close(URLchan)
	wg.Wait()
	if saveErr != nil {
		return saveErr
	}

	return nil
}

func (downloader *Downloader) save(url static.URL, fileURI string) error {

	file, err := os.Create(fileURI)
	if err != nil {
		return err
	}

	if downloader.data.Type == "application/x-mpegurl" {
		file.Close()
		file, err = os.OpenFile(fileURI, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		err = downloader.writeM3U(url.URL, file)
		if err != nil {
			return err
		}
		return nil
	}
	defer file.Close()

	//if stream size bigger than 10MB then use concurWrite
	if downloader.data.Streams[downloader.stream].Size > 10000000 && config.Workers > 1 {

		err = downloader.concurWriteFile(url.URL, file)
		if err != nil {
			return err
		}
		return nil
	}

	err = downloader.writeFile(url.URL, file)
	if err != nil {
		return err
	}

	return nil
}

func (downloader *Downloader) concurWriteFile(URL string, file *os.File) error {
	fileSize := downloader.data.Streams[downloader.stream].Size
	pieceSize := int64(10_000_000)

	var saveErr error
	lock := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(config.Workers)
	datachan := make(chan filePiece, config.Workers)

	for i := 0; i < config.Workers; i++ {
		go func() {
			defer wg.Done()
			for {
				d, ok := <-datachan
				if !ok {
					return
				}

				req, err := http.NewRequest(http.MethodGet, URL, nil)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}

				for k, v := range config.FakeHeaders {
					req.Header.Set(k, v)
				}

				if ref := req.Header.Get("Referer"); ref == "" {
					req.Header.Set("Referer", URL)
				}
				req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", d.offset, d.length))
				//fmt.Println(req.Header.Get("Range"))

				res, err := downloader.client.Do(req)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}
				//fmt.Printf("Url: %s, Status: %s, Size: %d", url, res.Status, res.ContentLength)
				if res.StatusCode != http.StatusPartialContent {
					time.Sleep(1 * time.Second)
					res, err = downloader.client.Get(URL)
					if err != nil {
						lock.Lock()
						saveErr = err
						lock.Unlock()
					}
				}
				defer res.Body.Close()
				//fmt.Println(res.ContentLength)

				buffer, err := ioutil.ReadAll(res.Body)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}

				lock.Lock()
				written, err := file.WriteAt(buffer, d.offset)
				if err != nil {
					saveErr = err
				}
				downloader.progressBar.Add(written)
				lock.Unlock()

				if saveErr != nil {
					return
				}
			}
		}()
	}

	if downloader.bar {
		downloader.progressBar = progressbar.NewOptions(
			int(fileSize),
			progressbar.OptionSetDescription(fmt.Sprintf("Downloading with workers %s ...", file.Name())),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	var offset int64
	for ; fileSize > 0; fileSize -= pieceSize {
		if pieceSize+pieceSize > fileSize {
			pieceSize += fileSize - pieceSize
			datachan <- filePiece{offset: offset, length: offset + pieceSize}
			break
		}
		datachan <- filePiece{offset: offset, length: offset + pieceSize - 1}
		offset += pieceSize
	}
	close(datachan)
	wg.Wait()

	return nil
}

func (downloader *Downloader) writeFile(URL string, file *os.File) error {
	// Supply http request with headers to ensure a higher possibility of success
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	if ref := req.Header.Get("Referer"); ref == "" {
		req.Header.Set("Referer", URL)
	}

	res, err := downloader.client.Do(req)
	if err != nil {
		return err
	}
	//fmt.Printf("Url: %s, Status: %s, Size: %d", url, res.Status, res.ContentLength)
	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)
		res, _ = downloader.client.Get(URL)
	}
	defer res.Body.Close()

	var writer io.Writer
	writer = file
	downloader.progressBar = nil
	if downloader.bar {
		//some sites do not return "content-type" in http header
		//it will render a blank progressbar
		downloader.progressBar = progressbar.NewOptions(
			int(res.ContentLength),
			progressbar.OptionSetBytes(int(res.ContentLength)),
			progressbar.OptionSetDescription(fmt.Sprintf("Downloading %s ...", file.Name())),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
		writer = io.MultiWriter(file, downloader.progressBar)
	}

	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	_, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return fmt.Errorf("file copy error: %s", copyErr)
	}
	return nil
}
