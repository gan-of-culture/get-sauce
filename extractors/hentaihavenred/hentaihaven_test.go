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
			Name: "Single Episode hentaihaven.red/",
			URL:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentaihaven.red/",
			URL:  "https://hentaihaven.red/ratings/",
			Want: 35,
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
				URL:     "https://hentaihaven.red/hentai/bitch-na-inane-sama-episode-4/",
				Title:   "Bitch na Inane-sama Episode 4",
				Quality: "1920x1080",
				Size:    0,
			},
		},
		{
			Name: "[OLD] Single Episode",
			Args: test.Args{
				URL:     "https://hentaihaven.red/hentai/mako-chan-kaihatsu-nikki-episode-2/",
				Title:   "Mako-chan Kaihatsu Nikki Episode 2",
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
