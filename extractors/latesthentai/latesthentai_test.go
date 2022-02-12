package latesthentai

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
			Name: "Single Episode Eng Sub",
			URL:  "https://latesthentai.com/hajimete-no-hitozuma-episode-1/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://latesthentai.com/hentai/hajimete-no-hitozuma/",
			Want: 3,
		}, {
			// this is the same logic for all extensions that group shows e.g. /genres/
			// its hard to make a test for the other groups since the number of episodes always changes
			Name: "Studio",
			URL:  "https://latesthentai.com/studio/new-generation/",
			Want: 31,
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
				URL:     "https://latesthentai.com/hajimete-no-hitozuma-episode-1/",
				Title:   "Hajimete no Hitozuma Episode 1",
				Quality: "1280x720",
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
