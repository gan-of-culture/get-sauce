package downloader

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/schollz/progressbar/v3"
)

type filePiece struct {
	offset int64
	length int64
}

type downloadInfo struct {
	URL     static.URL
	Title   string
	Headers map[string]string
}

// downloaderStruct instance
type downloaderStruct struct {
	stream      *static.Stream
	client      *http.Client
	filePath    string
	tmpFilePath string
	progressBar *progressbar.ProgressBar
	bar         bool
}

var reSanitizeTitle = regexp.MustCompile(`["&|:?<>/*\\ ]+`)

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

	if config.OutputName != "" {
		data.Title = config.OutputName
	}

	// sanitize filename here
	data.Title = strings.TrimSpace(reSanitizeTitle.ReplaceAllString(data.Title, " "))

	if config.Subdirectory {
		downloader.filePath = config.OutputPath
		downloader.filePath = filepath.Join(downloader.filePath, data.Title)
	}

	if downloader.filePath != "" {
		err := os.MkdirAll(downloader.filePath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fileURI, err := downloader.downloadStream(data)
	if err != nil {
		return err
	}

	// everything besides of video streams doesn't need the following logic to merge using FFmpeg
	if downloader.stream.Type != static.DataTypeVideo {
		return nil
	}
	var files []string
	files = append(files, fileURI)

	audioFilePath, err := downloader.downloadExtraAudio(data)
	if err != nil {
		return err
	}

	captionFilePath, err := downloader.downloadCaption(data)
	if err != nil {
		return err
	}

	if config.Keep {
		return nil
	}

	if audioFilePath != "" {
		files = append(files, audioFilePath)
	}
	if captionFilePath != "" {
		files = append(files, captionFilePath)
	}

	if downloader.stream.Ext == "" {
		downloader.stream.Ext = downloader.stream.URLs[0].Ext
	}
	return mergeMediaFiles(files, filepath.Join(downloader.filePath, data.Title+"_merged."+data.Streams[config.SelectStream].Ext))
}

func (downloader *downloaderStruct) downloadStream(data *static.Data) (string, error) {
	// select stream to download
	var ok bool
	if downloader.stream, ok = data.Streams[config.SelectStream]; !ok {
		log.Println(data.Streams)
		return "", fmt.Errorf("stream %s not found", config.SelectStream)
	}

	if !config.Quiet {
		printStreamInfo(data, config.SelectStream)
	}

	streamNeedsMerge := false
	if downloader.stream.Ext != "" {
		// ensure a different tmpDir for each download so concurrent processes won't colide
		h := sha1.New()
		h.Write([]byte(data.Title + config.SelectStream))
		downloader.tmpFilePath = filepath.Join(downloader.filePath, fmt.Sprintf("%x/", h.Sum(nil)[15:]))
		err := os.MkdirAll(downloader.tmpFilePath, os.ModePerm)
		if err != nil {
			return "", err
		}
		streamNeedsMerge = true
	}

	headers := config.FakeHeaders
	headers["Referer"] = data.URL
	maps.Copy(headers, downloader.stream.Headers)

	lenOfUrls := len(downloader.stream.URLs)
	appendEnum := false
	if lenOfUrls > 1 || config.Pages != "" {
		appendEnum = true
	}

	var saveErr error
	lock := sync.Mutex{}
	URLchan := make(chan downloadInfo, lenOfUrls)
	workers := min(config.Workers, lenOfUrls)
	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()
			for {
				dlInfo, ok := <-URLchan
				if !ok {
					return
				}
				err := downloader.save(dlInfo.URL, dlInfo.Title, dlInfo.Headers)
				if err != nil {
					lock.Lock()
					saveErr = err
					lock.Unlock()
				}
			}
		}()
	}

	// get page numbers if -p is supplied to name files correctly
	pageNumbers := utils.NeedDownloadList(lenOfUrls)

	var fileURI string
	for idx, URL := range downloader.stream.URLs {
		if appendEnum {
			fileURI = fmt.Sprintf("%s_%d", data.Title, pageNumbers[idx])
		} else {
			fileURI = data.Title
		}

		// build final file URI
		fileURI = filepath.Join(downloader.filePath, fileURI+"."+URL.Ext)
		if streamNeedsMerge {
			fileURI = filepath.Join(downloader.tmpFilePath, fmt.Sprintf("%d.%s", pageNumbers[idx], URL.Ext))
		}

		URLchan <- downloadInfo{*URL, fileURI, headers}
	}
	close(URLchan)
	wg.Wait()
	if saveErr != nil {
		return "", saveErr
	}

	if streamNeedsMerge {
		// build final file URI
		fileURI = filepath.Join(downloader.filePath, data.Title+"."+downloader.stream.Ext)
		err := downloader.MergeFilesWithSameExtension(fileURI)
		if err != nil {
			return "", err
		}
	}

	return fileURI, nil
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
	defer file.Close()

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

	return downloader.writeFile(URL.URL, file, headers)
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

	downloader.initPB(fileSize, fmt.Sprintf("Downloading %s using %d workers...", file.Name(), config.Workers), true)

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
	//fmt.Printf("Url: %s, Status: %s, Size: %d", URL, res.Status, res.ContentLength)
	if res.StatusCode != http.StatusOK {
		time.Sleep(1 * time.Second)
		res, _ = downloader.client.Get(URL)
	}
	defer res.Body.Close()

	var writer io.Writer
	writer = file
	// some sites do not return "content-type" or "content-length" in http header
	// it will render a spinner progressbar
	downloader.initPB(res.ContentLength, fmt.Sprintf("Downloading %s ...", file.Name()), true)
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

