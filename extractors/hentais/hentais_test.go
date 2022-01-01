package hentais

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
			Name: "Single default extraction",
			URL:  "https://www.hentais.tube/episodes/shishunki-sex-episode-4/",
			Want: 1,
		},
		{
			Name: "Whole default series extraction",
			URL:  "https://www.hentais.tube/tvshows/shishunki-sex/",
			Want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) != tt.Want {
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
				URL:     "https://www.hentais.tube/episodes/shishunki-sex-episode-4",
				Title:   "Shishunki Sex - Episode 4",
				Quality: "720p",
				Size:    164494351,
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
