package downloader

import (
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v2"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type filePiece struct {
	data   []byte
	offset int64
}

type Downloader struct {
	data        static.Data
	stream      string
	client      *http.Client
	filePath    string
	tmpFilePath string
	progressBar *progressbar.ProgressBar
	bar         bool
}

func New(data static.Data, stream string, bar bool) *Downloader {
	// ensure a different tmpDir for each m3u8 download so concurrent processes won't colide
	h := sha1.New()
	h.Write([]byte(data.Title))
	tmpFilePath := filepath.Join(config.OutputPath, fmt.Sprintf("%x/", h.Sum(nil)[15:]))

	return &Downloader{
		data:   data,
		stream: stream,
		client: &http.Client{
			Transport: &http.Transport{
				DisableCompression:  true,
				TLSHandshakeTimeout: 10 * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				IdleConnTimeout:     5 * time.Second,
				//DisableKeepAlives:   true,
			},
			Timeout: 15 * time.Minute,
		},
		filePath:    config.OutputPath,
		tmpFilePath: tmpFilePath,
		bar:         bar,
	}
}

func (downloader *Downloader) Download() error {

	// select stream to download
	var stream static.Stream
	var ok bool
	if stream, ok = downloader.data.Streams[config.SelectStream]; !ok {
		return fmt.Errorf("Stream %s not found", config.SelectStream)
	}

	appendEnum := false
	if len(stream.URLs) > 1 {
		appendEnum = true
	}

	var wg sync.WaitGroup
	var saveErr error
	var fileURI string
	for idx, URL := range stream.URLs {
		if appendEnum {
			fileURI = fmt.Sprintf("%s_%d", downloader.data.Title, idx+1)
		} else {
			fileURI = downloader.data.Title
		}

		if config.OutputName != "" && len(stream.URLs) == 1 {
			fileURI = config.OutputName
		}

		//sanitize filename here
		fileURI = strings.ReplaceAll(fileURI, "|", "_")
		fileURI = strings.ReplaceAll(fileURI, ":", "")

		//build final file URI
		fileURI = filepath.Join(downloader.filePath, fileURI+"."+URL.Ext)

		wg.Add(1)
		go func(URL static.URL, title string) {
			defer wg.Done()
			err := downloader.save(URL, title, downloader.data.Type)
			if err != nil {
				saveErr = err
			}
		}(URL, fileURI)
		if saveErr != nil {
			return saveErr
		}

	}
	wg.Wait()

	return nil
}

func (downloader *Downloader) save(url static.URL, fileURI string, mimeType string) error {

	file, err := os.Create(fileURI)
	if err != nil {
		return err
	}

	switch mimeType {
	case "application/x-mpegurl":
		file.Close()
		file, err = os.OpenFile(fileURI, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = downloader.writeM3U(url.URL, file)
		if err != nil {
			return err
		}
		return nil
	default:
		_, err = downloader.writeFile(url.URL, file, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (downloader *Downloader) writeFile(URL string, file *os.File, headers map[string]string) (int64, error) {
	res, err := downloader.client.Get(URL)
	if err != nil {
		return 0, err
	}
	//fmt.Printf("Url: %s, Status: %s, Size: %d", url, res.Status, res.ContentLength)
	if res.Status != "200 OK" {
		time.Sleep(1 * time.Second)
		res, err = downloader.client.Get(URL)
	}
	defer res.Body.Close()

	var writer io.Writer
	writer = file
	downloader.progressBar = nil
	if downloader.bar {
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
	written, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, fmt.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}

func (downloader *Downloader) parseSegments(URL string) ([]*m3u8.MediaSegment, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("Invalid m3u8 url %s", URL)
	}

	masterFileResp, err := downloader.client.Get(URL)
	if err != nil {
		return nil, err
	}
	defer masterFileResp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(masterFileResp.Body, true)
	if err != nil {
		return nil, err
	}

	var savedSegments []*m3u8.MediaSegment
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)

		for i, seg := range mediapl.Segments {
			if seg == nil {
				continue
			}

			if !strings.Contains(seg.URI, "http") {
				segmentURL, err := baseURL.Parse(seg.URI)
				if err != nil {
					return nil, err
				}

				seg.URI = segmentURL.String()
			}

			if seg.Key == nil && mediapl.Key != nil {
				seg.Key = mediapl.Key
			}

			seg.Title = fmt.Sprintf("%d", i)
			savedSegments = append(savedSegments, seg)
		}

	case m3u8.MASTER:
		return nil, fmt.Errorf("%s M3U File is a master! Needs to be a media list instead", p.String())
	}

	return savedSegments, nil
}

func (downloader *Downloader) writeM3U(url string, file *os.File) (int64, error) {
	segments, err := downloader.parseSegments(url)
	if err != nil {
		return 0, err
	}
	if len(segments) < 1 {
		return 0, fmt.Errorf("No segments found in %s", url)
	}

	err = os.MkdirAll(downloader.tmpFilePath, os.ModePerm)
	if err != nil {
		return 0, err
	}

	if downloader.bar {
		downloader.progressBar = progressbar.NewOptions(
			len(segments),
			progressbar.OptionSetDescription(fmt.Sprintf("Downloading segements of %s ...", file.Name())),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	var written int64
	for _, seg := range segments {

		w, err := downloader.writeSeg(seg)
		if err != nil {
			return 0, err
		}

		if downloader.bar {
			downloader.progressBar.Add(1)
		}
		written += w
	}

	err = downloader.mergeSegments(file, segments)
	if err != nil {
		return 0, err
	}

	err = os.RemoveAll(downloader.tmpFilePath)
	if err != nil {
		return 0, err
	}

	return written, nil

}

func (downloader *Downloader) writeSeg(segment *m3u8.MediaSegment) (int64, error) {
	// Supply http request with headers to ensure a higher possibility of success
	req, err := http.NewRequest(http.MethodGet, segment.URI, nil)
	if err != nil {
		return 0, err
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	if ref := req.Header.Get("Referer"); ref == "" {
		req.Header.Set("Referer", segment.URI)
	}

	res, err := downloader.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)

		res, err = downloader.client.Do(req)
		if err != nil {
			return 0, err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Println(segment.SeqId)
			fmt.Println(req.Header)
			return 0, errors.New(res.Status)
		}
	}

	file, err := os.Create(filepath.Join(downloader.tmpFilePath, segment.Title+".ts"))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	w, copyErr := io.Copy(file, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return 0, fmt.Errorf("file copy error: %s", copyErr)
	}

	return int64(w), nil
}

func (downloader *Downloader) mergeSegments(file *os.File, segments []*m3u8.MediaSegment) error {

	sort.Slice(segments, func(i, j int) bool {
		return segments[i].SeqId < segments[j].SeqId
	})

	if downloader.bar {
		downloader.progressBar = progressbar.NewOptions(
			len(segments),
			progressbar.OptionSetDescription(fmt.Sprintf("Merging into %s ...", file.Name())),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	for _, seg := range segments {

		d, err := decrypt(seg, filepath.Join(downloader.tmpFilePath, seg.Title+".ts"))
		if err != nil {
			return err
		}

		if _, err := file.Write(d); err != nil {
			return err
		}

		if downloader.bar {
			downloader.progressBar.Add(1)
		}

	}

	return nil
}
