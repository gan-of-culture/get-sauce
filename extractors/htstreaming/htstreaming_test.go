package htstreaming

import (
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
			Name: "Single Episode hentaihaven.com",
			URL:  "https://hentaihaven.com/soshite-watashi-wa-sensei-ni-episode-1/",
			Want: 1,
		}, {
			Name: "Single Episode uncensoredhentai.xxx",
			URL:  "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			Want: 1,
		}, {
			Name: "Overview uncensoredhentai.xxx",
			URL:  "https://uncensoredhentai.xxx/genres/ahegao/",
			Want: 18,
		}, {
			Name: "Single Episode hentai.pro",
			URL:  "https://hentai.pro/knight-of-erin-episode-2/",
			Want: 1,
		}, {
			Name: "Overview hentai.pro",
			URL:  "https://hentai.pro/tag/breasts/",
			Want: 50,
		}, {
			Name: "Single Episode hentaistream.xxx",
			URL:  "https://hentaistream.xxx/watch/ijirare-fukushuu-saimin-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentaistream.xxx",
			URL:  "https://hentaistream.xxx/genres/ahegao/",
			Want: 18,
		}, /*{
			Name: "Single Episode hentai.tv",
			URL:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentai.tv",
			URL:  "https://hentai.tv/trending/",
			Want: 24,
		}, */
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
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
			Name: "Single Episode hentaihaven.com",
			Args: test.Args{
				URL:     "https://hentaihaven.com/soshite-watashi-wa-sensei-ni-episode-1/",
				Title:   "Soshite Watashi wa Sensei ni… Episode 1",
				Quality: "1920x1080",
				Size:    633196032,
			},
		},
		{
			Name: "Single Episode uncensoredhentai.xxx",
			Args: test.Args{
				URL:     "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
				Title:   "Mako chan Kaihatsu Nikki Episode 1",
				Quality: "1920x1080",
				Size:    558305856,
			},
		},
		{
			Name: "Single Episode hentai.pro",
			Args: test.Args{
				URL:     "https://hentai.pro/imaizumin-chi-wa-douyara-gal-no-tamariba-ni-natteru-rashii-episode-2/",
				Title:   "Imaizumin-chi wa Douyara Gal no Tamariba ni Natteru Rashii Episode 2",
				Quality: "1920x1080",
				Size:    288810112,
			},
		},
		{
			Name: "Single Episode hentaistream.xxx",
			Args: test.Args{
				URL:     "https://hentaistream.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
				Title:   "Mako chan Kaihatsu Nikki Episode 1",
				Quality: "1920x1080",
				Size:    558305856,
			},
		},
		{
			Name: "Single Episode hentai.tv",
			Args: test.Args{
				URL:     "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
				Title:   "Chiisana Tsubomi no Sono Oku Ni…… Episode 1",
				Quality: "1280x720",
				Size:    234157572,
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
