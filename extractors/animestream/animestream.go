package animestream

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/utils"
)

// module to parseURL of the wordpress theme animestream
// https://themesia.com/animestream-wordpress-theme/

// ParseURL of the wordpress theme
func ParseURL(URL, site string) []string {

	reEpisodeURL := regexp.MustCompile(site + `(?:\d+|watch)/.+/`)
	reParseURLShow := regexp.MustCompile(site + `(?:hentai|anime)/[\w-%]+/`)

	if ok := reEpisodeURL.MatchString(URL); ok {
		return []string{URL}
	}

	htmlString, err := request.Get(URL)
	if err != nil {
		return nil
	}

	if reParseURLShow.MatchString(URL) {
		htmlString = strings.Split(htmlString, `<div class="bixbox"`)[0]
		return utils.RemoveAdjDuplicates(reEpisodeURL.FindAllString(htmlString, -1))
	}

	// contains list of show that need to be derefenced to episode level
	htmlString = strings.Split(htmlString, `<div id="sidebar">`)[0]

	out := []string{}
	for _, anime := range reParseURLShow.FindAllString(htmlString, -1) {
		out = append(out, ParseURL(anime, site)...)
	}
	return out
}

// ParseURL of the wordpress theme
func ParseURLwoSite(URL string) []string {
	u, err := url.Parse(URL)
	if err != nil {
		return nil
	}

	site := "https://" + u.Host + "/"

	return ParseURL(URL, site)
}
