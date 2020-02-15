package hanime

/*
import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/request"
	"github.com/gan-of-culture/go-hentai-scraper/static"
	"github.com/gan-of-culture/go-hentai-scraper/utils"
)

const site = "hanime.tv"

// Extractor for img/video data
func Extract(URL string) ([]static.Data, error) {

	URLs, err := ParseURL(URL)
	if err != nil {
		return nil, err
	}

	data := []static.Data{}
	for _, URL := range URLs {

		re := regexp.MustCompile("([0-9]{4,})\\.")
		ids := re.FindStringSubmatch(URL)
		title := ids[1]

		size, err := request.Size(URL, URL)
		if err != nil {
			return nil, err
		}

		ext := utils.GetLastItemString(strings.Split(URL, "."))
		ext = strings.Split(ext, "?")[0]
		contType := "image"
		if strings.HasPrefix(ext, "gif") {
			contType = "gif"
			ext = "gif"
		}

		streams := make(map[string]static.Stream)
		streams["0"] = static.Stream{
			URLs: []static.URL{
				{
					URL: URL,
					Ext: ext,
				},
			},
			Quality: fmt.Sprintf("%s x %s", "unknown", "unknown"),
			Size:    size,
		}

		data = append(data, static.Data{
			Site:    site,
			Title:   title,
			Type:    contType,
			Streams: streams,
			Url:     URL,
		})
	}

	return data, nil
}

// ParseURL of input url
func ParseURL(URL string) ([]string, error) {
	if strings.Contains(URL, "/uploads/") {
		return []string{URL}, nil
	}

	if !strings.Contains(URL, "browse/images") {
		return nil, errors.New("[HAnime] invalid URL")
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
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var nodes []*cdp.Node
	// run task list
	err := chromedp.Run(ctx,
		chromedp.Navigate(URL),
		chromedp.Nodes(".cuc.grows", &nodes, chromedp.NodeEnabled, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	URLs := []string{}
	for _, node := range nodes[0].Parent.Children {
		if node.NodeName == "A" {
			URLs = append(URLs, node.Attributes[1])
		}
	}

	return URLs, nil
}
*/