func (downloader *downloaderStruct) initPB(len int64, descr string, asBytes bool) {
	if !downloader.bar {
		return
	}
	if asBytes {
		downloader.progressBar = progressbar.NewOptions(
			int(len),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetDescription(descr),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
		return
	}
	downloader.progressBar = progressbar.NewOptions(
		int(len),
		progressbar.OptionShowIts(),
		progressbar.OptionSetDescription(descr),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
	)
}

func (downloader *downloaderStruct) MergeFilesWithSameExtension(fileURI string) error {
	lenOfUrls := len(downloader.stream.URLs)
	if lenOfUrls <= 1 {
		return nil
	}

	file, err := os.OpenFile(fileURI, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader.initPB(int64(lenOfUrls), fmt.Sprintf("Merging into %s ...", file.Name()), false)

	var d []byte
	for i, u := range downloader.stream.URLs {
		partURL := filepath.Join(downloader.tmpFilePath, fmt.Sprintf("%d.%s", i+1, u.Ext))
		if len(downloader.stream.Key) > 0 {
			d, err = decrypt(downloader.stream.Key, partURL)
			if err != nil {
				return err
			}
		} else {
			d, err = os.ReadFile(partURL)
			if err != nil {
				return err
			}
		}

		if _, err := file.Write(d); err != nil {
			return err
		}

		if downloader.bar {
			downloader.progressBar.Add(1)
		}

	}

	err = os.RemoveAll(downloader.tmpFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (downloader *downloaderStruct) downloadExtraAudio(data *static.Data) (string, error) {
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
		return "", nil
	}
	selectStreamOld := config.SelectStream
	config.SelectStream = streamID

	fileURI, err := downloader.downloadStream(data)
	if err != nil {
		return "", err
	}
	config.SelectStream = selectStreamOld

	return fileURI, nil
}

func (downloader *downloaderStruct) downloadCaption(data *static.Data) (string, error) {
	if len(data.Captions) <= config.Caption || config.Caption <= -1 {
		return "", nil
	}

	headers := config.FakeHeaders
	headers["Referer"] = data.URL

	fileURI := filepath.Join(downloader.filePath, fmt.Sprintf("%s_caption_%s.%s", data.Title, data.Captions[config.Caption].Language, data.Captions[config.Caption].URL.Ext))
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
