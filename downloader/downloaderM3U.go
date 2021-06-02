package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v2"
)

func (downloader *Downloader) parseSegments(URL string) ([]*m3u8.MediaSegment, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("Invalid m3u8 url %s", URL)
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, errors.New("Request can't be created")
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	//this is important other wise you might get something weird as response
	if _, ok := req.Header["Referer"]; !ok {
		req.Header.Set("Referer", URL)
	}

	mediaFileResp, err := downloader.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer mediaFileResp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(mediaFileResp.Body, true)
	if err != nil {
		return nil, err
	}

	var mediapl *m3u8.MediaPlaylist
	switch listType {
	case m3u8.MEDIA:
		mediapl = p.(*m3u8.MediaPlaylist)

	case m3u8.MASTER:
		return nil, fmt.Errorf("%s M3U File is a master! Needs to be a media list instead", p.String())
	}

	savedSegments := []*m3u8.MediaSegment{}
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
		if seg.Key != nil && !strings.Contains(seg.Key.URI, "http") {
			keyURL, err := baseURL.Parse(seg.Key.URI)
			if err != nil {
				return nil, err
			}

			seg.Key.URI = keyURL.String()
		}

		seg.Title = fmt.Sprintf("%d", i)
		savedSegments = append(savedSegments, seg)
	}

	return savedSegments, nil
}

func (downloader *Downloader) writeM3U(url string, file *os.File) error {
	segments, err := downloader.parseSegments(url)
	if err != nil {
		return err
	}
	if len(segments) < 1 {
		return fmt.Errorf("No segments found in %s", url)
	}

	err = os.MkdirAll(downloader.tmpFilePath, os.ModePerm)
	if err != nil {
		return err
	}

	if downloader.bar {
		downloader.progressBar = progressbar.NewOptions(
			len(segments),
			progressbar.OptionSetDescription(fmt.Sprintf("Downloading segements of %s ...", file.Name())),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	lenOfSegements := len(segments)
	var saveErr error
	lock := sync.Mutex{}
	segmentchan := make(chan *m3u8.MediaSegment, lenOfSegements)

	workers := config.Workers
	if config.Workers > lenOfSegements {
		workers = lenOfSegements
	}

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for {
				seg, ok := <-segmentchan
				if !ok {
					return
				}
				err := downloader.writeSeg(seg)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}

				if downloader.bar {
					downloader.progressBar.Add(1)
				}
			}
		}()
	}

	for _, seg := range segments {
		segmentchan <- seg
	}
	close(segmentchan)
	wg.Wait()
	if saveErr != nil {
		return saveErr
	}

	err = downloader.mergeSegments(file, segments)
	if err != nil {
		return err
	}

	err = os.RemoveAll(downloader.tmpFilePath)
	if err != nil {
		return err
	}

	return nil

}

func (downloader *Downloader) writeSeg(segment *m3u8.MediaSegment) error {
	// Supply http request with headers to ensure a higher possibility of success
	req, err := http.NewRequest(http.MethodGet, segment.URI, nil)
	if err != nil {
		return err
	}

	for k, v := range config.FakeHeaders {
		req.Header.Set(k, v)
	}

	if ref := req.Header.Get("Referer"); ref == "" {
		req.Header.Set("Referer", segment.URI)
	}

	res, err := downloader.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)

		res, err = downloader.client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			fmt.Println(segment.SeqId)
			fmt.Println(req.Header)
			return errors.New(res.Status)
		}
	}

	file, err := os.Create(filepath.Join(downloader.tmpFilePath, segment.Title+".ts"))
	if err != nil {
		return err
	}
	defer file.Close()

	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	_, copyErr := io.Copy(file, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return fmt.Errorf("file copy error: %s", copyErr)
	}

	return nil
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
