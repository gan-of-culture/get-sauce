package vraven

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
			Name: "Single Episode hentaihaven.xxx/",
			URL:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/episode-4/",
			Want: 1,
		}, {
			Name: "Series hentaihaven.xxx/",
			URL:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/",
			Want: 4,
		},
		{
			Name: "Single Episode hentaistream.tv",
			URL:  "https://hentaistream.tv/watch/papa-katsu/episode-1/",
			Want: 1,
		}, {
			Name: "Series hentaistream.tv",
			URL:  "https://hentaistream.tv/watch/kyonyuu-elf-oyako-saimin/",
			Want: 2,
		},
		{
			Name: "Single Episode hentaistream.io",
			URL:  "https://hentaistream.io/watch/kyonyuu-elf-oyako-saimin/episode-1/",
			Want: 1,
		}, {
			Name: "Series hentaistream.io",
			URL:  "https://hentaistream.io/watch/kyonyuu-elf-oyako-saimin/",
			Want: 2,
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
			Name: "Single Episode",
			Args: test.Args{
				URL:     "https://hentaihaven.xxx/watch/showtime-uta-no-onee-san-datte-shitai/episode-3/",
				Title:   "Showtime! Uta no Onee-san Datte Shitai - Episode 3",
				Quality: "1920x1080",
				Size:    135890160,
			},
		}, {
			Name: "Single Episode hentaistream.io",
			Args: test.Args{
				URL:     "https://hentaistream.io/watch/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation/episode-1/",
				Title:   "Onaho Kyoushitsu: Joshi Zenin Ninshin Keikaku The Animation - Episode 1",
				Quality: "1920x1080",
				Size:    275457600,
			},
		}, {
			Name: "Single Episode hentaistream.tv",
			Args: test.Args{
				URL:     "https://hentaistream.tv/watch/papa-katsu/episode-1/",
				Title:   "Papa Katsu! - Episode 1",
				Quality: "1920x1080",
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
