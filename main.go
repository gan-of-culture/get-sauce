package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/downloader"
	"github.com/gan-of-culture/get-sauce/extractors"
	"github.com/gan-of-culture/get-sauce/static"
)

// name of the application
const name = "get-sauce"

// version is overridden by ldflags in release
var version = "v0.0.0"

func init() {
	flag.IntVar(&config.Amount, "a", 0, "Amount of files to download")
	flag.IntVar(&config.Caption, "c", -1, "Download caption to a extra file")
	flag.StringVar(&config.UserHeaders, "h", "", "UserHeaders for the HTTP requests. To bypass Cloudflare or DDOS-GUARD protection")
	flag.BoolVar(&config.ShowInfo, "i", false, "Show info")
	flag.BoolVar(&config.ShowExtractedData, "j", false, "Show extracted data as json")
	flag.BoolVar(&config.Keep, "k", false, "Keep video, audio and subtitles. Don't merge using ffmpeg")
	flag.StringVar(&config.OutputName, "o", "", "Output name")
	flag.StringVar(&config.OutputPath, "O", "", "Output path (include ending delimiter)")
	flag.StringVar(&config.Pages, "p", "", "Enter pages like 1,2,3-4,6,7,8-9 for doujins")
	flag.BoolVar(&config.Quiet, "q", false, "Quiet mode - show minimal information")
	flag.StringVar(&config.SelectStream, "s", "0", "Select a stream")
	flag.BoolVar(&config.Subdirectory, "S", false, "Subdirectory for the downloaded content. The directory name defaults to a cleaned up version of the data title")
	flag.BoolVar(&config.Truncate, "t", false, "Truncate file if it already exists")
	flag.IntVar(&config.Timeout, "T", 10, "Timeout for the http.client in minutes")
	flag.BoolVar(&config.Version, "v", false, "Print the current version")
	flag.IntVar(&config.Workers, "w", 1, "Number of workers used for downloading")
}

func download(URL string) {
	if !config.Keep {
		_, err := exec.LookPath("ffmpeg")
		if err != nil {
			log.Println("No merging possible, because ffmpeg is not installed or not found in PATH")
			config.Keep = true
		}
	}

	data, err := extractors.Extract(URL)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if config.ShowExtractedData {
		for _, singleData := range data {
			jsonData, _ := json.MarshalIndent(singleData, "", "    ")
			fmt.Printf("%s\n", jsonData)
		}
		return
	}

	if config.SelectStream == "" {
		config.SelectStream = "0"
	}

	lenOfData := len(data)
	/*
		We have 3 main types of data that has to be downloaded concurrently
		1. lenOfData = 3000 e.g. mass scraping image boards
		2. lenOfData = 1 URLs = 200 e.g. doujin
		3. lenOfData = 1-10 but big file size e.g.hentai video
		here in main we will deal with the first type
	*/
	workers := min(config.Workers, lenOfData)

	var wg sync.WaitGroup
	wg.Add(workers)
	datachan := make(chan *static.Data, lenOfData)

	downloader := downloader.New(true)
	for range workers {
		go func() {
			defer wg.Done()
			for {
				d, ok := <-datachan
				if !ok {
					return
				}
				err := downloader.Download(d)
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}

	for _, d := range data {
		datachan <- d
	}
	close(datachan)
	wg.Wait()
}

func main() {
	flag.Parse()
	args := flag.Args()
	if config.Version {
		fmt.Printf("\n%s: version %s\n\n", name, version)
	}

	if len(args) < 1 {
		fmt.Println("Too few arguments")
		fmt.Println("Usage: get-sauce [args] URLs...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, a := range args {
		download(a)
	}
}
