package nhgroup

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
			Name: "Single Episode animeidhentai.com",
			URL:  "https://animeidhentai.com/31678/mako-chan-kaihatsu-nikki-episode-1-subbed/",
			Want: 1,
		}, {
			Name: "Series animeidhentai.com",
			URL:  "https://animeidhentai.com/hentai/mako-chan-kaihatsu-nikki/",
			Want: 4,
		}, {
			Name: "Single Episode hentaihaven.co/",
			URL:  "https://hentaihaven.co/watch/seika-jogakuin-koutoubu-kounin-sao-oji-san-episode-3/",
			Want: 1,
		}, {
			Name: "Overview hentaihaven.co/",
			URL:  "https://hentaihaven.co/brand/bunnywalker/",
			Want: 36,
		}, {
			Name: "Single Episode hentaihaven.red/",
			URL:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentaihaven.red/",
			URL:  "https://hentaihaven.red/brand/bunnywalker/",
			Want: 36,
		}, {
			Name: "Single Episode hentai.tv",
			URL:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentai.tv",
			URL:  "https://hentai.tv/trending/",
			Want: 48,
		}, {
			Name: "Single Episode hentaistream.xxx",
			URL:  "https://hentaistream.xxx/watch/ijirare-fukushuu-saimin-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentaistream.xxx",
			URL:  "https://hentaistream.xxx/genres/ahegao/",
			Want: 36,
		}, {
			Name: "Single Episode Eng Sub",
			URL:  "https://latesthentai.com/watch/hajimete-no-hitozumaepisode-1-sub",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://latesthentai.com/serie/hajimete-no-hitozuma/",
			Want: 5,
		}, {
			// this is the same logic for all extensions that group shows e.g. /genres/
			// its hard to make a test for the other groups since the number of episodes always changes
			Name: "Studio",
			URL:  "https://latesthentai.com/genre/ahegao/",
			Want: 36,
		}, {
			Name: "Single Episode uncensoredhentai.xxx",
			URL:  "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			Want: 1,
		}, {
			Name: "Overview uncensoredhentai.xxx",
			URL:  "https://uncensoredhentai.xxx/genres/ahegao/",
			Want: 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := ParseURL(tt.URL)
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
			Name: "Single Episode animeidhentai.com",
			Args: test.Args{
				URL:   "https://animeidhentai.com/36364/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3/",
				Title: "Usamimi Bouken-tan: Sekuhara Shinagara Sekai o Sukue Episode 3",
				Size:  288687974,
			},
		},
		{
			Name: "Single Episode hentaihaven.co",
			Args: test.Args{
				URL:   "https://hentaihaven.co/watch/seika-jogakuin-koutoubu-kounin-sao-oji-san-episode-3/",
				Title: "Seika Jogakuin Koutoubu Kounin Sao Oji-San Episode 3",
				Size:  391364277,
			},
		},
		{
			Name: "Single Episode hentai.tv",
			Args: test.Args{
				URL:   "https://hentai.tv/hentai/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3/",
				Title: "Usamimi Bouken-tan: Sekuhara Shinagara Sekai o Sukue Episode 3",
				Size:  288687974,
			},
		},
		{
			Name: "Single Episode hentaistream.xxx",
			Args: test.Args{
				URL:   "https://hentaistream.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
				Title: "Mako-chan Kaihatsu Nikki Episode 1",
				Size:  398919647,
			},
		},
		{
			Name: "Single Episode",
			Args: test.Args{
				URL:   "https://latesthentai.com/watch/hajimete-no-hitozuma-episode-1/",
				Title: "Hajimete no Hitozuma Episode 1",
				Size:  278979392,
			},
		},
		{
			Name: "Single Episode uncensoredhentai.xxx",
			Args: test.Args{
				URL:   "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
				Title: "Mako-chan Kaihatsu Nikki Episode 1",
				Size:  398919647,
			},
		},
		{
			Name: "[New] Single Episode uncensoredhentai.xxx",
			Args: test.Args{
				URL:   "https://uncensoredhentai.xxx/watch/joshi-luck-episode-5/",
				Title: "Joshi Luck! Episode 5",
				Size:  423727631,
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
