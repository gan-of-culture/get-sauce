package hanime

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
			Name: "Single Episode",
			URL:  "https://hanime.tv/videos/hentai/inkou-kyoushi-no-saimin-seikatsu-shidouroku-2",
			Want: 1,
		}, {
			Name: "Overview",
			URL:  "https://hanime.tv/browse/tags/fantasy",
			Want: 24,
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
				URL:     "https://hanime.tv/videos/hentai/inkou-kyoushi-no-saimin-seikatsu-shidouroku-2",
				Title:   "Inkou Kyoushi no Saimin Seikatsu Shidouroku 2",
				Quality: "720p; 1280 x 720",
				Size:    115000000,
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
