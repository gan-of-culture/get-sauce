package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
)

func progressBar(size int64) *pb.ProgressBar {
	bar := pb.New64(size).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(1000)
	return bar
}

// Download data
func Download(data static.Data) error {
	
	var waitgroup sync.WaitGroup 

	bar := progressBar(data.Size)
	bar.Start()

	var save_err error

	for _, URL := range data.Streams[0].URLs {
		wg.add()
		go func( URL URL, title data.Title, bar *pb.ProgressBar ) {
			defer wgp.Done()
			err := save(URL, title, config.FakeHeaders, bar)
			if err != nil {
				save_err = err
			}
		}
		if save_err != nil{
			return save_err
		}
	}
	bar.Finish()
	return nil
}

func save(url static.URL, fileName string, headers map[string]string, bar pb.ProgressBar) error {
	if config.OuputName != "" {
		fileName = config.OutputName
	}

	var filePath string
	if config.OuputPath != "" {
		filePath = config.OutputPath
	}

	file, err := os.Create(filePath + fileName + "." + url.Ext)
	if err != nil {
		return err
	}

	written, err := writeFile(url.URL, file, bar)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(url string, file *os.File, bar *pb.ProgressBar) (int64, error) {
	res, err := request.Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	writer := io.MultiWriter(file, bar)
	// Note that io.Copy reads 32kb(maximum) from input and writes them to output, then repeats.
	// So don't worry about memory.
	written, copyErr := io.Copy(writer, res.Body)
	if copyErr != nil && copyErr != io.EOF {
		return written, fmt.Errorf("file copy error: %s", copyErr)
	}
	return written, nil
}
