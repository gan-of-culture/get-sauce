package danbooru

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

const site = "https://danbooru.donmai.us"

// Extractor for danbooru pages
func Extractor(url string) ([]static.Data, error) {
	posts, err := ParseURL(url)
	if err != nil {
		return nil, err
	}

	data := []static.Data{}
	for _, post := range posts {
		contentData, err := extractData(site + post)
		if err != nil {
			return nil, err
		}
		data = append(data, contentData)
	}

	return data, nil
}

// ParseURL for danbooru pages
func ParseURL(url string) ([]string, error) {
	re := regexp.MustCompile("page=([0-9]+)")
	pageNo := re.FindAllString(url, -1)
	// pageNo = url?page=number -> if it's there it means overview page otherwise single post or invalid
	if len(pageNo) == 0 {

		re := regexp.MustCompile("[/]posts[/]([0-9]+)")
		linkToPost := re.FindString(url)
		if linkToPost == "" {
			return nil, errors.New("[Danbooru]Invalid Url no post found")
		}

		out := []string{}
		out = append(out, linkToPost)
		return out, nil
	}

	htmlString, err := request.Get(url)
	if err != nil {
		return nil, err
	}

	doc := soup.HTMLParse(htmlString)
	container := doc.Find("div", "id", "posts-container")
	if container.Error != nil {
		return nil, errors.New("[Danbooru] " + container.Error.Error())
	}

	items := container.FindAll("article")
	if len(items) == 0 {
		return nil, errors.New("[Danbooru]No articles found in overview page")
	}

	out := []string{}
	for _, item := range items {
		out = append(out, "/posts/"+item.Attrs()["data-id"])
	}

	return out, nil
}

func extractData(postURL string) (static.Data, error) {
	htmlString, err := request.Get(postURL)
	if err != nil {
		return static.Data{}, err
	}

	doc := soup.HTMLParse(htmlString)
	imageContainer := doc.Find("section", "id", "image-container")
	if imageContainer.Error != nil {
		return static.Data{}, errors.New("[Danbooru] " + imageContainer.Error.Error())
	}

	attrs := imageContainer.Attrs()
	size, err := request.Size(attrs["data-large-file-url"], postURL)
	if err != nil {
		return static.Data{}, errors.New("[Danbooru]No image size not found")
	}

	streams := make(map[string]static.Stream, 1)
	streams["0"] = static.Stream{
		URLs: []URL{
			{
				URL: attrs["data-large-file-url"],
				Ext: utils.GetLastItem(strings.Split(attrs["data-large-file-url"], ".")),
			},
		},
		Quality: fmt.Sprintf("%s x %s", attrs["data-width"], attrs["data-height"]),
		Size:    size,
	}

	title := getTitle(doc.Find("section", "id", "tag-list"))
	if title == "" {
		title = attrs["data-id"]
	}

	return static.Data{
		Site:  site,
		Title: title,
		Type:  "image",

		Streams: streams,
		Url:     postURL,
	}, nil
}

func getTitle(tagList soup.Root) string {
	if tagList == (soup.Root{}) {
		return ""
	}

	title := ""
	tagLists := []string{"copyright", "character", "artist"}
	for _, listName := range tagLists {
		list := tagList.Find("ul", "class", fmt.Sprintf("%s-tag-list", listName))
		if list.Error != nil {
			continue
		}

		tagText := list.Children()[0].Find("a", "class", "search-tag")
		if tagText.Error != nil {
			continue
		}

		title = title + " " + tagText.Text()
	}

	title = strings.ReplaceAll(title, "(", " ")
	title = strings.ReplaceAll(title, ")", " ")
	title = strings.ReplaceAll(title, "|", " ")

	return strings.ReplaceAll(title, "/", " ")
}
