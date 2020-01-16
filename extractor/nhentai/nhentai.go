package nhentai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

const site = "https://nhentai.net"

// Extract data from supplied url
func Extract(url string) ([]static.Data, error) {
	magicNumber, page := ParseURL(url)
	if magicNumber == "" && page == "" {
		return nil, errors.New("[NHentai]No magic number found")
	}

	if config.Pages != "" {
		page = config.Pages
	}

	htmlString, err := request.Get(fmt.Sprintf("https://nhentai.net/g/%s/", magicNumber))
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(htmlString)
	title := doc.Find("h1").Text()

	data := []static.Data{}

	// get all img links
	pages := doc.FindAll("a", "class", "gallerythumb")

	// if one page is selected get that TODO multiple pages
	if page != "" {
		pageNo, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
		pages = pages[pageNo-1 : pageNo]
	}

	for _, page := range pages {
		stream, err := extractImageData(page.Attrs()["href"])
		if err != nil {
			return nil, err
		}
		data = append(data, static.Data{
			Site:    site,
			Title:   title,
			Type:    "image",
			Streams: stream,
			Url:     url,
		})
	}

	return data, nil
}

// ParseURL data
func ParseURL(url string) (string, string) {
	re := regexp.MustCompile("[0-9]+")
	urlNumbers := re.FindAllString(url, -1)

	if len(urlNumbers) <= 0 {
		return "", ""
	}

	// if there are two "int" values it means the exact page was supplied
	var page string
	if len(urlNumbers) > 1 {
		page = urlNumbers[1]
	}

	return urlNumbers[0], page
}

func extractImageData(id string) (map[string]static.Stream, error) {

	htmlString, err := request.Get(site + id)
	if err != nil {
		return nil, err
	}
	// some times you need to retry
	if strings.Contains(htmlString, "<title>503 Service Temporarily Unavailable</title>") {
		htmlString, err = request.Get(site + id)
	}

	doc := soup.HTMLParse(htmlString)
	imgTag := doc.Find("section", "id", "image-container").Find("img").Attrs()

	stream := make(map[string]static.Stream)
	stream["0"] = static.Stream{
		Url:     imgTag["src"],
		Quality: fmt.Sprintf("%s x %s", imgTag["width"], imgTag["height"]),
		Size:    0,
	}

	return stream, nil
}
