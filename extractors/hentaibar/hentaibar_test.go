package hentaibar

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
			URL:  "https://hentaibar.com/videos/2688/soshite-watashi-wa-ojisan-ni-episode-4-english-subbed/",
			Want: 1,
		}, {
			Name: "Group",
			URL:  "https://hentaibar.com/tags/1080p/",
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
				URL:     "https://hentaibar.com/videos/2670/torokase-orgasm-the-animation-episode-1-english-subbed/",
				Title:   "Torokase Orgasm The Animation Episode 1 English Subbed",
				Quality: "720p",
				Size:    381023060,
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
