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
			Name: "Single Episode hanime.io/",
			URL:  "https://hanime.io/watch/torokase-orgasm-the-animation-episode-1/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hanime.io/hentai/kyonyuu-elf-oyako-saimin/",
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
				URL:     "https://hanime.io/watch/torokase-orgasm-the-animation-episode-1/",
				Title:   "Torokase Orgasm The Animation Episode 1",
				Quality: "1920x1080",
				Size:    485130240,
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
