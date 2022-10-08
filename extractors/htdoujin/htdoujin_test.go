package htdoujin

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
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
			Name: "Single Gallery HentaiEra 2",
			URL:  "https://hentaiera.com/gallery/610929/",
			Want: 1,
		}, {
			Name: "Tag HentaiEra",
			URL:  "https://hentaiera.com/tag/ahegao/",
			Want: 25,
		}, {
			Name: "Single Gallery HentaiEnvy",
			URL:  "https://hentaienvy.com/gallery/808735/",
			Want: 1,
		}, {
			Name: "Tag HentaiEnvy",
			URL:  "https://hentaienvy.com/parody/azur-lane/",
			Want: 28,
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
			Name: "Single Gallery HentaiRox",
			URL:  "https://hentairox.com/gallery/397913/",
			Want: 1,
		}, {
			Name: "Tag HentaiEra",
			URL:  "https://hentairox.com/tag/mosaic-censorship/",
			Want: 20,
		}, {
			Name: "Single Gallery HentaiZap",
			URL:  "https://hentaizap.com/gallery/843645/",
			Want: 1,
		}, {
			Name: "Tag HentaiZap",
			URL:  "https://hentaizap.com/tag/ahegao/",
			Want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u, err := url.Parse(tt.URL)
			test.CheckError(t, err)

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
		Args test.Args
	}{
		{
			Name: "Single Gallery HentaiEra",
			Args: test.Args{
				URL:     "https://hentaiera.com/gallery/610929/",
				Title:   "Senran Princess G",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Gallery HentaiEnvy",
			Args: test.Args{
				URL:     "https://hentaienvy.com/gallery/808737/",
				Title:   "Patreon Rewards 07-2022",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Gallery HentaiFox",
			Args: test.Args{
				URL:     "https://hentaifox.com/gallery/84580/",
				Title:   "Mirai Kairozu ni Android + 2 Plus",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Gallery HentaiRox",
			Args: test.Args{
				URL:     "https://hentairox.com/gallery/397913/",
				Title:   "Hanamizuki Vol.1",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Gallery HentaiZap",
			Args: test.Args{
				URL:     "https://hentaizap.com/gallery/843645/",
				Title:   "SUMMER FOX HUNTING",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Gallery IMHentai",
			Args: test.Args{
				URL:     "https://imhentai.xxx/gallery/684976/",
				Title:   "Otona ni Naru Hi",
				Quality: "",
				Size:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}
