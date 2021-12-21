package htstreaming

import (
	"strings"
	"testing"
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
		}, */{
			Name: "Single Episode animeidhentai.com",
			URL:  "https://animeidhentai.com/31678/mako-chan-kaihatsu-nikki-episode-1/",
			Want: 1,
		}, {
			Name: "Series animeidhentai.com",
			URL:  "https://animeidhentai.com/hentai/mako-chan-kaihatsu-nikki/",
			Want: 4,
		},
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
			Name: "Single Episode hentai.pro",
			URL:  "https://hentai.pro/imaizumin-chi-wa-douyara-gal-no-tamariba-ni-natteru-rashii-episode-2/",
			Want: 1,
		}, {
			Name: "Single Episode hentaistream.xxx",
			URL:  "https://hentaistream.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			Want: 1,
		}, /* {
			Name: "Single Episode hentai.tv",
			URL:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			Want: 1,
		}, */{
			Name: "Single Episode animeidhentai.com",
			URL:  "https://animeidhentai.com/31680/mako-chan-kaihatsu-nikki-episode-2/",
			Want: 1,
		}, {
			Name: "Series animeidhentai.com",
			URL:  "https://animeidhentai.com/hentai/mako-chan-kaihatsu-nikki/",
			Want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil && !strings.Contains(err.Error(), "Video not found") {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
