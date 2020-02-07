package exhentai

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func Extract(URL string) ([]static.Data, error) {

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

	return nil, nil

}

func ParseURL(URL string) []string {
	if strings.Contains(URL, "https://exhentai.org/g/") {
		return []string{URL}
	}
	return nil
}
