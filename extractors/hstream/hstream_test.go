package hstream

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
			URL:  "https://hstream.moe/hentai/maki-chan-to-now/1",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hstream.moe/hentai/maki-chan-to-now",
			Want: 4,
		}, {
			Name: "Single Episode 4k",
			URL:  "https://hstream.moe/hentai/aku-no-onna-kanbu-full-moon-night-r/1",
			Want: 1,
		}, {
			Name: "Series 4k",
			URL:  "https://hstream.moe/hentai/aku-no-onna-kanbu-full-moon-night-r",
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
			Name: "Single Episode 4k",
			Args: test.Args{
				URL:     "https://hstream.moe/hentai/wizard-girl-ambitious/1",
				Title:   "Wizard Girl Ambitious – 1",
				Quality: "av1.2160p.webm",
				Size:    804660690,
			},
		},
		{
			Name: "Single Episode",
			Args: test.Args{
				URL:     "https://hstream.moe/hentai/maki-chan-to-now/1",
				Title:   "Maki-chan to Now. – 1",
				Quality: "av1.1080p.webm",
				Size:    180211896,
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
