package extractors

import (
	"log"
	"net/url"

	"github.com/gan-of-culture/get-sauce/extractors/booru"
	"github.com/gan-of-culture/get-sauce/extractors/danbooru"
	"github.com/gan-of-culture/get-sauce/extractors/doujin"
	"github.com/gan-of-culture/get-sauce/extractors/ehentai"
	"github.com/gan-of-culture/get-sauce/extractors/exhentai"
	"github.com/gan-of-culture/get-sauce/extractors/haho"
	"github.com/gan-of-culture/get-sauce/extractors/hanime"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2read"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2w"
	"github.com/gan-of-culture/get-sauce/extractors/hentaibar"
	"github.com/gan-of-culture/get-sauce/extractors/hentaicloud"
	"github.com/gan-of-culture/get-sauce/extractors/hentaidude"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiff"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiguru"
	"github.com/gan-of-culture/get-sauce/extractors/hentaihavenred"
	"github.com/gan-of-culture/get-sauce/extractors/hentaihd"
	"github.com/gan-of-culture/get-sauce/extractors/hentaimama"
	"github.com/gan-of-culture/get-sauce/extractors/hentaimoon"
	"github.com/gan-of-culture/get-sauce/extractors/hentaipulse"
	"github.com/gan-of-culture/get-sauce/extractors/hentais"
	"github.com/gan-of-culture/get-sauce/extractors/hentaistream"
	"github.com/gan-of-culture/get-sauce/extractors/hentaivideos"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiyes"
	"github.com/gan-of-culture/get-sauce/extractors/hitomi"
	"github.com/gan-of-culture/get-sauce/extractors/htdoujin"
	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/extractors/iwara"
	"github.com/gan-of-culture/get-sauce/extractors/koushoku"
	"github.com/gan-of-culture/get-sauce/extractors/latesthentai"
	"github.com/gan-of-culture/get-sauce/extractors/manhwa18"
	"github.com/gan-of-culture/get-sauce/extractors/muchohentai"
	"github.com/gan-of-culture/get-sauce/extractors/nhentai"
	"github.com/gan-of-culture/get-sauce/extractors/nhgroup"
	"github.com/gan-of-culture/get-sauce/extractors/ninehentai"
	"github.com/gan-of-culture/get-sauce/extractors/ohentai"
	"github.com/gan-of-culture/get-sauce/extractors/pururin"
	"github.com/gan-of-culture/get-sauce/extractors/rule34"
	"github.com/gan-of-culture/get-sauce/extractors/simplyhentai"
	"github.com/gan-of-culture/get-sauce/extractors/thehentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/universal"
	"github.com/gan-of-culture/get-sauce/extractors/vraven"
	"github.com/gan-of-culture/get-sauce/extractors/zhentube"
	"github.com/gan-of-culture/get-sauce/static"
)

var extractorsMap map[string]static.Extractor

func init() {
	htstreamingExtactor := htstreaming.New()
	htdoujinExtractor := htdoujin.New()
	ninehentaiExtractor := ninehentai.New()
	vravenExtractor := vraven.New()
	nhgroupExtractor := nhgroup.New()

	extractorsMap = map[string]static.Extractor{
		"": universal.New(),

		"9hentai.to":            ninehentaiExtractor,
		"www1.9hentai.ru":       ninehentaiExtractor,
		"animeidhentai.com":     nhgroupExtractor,
		"booru.io":              booru.New(),
		"comicporn.xxx":         htdoujinExtractor,
		"danbooru.donmai.us":    danbooru.New(),
		"doujin.sexy":           doujin.New(),
		"e-hentai.org":          ehentai.New(),
		"ecchi.iwara.tv":        iwara.New(),
		"exhentai.org":          exhentai.New(),
		"haho.moe":              haho.New(),
		"hanime.io":             hanime.New(),
		"hentai.guru":           hentaiguru.New(),
		"hentai.pro":            htstreamingExtactor,
		"hentai.tv":             nhgroupExtractor,
		"hentai-moon.com":       hentaimoon.New(),
		"hentai2read.com":       hentai2read.New(),
		"hentai2w.com":          hentai2w.New(),
		"hentaibar.com":         hentaibar.New(),
		"www.hentaicloud.com":   hentaicloud.New(),
		"hentaidude.com":        hentaidude.New(),
		"hentaiera.com":         htdoujinExtractor,
		"hentaiff.com":          hentaiff.New(),
		"hentaifox.com":         htdoujinExtractor,
		"hentaihaven.com":       nhgroupExtractor,
		"hentaihaven.red":       hentaihavenred.New(),
		"hentaihaven.xxx":       vravenExtractor,
		"v2.hentaihd.net":       hentaihd.New(),
		"hentaimama.io":         hentaimama.New(),
		"hentaipulse.com":       hentaipulse.New(),
		"hentairox.com":         htdoujinExtractor,
		"www.hentais.tube":      hentais.New(),
		"hentaistream.moe":      hentaistream.New(),
		"hentaistream.io":       vravenExtractor,
		"hentaistream.tv":       vravenExtractor,
		"hentaistream.xxx":      htstreamingExtactor,
		"hentaivideos.net":      hentaivideos.New(),
		"hentaiworld.tv":        hentaiworld.New(),
		"hentaiyes.com":         hentaiyes.New(),
		"hitomi.la":             hitomi.New(),
		"imhentai.xxx":          htdoujinExtractor,
		"koushoku.org":          koushoku.New(),
		"latesthentai.com":      latesthentai.New(),
		"manhwa18.tv":           manhwa18.New(),
		"muchohentai.com":       muchohentai.New(),
		"nhentai.net":           nhentai.New(),
		"ohentai.org":           ohentai.New(),
		"pururin.to":            pururin.New(),
		"rule34.paheal.net":     rule34.New(),
		"www.simply-hentai.com": simplyhentai.New(),
		"thehentaiworld.com":    thehentaiworld.New(),
		"uncensoredhentai.xxx":  htstreamingExtactor,
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
