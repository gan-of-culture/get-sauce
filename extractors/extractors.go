package extractors

import (
	"log"
	"net/url"

	"github.com/gan-of-culture/get-sauce/v2/extractors/booru"
	"github.com/gan-of-culture/get-sauce/v2/extractors/damn"
	"github.com/gan-of-culture/get-sauce/v2/extractors/danbooru"
	"github.com/gan-of-culture/get-sauce/v2/extractors/ehentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/exhentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hanime"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentai2read"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentai2w"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaibar"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaicloud"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaidude"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaiff"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaihaven"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaihavenred"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaimama"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaimoon"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaipulse"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentais"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaistream"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaiworld"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hentaiyes"
	"github.com/gan-of-culture/get-sauce/v2/extractors/hitomi"
	"github.com/gan-of-culture/get-sauce/v2/extractors/htdoujin"
	"github.com/gan-of-culture/get-sauce/v2/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/v2/extractors/iwara"
	"github.com/gan-of-culture/get-sauce/v2/extractors/miohentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/nhentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/ninehentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/ohentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/pururin"
	"github.com/gan-of-culture/get-sauce/v2/extractors/rule34"
	"github.com/gan-of-culture/get-sauce/v2/extractors/simplyhentai"
	"github.com/gan-of-culture/get-sauce/v2/extractors/thehentaiworld"
	"github.com/gan-of-culture/get-sauce/v2/extractors/universal"
	"github.com/gan-of-culture/get-sauce/v2/extractors/zhentube"
	"github.com/gan-of-culture/get-sauce/v2/static"
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
		"hanime.io":             hanime.New(),
		"hentai-moon.com":       hentaimoon.New(),
		"hentai2read.com":       hentai2read.New(),
		"hentai2w.com":          hentai2w.New(),
		"hentaibar.com":         hentaibar.New(),
		"www.hentaicloud.com":   hentaicloud.New(),
		"hentaidude.com":        hentaidude.New(),
		"hentaiera.com":         htdoujinExtractor,
		"hentaiff.com":          hentaiff.New(),
		"hentaifox.com":         htdoujinExtractor,
		"hentaihaven.xxx":       hentaihaven.New(),
		"hentaimama.io":         hentaimama.New(),
		"www.hentais.tube":      hentais.New(),
		"hentaistream.moe":      hentaistream.New(),
		"hentaihaven.com":       htstreamingExtactor,
		"hentaistream.xxx":      htstreamingExtactor,
		"hentaihaven.red":       hentaihavenred.New(),
		"hentai.tv":             htstreamingExtactor,
		"animeidhentai.com":     htstreamingExtactor,
		"hentai.pro":            htstreamingExtactor,
		"uncensoredhentai.xxx":  htstreamingExtactor,
		"hentaipulse.com":       hentaipulse.New(),
		"hentaiworld.tv":        hentaiworld.New(),
		"hentaiyes.com":         hentaiyes.New(),
		"hitomi.la":             hitomi.New(),
		"imhentai.xxx":          htdoujinExtractor,
		"miohentai.com":         miohentai.New(),
		"nhentai.net":           nhentai.New(),
		"ohentai.org":           ohentai.New(),
		"pururin.to":            pururin.New(),
		"rule34.paheal.net":     rule34.New(),
		"www.simply-hentai.com": simplyhentaiExtractor,
		"thehentaiworld.com":    thehentaiworld.New(),
		"zhentube.com":          zhentube.New(),
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
