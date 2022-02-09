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
		},
		{
			Name: "Single Episode hentaihaven.com",
			URL:  "https://hentaihaven.com/soshite-watashi-wa-sensei-ni-episode-1/",
			Want: 1,
		},
		{
			Name: "Single Episode hentai.tv",
			URL:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentai.tv",
			URL:  "https://hentai.tv/trending/",
			Want: 48,
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
		Args test.Args
	}{
		{
			Name: "Single Episode animeidhentai.com using htstreaming directly",
			Args: test.Args{
				URL:     "https://animeidhentai.com/31821/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-2-subbed/",
				Title:   "Usamimi Bouken-tan: Sekuhara Shinagara Sekai o Sukue Episode 2",
				Quality: "1920x1080",
				Size:    0,
			},
		},
		{
			Name: "Single Episode animeidhentai.com",
			Args: test.Args{
				URL:     "https://animeidhentai.com/36364/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3/",
				Title:   "Usamimi Bouken-Tan Sekuhara Shinagara Sekai O Sukue Episode 3",
				Quality: "1920x1080",
				Size:    0,
			},
		}, {
			Name: "Single Episode hentaihaven.com",
			Args: test.Args{
				URL:     "https://hentaihaven.com/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3/",
				Title:   "Usamimi Bouken-Tan Sekuhara Shinagara Sekai O Sukue Episode 3",
				Quality: "1920x1080",
				Size:    0,
			},
		},
		{
			Name: "Single Episode hentai.tv",
			Args: test.Args{
				URL:     "https://hentai.tv/hentai/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3/",
				Title:   "Usamimi Bouken-Tan Sekuhara Shinagara Sekai O Sukue Episode 3",
				Quality: "1920x1080",
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