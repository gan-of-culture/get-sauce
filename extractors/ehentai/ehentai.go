package ehentai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

const site = "https://e-hentai.org/"

var reNumbOfPages = regexp.MustCompile(`([0-9]+) pages`)
var reIMGURLs = regexp.MustCompile(`https://e-hentai.org/s[^"]+-[0-9]+`)
var reFileInfo = regexp.MustCompile(`<div>[^.]+\.([^::]+):: ([^::]+) :: ([^.]+.[0-9]+) ([A-Za-z]{2,3})`)
var reSourceURL = regexp.MustCompile(`<img id="img" src="([^"]+)`)

type extractor struct{}

// New returns a e-hentai extractor
func New() static.Extractor {
	return &extractor{}
}

// Extract data
func (e *extractor) Extract(URL string) ([]*static.Data, error) {
	URLs := parseURL(URL)
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		rData, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, rData...)
	}
	return data, nil
}

func parseURL(URL string) []string {
	if strings.Contains(URL, "https://e-hentai.org/g/") {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return []string{}
	}

	re := regexp.MustCompile(`https://e-hentai.org/g/[^"]+`)
	galleries := re.FindAllStringSubmatch(htmlString, -1)
	if len(galleries) == 0 {
		return []string{}
	}

	out := []string{}

	for _, gallery := range galleries {
		out = append(out, gallery[0])
	}
	return out
}

func extractData(URL string) ([]*static.Data, error) {
	if !strings.Contains(URL, "?nw=session") {
		URL = URL + "?nw=session"
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	htmlNumberOfPages := reNumbOfPages.FindStringSubmatch(htmlString)
	if len(htmlNumberOfPages) != 2 {
		return nil, errors.New("error while trying to access the gallery images")
	}
	numberOfPages, err := strconv.Atoi(htmlNumberOfPages[1])
	if err != nil {
		return nil, errors.New("couldn't get number of pages")
	}

	imgURLs := reIMGURLs.FindAllString(htmlString, -1)

	// if gallery has more than 40 images -> walk other pages for links aswell
	for page := 1; len(imgURLs) < numberOfPages; page++ {
		htmlString, err := request.Get(fmt.Sprintf("%s?p=%d", URL, page))
		if err != nil {
			return nil, err
		}
		imgURLs = append(imgURLs, reIMGURLs.FindAllString(htmlString, -1)...)
	}

	data := []*static.Data{}
	for _, idx := range utils.NeedDownloadList(len(imgURLs)) {
		htmlString, err := request.Get(imgURLs[idx-1])
		if err != nil {
			return nil, err
		}

		title := utils.GetH1(&htmlString, 0)
		if title == "" {
			return nil, errors.New("invaild image title")
		}

		matchedFileInfo := reFileInfo.FindAllStringSubmatch(htmlString, -1)
		if len(matchedFileInfo) == 0 {
			return nil, errors.New("invaild image file info")
		}
		fileInfo := matchedFileInfo[0]

		// sometimes the "full image URL is not provided"
		matchedSrcURL := reSourceURL.FindAllStringSubmatch(htmlString, -1)
		if len(matchedSrcURL) != 1 {
			return nil, static.ErrDataSourceParseFailed
		}

		fSize, _ := strconv.ParseFloat(fileInfo[3], 64)

		data = append(data, &static.Data{
			Site:  site,
			Title: fmt.Sprintf("%s - %d", title, idx),
			Type:  static.DataTypeImage,
			Streams: map[string]*static.Stream{
				"0": {
					Type: static.DataTypeImage,
					URLs: []*static.URL{
						{
							URL: matchedSrcURL[0][1],
							Ext: fileInfo[1],
						},
					},
					Quality: fileInfo[2],
					Size:    utils.CalcSizeInByte(fSize, fileInfo[4]),
				},
			},
			URL: imgURLs[idx-1],
		})

	}

	return data, nil
}
