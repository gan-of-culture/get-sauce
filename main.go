package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/downloader"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/booru"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/danbooru"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/ehentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/exhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/hentaimama"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/hentais"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/hentaistream"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/hentaiworld"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/imgboard"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/nhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/rule34"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/universal"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func init() {
	flag.IntVar(&config.Amount, "a", 0, "Amount of files to download")
	flag.StringVar(&config.OutputName, "o", "", "Output name")
	flag.StringVar(&config.OutputPath, "O", "", "Output path (include ending slash)")
	flag.StringVar(&config.Pages, "p", "", "Enter pages like 1,2,3-4,6,7,8-9 for doujins")
	flag.BoolVar(&config.RestrictContent, "r", false, "Don't scrape Restricted Content")
	flag.StringVar(&config.SelectStream, "s", "0", "Select a stream")
	flag.BoolVar(&config.ShowInfo, "i", false, "Show info")
	flag.IntVar(&config.Threads, "t", 1, "Number of threads used for downloading")
	flag.StringVar(&config.Username, "un", "", "Username for exhentai/forum e hentai")
	flag.StringVar(&config.UserPassword, "up", "", "User password for exhentai/forum e hentai")
}

func download(URL string) {
	var err error
	var data []static.Data
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Identified site: %s", u.Host)

	switch u.Host {
	case "booru.io":
		data, err = booru.Extract(URL)
	case "danbooru.donmai.us":
		data, err = danbooru.Extract(URL)
	case "e-hentai.org":
		data, err = ehentai.Extract(URL)
	case "exhentai.org":
		data, err = exhentai.Extract(URL)
	case "hentaimama.io":
		data, err = hentaimama.Extract(URL)
	case "hentais.tube":
		data, err = hentais.Extract(URL)
	case "hentaistream.moe":
		data, err = hentaistream.Extract(URL)
	case "hentaiworld.tv":
		data, err = hentaiworld.Extract(URL)
	case "nhentai.net":
		data, err = nhentai.Extract(URL)
	case "rule34.paheal.net":
		data, err = rule34.Extract(URL)
	case "rule34.xxx":
		data, err = imgboard.Extract(URL)
	default:
		data, err = imgboard.Extract(URL)
		if err != nil {
			data, err = universal.Extract(URL, u.Host)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	if config.ShowInfo {
		for _, singleData := range data {
			jsonData, _ := json.MarshalIndent(singleData, "", "    ")
			fmt.Printf("%s\n", jsonData)
		}
		return
	}

	var wg sync.WaitGroup
	wg.Add(config.Threads)
	datachan := make(chan static.Data, len(data))

	for i := 0; i < config.Threads; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case d, ok := <-datachan:
					if !ok {
						return
					}
					err := downloader.Download(d)
					if err != nil {
						log.Println(err)
					}
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
	download(args[0])
}
