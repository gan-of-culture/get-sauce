package hentaiplay

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
			Name: "Single Gallery",
			URL:  "https://hentaiplay.net/sister-breeder-episode-4/",
			Want: 1,
		}, {
			Name: "Episode-List",
			URL:  "https://hentaiplay.net/episode-list/a-forbidden-time/",
			Want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs, err := parseURL(tt.URL)
			test.CheckError(t, err)
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
			Name: "Single Gallery",
			Args: test.Args{
				URL:   "https://hentaiplay.net/sister-breeder-episode-4/",
				Title: "Sister Breeder! Episode 4 English",
				Size:  84322100,
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
