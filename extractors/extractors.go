package extractors

import (
	"log"
	"net/url"

	"github.com/gan-of-culture/get-sauce/extractors/booru"
	"github.com/gan-of-culture/get-sauce/extractors/damn"
	"github.com/gan-of-culture/get-sauce/extractors/danbooru"
	"github.com/gan-of-culture/get-sauce/extractors/ehentai"
	"github.com/gan-of-culture/get-sauce/extractors/exhentai"
	"github.com/gan-of-culture/get-sauce/extractors/hanime"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2read"
	"github.com/gan-of-culture/get-sauce/extractors/hentai2w"
	"github.com/gan-of-culture/get-sauce/extractors/hentaicloud"
	"github.com/gan-of-culture/get-sauce/extractors/hentaidude"
	"github.com/gan-of-culture/get-sauce/extractors/hentaihaven"
	"github.com/gan-of-culture/get-sauce/extractors/hentaimama"
	"github.com/gan-of-culture/get-sauce/extractors/hentaipulse"
	"github.com/gan-of-culture/get-sauce/extractors/hentais"
	"github.com/gan-of-culture/get-sauce/extractors/hentaistream"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/hentaiyes"
	"github.com/gan-of-culture/get-sauce/extractors/hitomi"
	"github.com/gan-of-culture/get-sauce/extractors/htdoujin"
	"github.com/gan-of-culture/get-sauce/extractors/htstreaming"
	"github.com/gan-of-culture/get-sauce/extractors/iwara"
	"github.com/gan-of-culture/get-sauce/extractors/miohentai"
	"github.com/gan-of-culture/get-sauce/extractors/nhentai"
	"github.com/gan-of-culture/get-sauce/extractors/ninehentai"
	"github.com/gan-of-culture/get-sauce/extractors/ohentai"
	"github.com/gan-of-culture/get-sauce/extractors/pururin"
	"github.com/gan-of-culture/get-sauce/extractors/rule34"
	"github.com/gan-of-culture/get-sauce/extractors/simplyhentai"
	"github.com/gan-of-culture/get-sauce/extractors/thehentaiworld"
	"github.com/gan-of-culture/get-sauce/extractors/universal"
	"github.com/gan-of-culture/get-sauce/static"
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

		"9hentai.to":          ninehentaiExtractor,
		"www1.9hentai.ru":     ninehentaiExtractor,
		"booru.io":            booru.New(),
		"comicporn.xxx":       htdoujinExtractor,
		"www.damn.stream":     damnExtractor,
		"damn.stream":         damnExtractor,
		"danbooru.donmai.us":  danbooru.New(),
		"doujin.sexy":         simplyhentaiExtractor,
		"e-hentai.org":        ehentai.New(),
		"ecchi.iwara.tv":      iwara.New(),
		"exhentai.org":        exhentai.New(),
		"hanime.tv":           hanime.New(),
		"hentai2read.com":     hentai2read.New(),
		"hentai2w.com":        hentai2w.New(),
		"www.hentaicloud.com": hentaicloud.New(),
		"hentaidude.com":      hentaidude.New(),
		"hentaiera.com":       htdoujinExtractor,
		"hentaifox.com":       htdoujinExtractor,
		"hentaihaven.xxx":     hentaihaven.New(),
		"hentaimama.io":       hentaimama.New(),
		"www.hentais.tube":    hentais.New(),
		"hentaistream.moe":    hentaistream.New(),
		"hentaistream.xxx":    htstreamingExtactor,
		"hentaihaven.red":     htstreamingExtactor,
		"hentai.tv":           htstreamingExtactor,
		"animeidhentai.com":   htstreamingExtactor,
		"hentai.pro":          htstreamingExtactor,
		"hentaipulse.com":     hentaipulse.New(),
		"hentaiworld.tv":      hentaiworld.New(),
		"hentaiyes.com":       hentaiyes.New(),
		"hitomi.la":           hitomi.New(),
		"imhentai.xxx":        htdoujinExtractor,
		"miohentai.com":       miohentai.New(),
		//"muchohentai.com":       muchohentai.New(),
		"nhentai.net":           nhentai.New(),
		"ohentai.org":           ohentai.New(),
		"pururin.to":            pururin.New(),
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
