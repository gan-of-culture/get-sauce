package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/gan-of-culture/go-hentai-scraper/request"
)

func progressBar(size int64) *pb.ProgressBar {
	bar := pb.New64(size).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.ShowFinalTime = true
	bar.SetMaxWidth(1000)
	return bar
}

func Download(data static.Data) error {
	//TODO add waitgroup
	bar := progressBar(data.Size)
	bar.Start()
	for _, URL := range data.Streams[0].URLs {
		err := save(URL, data.Title, config.FakeHeaders, bar)
		if err != nil {
			return
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
