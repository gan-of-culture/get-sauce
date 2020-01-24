package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/downloader"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/danbooru"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/hanime"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/nhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/rule34"
	"github.com/gan-of-culture/go-hentai-scraper/extractor/underhentai"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func init() {
	flag.StringVar(&config.Pages, "p", "", "Enter pages like 1,2,3-4,6,7,8-9")
	flag.StringVar(&config.OutputName, "o", "", "Output name")
	flag.StringVar(&config.OutputPath, "O", "", "Output path")
	flag.BoolVar(&config.ShowInfo, "i", false, "Show info")
	flag.StringVar(&config.SelectStream, "s", "0", "Select a stream")
}

func download(url string) {
	var err error
	var data []static.Data
	re := regexp.MustCompile("http.://([^\\.]+)")
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 1 {
		log.Fatal("Can't parse URL")
	}
	switch matches[1] {
	case "nhentai":
		data, err = nhentai.Extract(url)
	case "rule34":
		data, err = rule34.Extract(url)
	case "danbooru":
		data, err = danbooru.Extract(url)
	case "hanime":
		data, err = hanime.Extract(url)
	case "underhentai":
		data, err = underhentai.Extract(url)
	}
	if err != nil {
		log.Fatal(err)
	}

	if config.ShowInfo {
		fmt.Println(data)
		return
	}

	for _, d := range data {
		downloader.Download(d)
	}

}

func main() {
	flag.Parse()
	args := flag.Args()
	download(args[0])
}
