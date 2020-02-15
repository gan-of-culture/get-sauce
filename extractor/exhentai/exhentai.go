package exhentai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "https://exhentai.org"

func Extract(URL string) ([]static.Data, error) {
	URLs := ParseURL(URL)
	if URLs == nil {
		return nil, errors.New("[ExHentai no URL parsed]")
	}

	//Opts
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(config.FakeHeaders["User-Agent"]),
		chromedp.WindowSize(400, 400),
	)

	//ExecAllocator
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	// Login user
	err := chromedp.Run(ctx, login())
	if err != nil {
		return nil, err
	}

	var innerHTMLBody string
	// run task list
	err = chromedp.Run(ctx,
		chromedp.Navigate(URL),
		chromedp.WaitReady("#gdt", chromedp.ByID),
		chromedp.InnerHTML("/html/body", &innerHTMLBody, chromedp.BySearch),
	)
	if err != nil || innerHTMLBody == "" {
		return nil, errors.New("[ExHentai] no matching gallery found")
	}

	re := regexp.MustCompile("([0-9]+) pages")
	htmlNumberOfPages := re.FindStringSubmatch(innerHTMLBody)
	if len(htmlNumberOfPages) != 2 {
		return nil, errors.New("[ExHentai] error while trying to access the gallery images")
	}
	numberOfPages, err := strconv.Atoi(htmlNumberOfPages[1])
	if err != nil {
		return nil, errors.New("[ExHentai] couldn't get number of pages")
	}

	re = regexp.MustCompile("https://exhentai.org/s[^\"]+-[0-9]+")
	matchedImgURLs := re.FindAllStringSubmatch(innerHTMLBody, -1)
	imgURLs := []string{}
	for _, imgURL := range matchedImgURLs {
		imgURLs = append(imgURLs, imgURL[0])
	}

	for page := 1; len(imgURLs) < numberOfPages; page++ {
		var iHTML string
		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("%s?p=%d", URL, page)),
			chromedp.InnerHTML("/html/body", &iHTML, chromedp.BySearch),
		)
		if err != nil {
			return nil, errors.New("[ExHentai] unvaild page URL")
		}
		imgURLs = append(imgURLs, re.FindStringSubmatch(iHTML)...)
	}

	data := []static.Data{}
	for idx, URL := range imgURLs {
		err := chromedp.Run(ctx,
			chromedp.Navigate(URL),
			chromedp.InnerHTML("/html/body", &innerHTMLBody, chromedp.BySearch),
		)
		if err != nil {
			return nil, errors.New("[ExHentai] unvaild image URL")
		}

		re := regexp.MustCompile("<h1>([^<]+)")
		matchedTitle := re.FindAllStringSubmatch(innerHTMLBody, -1)
		if len(matchedTitle) == 0 {
			return nil, errors.New("[ExHentai] unvaild image title")
		}

		re = regexp.MustCompile("<div>[^.]+([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Z]{2})")
		matchedFileInfo := re.FindAllStringSubmatch(innerHTMLBody, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("[ExHentai] unvaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		re = regexp.MustCompile("https://exhentai.org/fullimg[^\"]+")
		srcURL := re.FindStringSubmatch(innerHTMLBody)
		if len(srcURL) != 1 {

			// sometimes the "full image url is not provided"
			re = regexp.MustCompile("<img id=\"img\" src=\"([^\"]+)")
			matchedSrcURL := re.FindAllStringSubmatch(innerHTMLBody, -1)
			if len(matchedSrcURL) != 1 {
				return nil, errors.New("[ExHentai] unvaild image src")
			}
			srcURL = []string{matchedSrcURL[0][1]}
		}

		// size will be empty if err occurs
		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		data = append(data, static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s - %d", matchedTitle[0][1], idx+1),
			Type:  "image",
			Streams: map[string]static.Stream{
				"0": {
					URLs: []static.URL{
						{
							URL: srcURL[0],
							Ext: fileInfo[1],
						},
					},
					Quality: fileInfo[2],
					// ex						735       KB 	== 735000Bytes
					Size: utils.CalcSizeInByte(fSize, fileInfo[4]),
				},
			},
			Url: URL,
		})

	}

	return data, nil
}

func ParseURL(URL string) []string {
	if strings.Contains(URL, "https://exhentai.org/g/") {
		return []string{URL}
	}
	return nil
}

func login() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("https://forums.e-hentai.org/index.php?"),
		chromedp.Click("//*[@id=\"userlinksguest\"]/p/a[1]", chromedp.BySearch),
		chromedp.WaitReady("//*[@id=\"border\"]/tbody/tr/td/table[3]/tbody/tr/td/table/tbody/tr[2]/td/div/div/div/form/div/div/table/tbody/tr[1]/td[1]/fieldset/table/tbody/tr[2]/td[2]/input", chromedp.BySearch),
		chromedp.SendKeys("//*[@id=\"border\"]/tbody/tr/td/table[3]/tbody/tr/td/table/tbody/tr[2]/td/div/div/div/form/div/div/table/tbody/tr[1]/td[1]/fieldset/table/tbody/tr[1]/td[2]/input", config.Username),
		chromedp.SendKeys("//*[@id=\"border\"]/tbody/tr/td/table[3]/tbody/tr/td/table/tbody/tr[2]/td/div/div/div/form/div/div/table/tbody/tr[1]/td[1]/fieldset/table/tbody/tr[2]/td[2]/input", config.UserPassword),
		chromedp.Click("//*[@id=\"border\"]/tbody/tr/td/table[3]/tbody/tr/td/table/tbody/tr[2]/td/div/div/div/form/div/div/table/tbody/tr[2]/td/input", chromedp.BySearch),
		chromedp.WaitReady("//*[@id=\"userlinks\"]/p[2]/b/a", chromedp.BySearch),
	}
}
