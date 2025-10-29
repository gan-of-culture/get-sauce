package extractors

import (
	"log"
	"net/url"

	"github.com/gan-of-culture/get-sauce/extractors/booru"
	"github.com/gan-of-culture/get-sauce/extractors/danbooru"
	"github.com/gan-of-culture/get-sauce/extractors/ehentai"
	"github.com/gan-of-culture/get-sauce/extractors/haho"
	"github.com/gan-of-culture/get-sauce/extractors/hanime"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2read"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2w"
	"github.com/gan-of-culture/get-sauce/extractors/hentaicloud"
	"github.com/gan-of-culture/get-sauce/extractors/hentaifoundry"
	"github.com/gan-of-culture/get-sauce/extractors/hentaimama"
	"github.com/gan-of-culture/get-sauce/extractors/hentaimoon"
	"github.com/gan-of-culture/get-sauce/extractors/hentainexus"
	"github.com/gan-of-culture/get-sauce/extractors/hentaipulse"
	"github.com/gan-of-culture/get-sauce/extractors/hentaivideos"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiyes"
	"github.com/gan-of-culture/get-sauce/extractors/hitomi"
	"github.com/gan-of-culture/get-sauce/extractors/hstream"
	"github.com/gan-of-culture/get-sauce/extractors/htdoujin"
	"github.com/gan-of-culture/get-sauce/extractors/iwara"
	"github.com/gan-of-culture/get-sauce/extractors/miohentai"
	"github.com/gan-of-culture/get-sauce/extractors/muchohentai"
	"github.com/gan-of-culture/get-sauce/extractors/nhentai"
	"github.com/gan-of-culture/get-sauce/extractors/nhgroup"
	"github.com/gan-of-culture/get-sauce/extractors/ninehentai"
	"github.com/gan-of-culture/get-sauce/extractors/ohentai"
	"github.com/gan-of-culture/get-sauce/extractors/oppai"
	"github.com/gan-of-culture/get-sauce/extractors/rule34"
	"github.com/gan-of-culture/get-sauce/extractors/rule34video"
	"github.com/gan-of-culture/get-sauce/extractors/simplyhentai"
	"github.com/gan-of-culture/get-sauce/extractors/thehentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/universal"
	"github.com/gan-of-culture/get-sauce/static"
)

var extractorsMap map[string]static.Extractor

func init() {
	htdoujinExtractor := htdoujin.New()
	ninehentaiExtractor := ninehentai.New()
	nhgroupExtractor := nhgroup.New()
	iwaraExtractor := iwara.New()

	extractorsMap = map[string]static.Extractor{
		"": universal.New(),

		"9hentai.to":             ninehentaiExtractor,
		"www1.9hentai.ru":        ninehentaiExtractor,
		"animeidhentai.com":      nhgroupExtractor,
		"booru.io":               booru.New(),
		"comicporn.xxx":          htdoujinExtractor,
		"danbooru.donmai.us":     danbooru.New(),
		"e-hentai.org":           ehentai.New(),
		"haho.moe":               haho.New(),
		"hanime.tv":              hanime.New(),
		"hentai.tv":              nhgroupExtractor,
		"hentai-moon.com":        hentaimoon.New(),
		"hentai2read.com":        hentai2read.New(),
		"hentai2w.com":           hentai2w.New(),
		"www.hentaicloud.com":    hentaicloud.New(),
		"hentaiera.com":          htdoujinExtractor,
		"hentaienvy.com":         htdoujinExtractor,
		"www.hentai-foundry.com": hentaifoundry.New(),
		"hentaifox.com":          htdoujinExtractor,
		"hentaihaven.co":         nhgroupExtractor,
		"hentaimama.io":          hentaimama.New(),
		"hentainexus.com":        hentainexus.New(),
		"hentaipulse.com":        hentaipulse.New(),
		"hentairox.com":          htdoujinExtractor,
		"hstream.moe":            hstream.New(),
		"hentaistream.xxx":       nhgroupExtractor,
		"hentaivideos.net":       hentaivideos.New(),
		"hentaiworld.tv":         hentaiworld.New(),
		"hentaiyes.com":          hentaiyes.New(),
		"hentaizap.com":          htdoujinExtractor,
		"hitomi.la":              hitomi.New(),
		"imhentai.xxx":           htdoujinExtractor,
		"iwara.tv":               iwaraExtractor,
		"www.iwara.tv":           iwaraExtractor,
		"latesthentai.com":       nhgroupExtractor,
		"miohentai.com":          miohentai.New(),
		"muchohentai.com":        muchohentai.New(),
		"nhentai.net":            nhentai.New(),
		"ohentai.org":            ohentai.New(),
		"oppai.stream":           oppai.New(),
		"rule34.paheal.net":      rule34.New(),
		"rule34video.com":        rule34video.New(),
		"www.simply-hentai.com":  simplyhentai.New(),
		"thehentaiworld.com":     thehentaiworld.New(),
		"uncensoredhentai.xxx":   nhgroupExtractor,
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
