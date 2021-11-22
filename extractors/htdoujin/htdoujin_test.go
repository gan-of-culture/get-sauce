package htdoujin

import (
	"fmt"
	"net/url"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery HentaiEra",
			URL:  "https://hentaiera.com/gallery/150354/",
			Want: 1,
		}, {
			Name: "Tag HentaiEra",
			URL:  "https://hentaiera.com/tag/ahegao/",
			Want: 25,
		}, {
			Name: "Single Gallery IMHentai",
			URL:  "https://imhentai.xxx/gallery/684976/",
			Want: 1,
		}, {
			Name: "Tag IMHentai",
			URL:  "https://imhentai.xxx/artist/100yen-locker/",
			Want: 20,
		}, {
			Name: "Single Gallery HentaiFox",
			URL:  "https://hentaifox.com/gallery/84580/",
			Want: 1,
		}, {
			Name: "Tag HentaiFox",
			URL:  "https://hentaifox.com/tag/age-progression/",
			Want: 20,
		}, {
			Name: "Single Gallery HentaiEra",
			URL:  "https://hentaiera.com/gallery/150354/",
			Want: 1,
		}, {
			Name: "Tag HentaiEra",
			URL:  "https://hentaiera.com/tag/ahegao/",
			Want: 25,
		}, {
			Name: "Single Gallery HentaiRox",
			URL:  "https://hentairox.com/gallery/397913/",
			Want: 1,
		}, {
			Name: "Tag HentaiEra",
			URL:  "https://hentairox.com/tag/mosaic-censorship/",
			Want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u, err := url.Parse(tt.URL)
			if err != nil {
				t.Error(err)
			}

			if cdnPrefix, ok := sites[u.Host]; ok {
				site = "https://" + u.Host + "/"
				cdn = fmt.Sprintf("https://%s.%s/", cdnPrefix, u.Host)
			}

			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery IMHentai",
			URL:  "https://imhentai.xxx/gallery/684976/",
			Want: 1,
		}, {
			Name: "Single Gallery HentaiFox",
			URL:  "https://hentaifox.com/gallery/84580/",
			Want: 1,
		}, {
			Name: "Single Gallery HentaiEra",
			URL:  "https://hentaiera.com/gallery/488946/",
			Want: 1,
		}, {
			Name: "Single Gallery HentaiRox",
			URL:  "https://hentairox.com/gallery/397913/",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
