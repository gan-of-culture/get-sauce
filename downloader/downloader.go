package downloader

import (
	"cmp"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/merger"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

type filePiece struct {
	offset int64
	length int64
}

type downloadInfo struct {
	URL     static.URL
	FileURI string
	Headers map[string]string
}

// downloaderStruct instance
type downloaderStruct struct {
	stream      *static.Stream
	client      *http.Client
	filePath    string
	tmpFilePath string
	filename    string
	progressBar *progressbar.ProgressBar
	bar         bool
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// New instance of Downloader
func New(bar bool) *downloaderStruct {
	return &downloaderStruct{
		client:   request.DefaultClient(),
		filePath: config.OutputPath,
		bar:      bar,
	}
}

// Download extracted data
func (downloader *downloaderStruct) Download(data *static.Data) error {
	if config.ShowInfo {
		printInfo(data)
		return nil
	}

	data.Title = cmp.Or(config.OutputName, data.Title)
	// sanitize filename here
	downloader.filename = sanitizeTitle(data.Title)

	if config.Subdirectory {
		downloader.filePath = config.OutputPath
		downloader.filePath = filepath.Join(downloader.filePath, downloader.filename)
	}

	if downloader.filePath != "" {
		err := os.MkdirAll(downloader.filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fileURIs, err := downloader.downloadStream(data)
	if err != nil {
		return err
	}

	additionalFiles, err := downloader.downloadAdditionalStreams(data)
	if err != nil {
		return err
	}

	if config.Merge == config.MergeOptNone {
		return nil
	}

	mergeFiles := make([]*merger.MergeFile, len(fileURIs))
	for i, f := range fileURIs {
		mergeFiles[i] = &merger.MergeFile{Path: f, DataType: static.DataTypeImage}
	}

	switch config.Merge {
	case config.MergeOptDefault:
		if len(additionalFiles) < 1 {
			break
		}
		downloader.stream.Ext = cmp.Or(downloader.stream.Ext, downloader.stream.URLs[0].Ext)
		mergeFiles = append(mergeFiles, additionalFiles...)
		return merger.NewDataMerger().Merge(mergeFiles, filepath.Join(downloader.filePath, fmt.Sprintf("%s_merged.%s", downloader.filename, data.Streams[config.SelectStream].Ext)))
	case config.MergeOptCBZ:
		return merger.NewArchiveMerger(downloader.bar, data).Merge(mergeFiles, filepath.Join(downloader.filePath, fmt.Sprintf("%s.cbz", downloader.filename)))
	}

	return nil
}

func (downloader *downloaderStruct) downloadStream(data *static.Data) ([]string, error) {
	// select stream to download
	var ok bool
	if downloader.stream, ok = data.Streams[config.SelectStream]; !ok {
		log.Println(data.Streams)
		return nil, fmt.Errorf("stream %s not found", config.SelectStream)
	}

	if !config.Quiet {
		printStreamInfo(data, config.SelectStream)
	}

	streamNeedsMerge := (downloader.stream.Ext != "")
	if streamNeedsMerge {
		// ensure a different tmpDir for each download so concurrent processes won't colide
		h := sha1.New()
		h.Write([]byte(data.Title + config.SelectStream))
		downloader.tmpFilePath = filepath.Join(downloader.filePath, fmt.Sprintf("%x/", h.Sum(nil)[15:]))
		err := os.MkdirAll(downloader.tmpFilePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	headers := config.FakeHeaders
	headers["Referer"] = data.URL
	maps.Copy(headers, downloader.stream.Headers)

	lenOfUrls := len(downloader.stream.URLs)
	appendEnum := false
	if lenOfUrls > 1 || config.Pages != "" {
		appendEnum = true
	}

	URLchan := make(chan downloadInfo, lenOfUrls)
	workers := min(config.Workers, lenOfUrls)
	errs, _ := errgroup.WithContext(context.TODO())

	for range workers {
		errs.Go(func() error {
			for {
				dlInfo, ok := <-URLchan
				if !ok {
					return nil
				}
				err := downloader.save(dlInfo.URL, dlInfo.FileURI, dlInfo.Headers)
				if err != nil {
					return err
				}
			}
		})
	}

	// get page numbers if -p is supplied to name files correctly
	pageNumbers := utils.NeedDownloadList(lenOfUrls)

	var fileURIs []string
	var fileURI string
	for idx, URL := range downloader.stream.URLs {
		if appendEnum {
			if config.Merge == config.MergeOptCBZ {
				fileURI = fmt.Sprint(pageNumbers[idx])
			} else {
				fileURI = fmt.Sprintf("%s_%d", downloader.filename, pageNumbers[idx])
			}
		} else {
			fileURI = downloader.filename
		}

		// build final file URI
		fileURI = filepath.Join(downloader.filePath, fileURI+"."+URL.Ext)
		if streamNeedsMerge {
			fileURI = filepath.Join(downloader.tmpFilePath, fmt.Sprintf("%d.%s", pageNumbers[idx], URL.Ext))
		}
		fileURIs = append(fileURIs, fileURI)

		URLchan <- downloadInfo{*URL, fileURI, headers}
	}
	close(URLchan)
	if err := errs.Wait(); err != nil {
		return nil, err
	}

	if streamNeedsMerge {
		// build final file URI
		fileURI = filepath.Join(downloader.filePath, downloader.filename+"."+downloader.stream.Ext)

		tmpFiles := []*merger.MergeFile{}
		for i, u := range downloader.stream.URLs {
			tmpFiles = append(tmpFiles, &merger.MergeFile{Path: filepath.Join(downloader.tmpFilePath, fmt.Sprintf("%d.%s", i, u.Ext)), DataType: downloader.stream.Type})
		}

		err := merger.NewFragmentMerger(downloader.stream.Key, downloader.bar).Merge(tmpFiles, fileURI)
		if err != nil {
			return nil, err
		}

		err = os.RemoveAll(downloader.tmpFilePath)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return []string{fileURI}, nil
	}

	return fileURIs, nil
}

func (downloader *downloaderStruct) save(URL static.URL, fileURI string, headers map[string]string) error {

	openOpts := os.O_RDWR | os.O_CREATE
	if config.Truncate {
		openOpts |= os.O_TRUNC
	}

	file, err := os.OpenFile(fileURI, openOpts, 0644)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		if !config.Quiet {
			log.Printf(`file "%s" already exists and will be skipped`, fileURI)
		}
		return nil
	}

	// if stream size bigger than 10MB then use concurWrite
	if downloader.stream.Size > 10_000_000 && config.Workers > 1 && downloader.stream.Ext == "" {
		return downloader.concurWriteFile(URL.URL, file, headers)
	}

	if err = downloader.writeFile(URL.URL, file, headers); err != nil {
		file.Close()
		os.Remove(fileURI)
		return err
	}

	file.Close()
	return nil
}

func (downloader *downloaderStruct) concurWriteFile(URL string, file *os.File, headers map[string]string) error {
	fileSize := downloader.stream.Size
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

				for k, v := range headers {
					req.Header.Set(k, v)
				}

				req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", d.offset, d.length))
				//fmt.Println(req.Header.Get("Range"))

				res, err := downloader.client.Do(req)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()

				}
				//fmt.Printf("Url: %s, Status: %s, Size: %d", URL, res.Status, res.ContentLength)
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

				buffer, err := io.ReadAll(res.Body)
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
				if downloader.bar {
					downloader.progressBar.Add(written)
				}
				lock.Unlock()

				if saveErr != nil {
					return
				}
			}
		}()
	}

	downloader.progressBar = utils.InitPB(utils.ProgressBarConfig{
		Length:      fileSize,
		Description: fmt.Sprintf("Downloading %s using %d workers...", file.Name(), config.Workers),
		AsBytes:     true,
	})

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

func (downloader *downloaderStruct) writeFile(URL string, file *os.File, headers map[string]string) error {
	// Supply http request with headers to ensure a higher possibility of success
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := downloader.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)
		res, err = downloader.client.Get(URL)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return fmt.Errorf("downloading URL: '%s' returned status %d even after retrying", URL, res.StatusCode)
		}
	}
	defer res.Body.Close()

	var writer io.Writer
	writer = file
	// some sites do not return "content-type" or "content-length" in http header
	// it will render a spinner progressbar
	downloader.progressBar = utils.InitPB(utils.ProgressBarConfig{
		Length:      res.ContentLength,
		Description: fmt.Sprintf("Downloading %s ...", file.Name()),
		AsBytes:     true,
	})
	if downloader.bar {
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

// downloadAdditionalStreams needed for e.g. a video file to be complete. Downloads audio and captions if separate and requested
func (downloader *downloaderStruct) downloadAdditionalStreams(data *static.Data) ([]*merger.MergeFile, error) {
	// everything besides of video streams doesn't need the following logic to merge using FFmpeg
	if downloader.stream.Type != static.DataTypeVideo {
		return nil, nil
	}

	var files []*merger.MergeFile
	audioFileURIs, err := downloader.downloadExtraAudio(data)
	if err != nil {
		return nil, err
	}

	captionFileURI, err := downloader.downloadCaption(data)
	if err != nil {
		return nil, err
	}

	if len(audioFileURIs) > 0 {
		for _, a := range audioFileURIs {
			files = append(files, &merger.MergeFile{Path: a, DataType: static.DataTypeAudio})
		}
	}
	if captionFileURI != "" {
		files = append(files, &merger.MergeFile{Path: captionFileURI, DataType: static.DataTypeText})
	}

	return files, nil
}

func (downloader *downloaderStruct) downloadExtraAudio(data *static.Data) ([]string, error) {
	// if audio is in separate stream -> download it with the video stream.
	// normally audio is included in the video streams. With this only special cases where this is not
	// the case are handled.

	streamID := ""
	for k, v := range data.Streams {
		if v.Type != static.DataTypeAudio {
			continue
		}
		streamID = k
	}
	if streamID == "" {
		return nil, nil
	}
	selectStreamOld := config.SelectStream
	config.SelectStream = streamID

	fileURIs, err := downloader.downloadStream(data)
	if err != nil {
		return nil, err
	}
	config.SelectStream = selectStreamOld

	return fileURIs, nil
}

func (downloader *downloaderStruct) downloadCaption(data *static.Data) (string, error) {
	if len(data.Captions) <= config.Caption || config.Caption <= -1 {
		return "", nil
	}

	headers := config.FakeHeaders
	headers["Referer"] = data.URL

	fileURI := filepath.Join(downloader.filePath, fmt.Sprintf("%s_caption_%s.%s", downloader.filename, data.Captions[config.Caption].Language, data.Captions[config.Caption].URL.Ext))
	err := downloader.save(data.Captions[config.Caption].URL, fileURI, headers)
	if err != nil {
		return "", err
	}
	if data.Captions[config.Caption].URL.Ext == "vtt" {
		err = sanitizeVTT(fileURI)
		if err != nil {
			return "", err
		}
	}

	return fileURI, nil
}
