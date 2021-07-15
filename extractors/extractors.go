package extractors

import (
	"log"
	"net/url"

	"github.com/gan-of-culture/go-hentai-scraper/extractors/booru"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/damn"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/danbooru"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/ehentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/exhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hanime"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentai2read"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentai2w"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaicloud"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaidude"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaihaven"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaimama"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentais"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaistream"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaiworld"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hentaiyes"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/hitomi"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/htdoujin"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/htstreaming"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/iwara"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/miohentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/muchohentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/nhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/ninehentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/ohentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/pururin"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/rule34"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/simplyhentai"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/thehentaiworld"
	"github.com/gan-of-culture/go-hentai-scraper/extractors/universal"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

var extractorsMap map[string]static.Extractor

func init() {
	damnExtractor := damn.New()
	htstreamingExtactor := htstreaming.New()
	htdoujinExtractor := htdoujin.New()
	ninehentaiExtractor := ninehentai.New()
	simplyhentaiExtractor := simplyhentai.New()

	extractorsMap = map[string]static.Extractor{
		"": universal.New(),

		"9hentai.to":            ninehentaiExtractor,
		"www1.9hentai.ru":       ninehentaiExtractor,
		"booru.io":              booru.New(),
		"comicporn.xxx":         htdoujinExtractor,
		"www.damn.stream":       damnExtractor,
		"damn.stream":           damnExtractor,
		"danbooru.donmai.us":    danbooru.New(),
		"doujin.sexy":           simplyhentaiExtractor,
		"e-hentai.org":          ehentai.New(),
		"ecchi.iwara.tv":        iwara.New(),
		"exhentai.org":          exhentai.New(),
		"hanime.tv":             hanime.New(),
		"hentai2read.com":       hentai2read.New(),
		"hentai2w.com":          hentai2w.New(),
		"www.hentaicloud.com":   hentaicloud.New(),
		"hentaidude.com":        hentaidude.New(),
		"hentaiera.com":         htdoujinExtractor,
		"hentaifox.com":         htdoujinExtractor,
		"hentaihaven.xxx":       hentaihaven.New(),
		"hentaimama.io":         hentaimama.New(),
		"www.hentais.tube":      hentais.New(),
		"hentaistream.moe":      hentaistream.New(),
		"hentaistream.xxx":      htstreamingExtactor,
		"hentaihaven.red":       htstreamingExtactor,
		"hentai.tv":             htstreamingExtactor,
		"animeidhentai.com":     htstreamingExtactor,
		"hentai.pro":            htstreamingExtactor,
		"hentaiworld.tv":        hentaiworld.New(),
		"hentaiyes.com":         hentaiyes.New(),
		"hitomi.la":             hitomi.New(),
		"imhentai.xxx":          htdoujinExtractor,
		"miohentai.com":         miohentai.New(),
		"muchohentai.com":       muchohentai.New(),
		"nhentai.net":           nhentai.New(),
		"ohentai.org":           ohentai.New(),
		"pururin.io":            pururin.New(),
		"rule34.paheal.net":     rule34.New(),
		"www.simply-hentai.com": simplyhentaiExtractor,
		"thehentaiworld.com":    thehentaiworld.New(),
		//"www.tsumino.com":       tsumino.New(),
	}
}

// Extract call the other extractors
func Extract(URL string) ([]*static.Data, error) {
	u, err := url.Parse(URL)
	if err != nil {
		log.Fatal(err)
	}

	extractor := extractorsMap[u.Host]
	if extractor == nil {
		extractor = extractorsMap[""]
	}
	return extractor.Extract(URL)
}
