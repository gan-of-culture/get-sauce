package hentaihavenred

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
			URL:  "https://hentaihaven.red/hentai/ikusa-otome-suvia-episode-4/",
			Want: 1,
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
				URL:     "https://hentaihaven.red/hentai/ikusa-otome-suvia-episode-4/",
				Title:   "Ikusa Otome Suvia Episode 4",
				Quality: "1280x720",
			},
		}, {
			Name: "Single Episode 2",
			Args: test.Args{
				URL:     "https://hentaihaven.red/hentai/class-de-otoko-wa-boku-hitori-episode-1/",
				Title:   "Class de Otoko wa Boku Hitori!? Episode 1",
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
