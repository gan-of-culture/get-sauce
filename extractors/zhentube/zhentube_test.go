package zhentube

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
			URL:  "https://zhentube.com/soshite-watashi-wa-ojisan-ni-episode-4/",
			Want: 1,
		}, {
			Name: "Category",
			URL:  "https://zhentube.com/category/2021/",
			Want: 30,
		}, {
			Name: "Tag",
			URL:  "https://zhentube.com/tag/new-hentai-stream/",
			Want: 30,
		}, {
			Name: "Actor",
			URL:  "https://zhentube.com/actor/kotomi/",
			Want: 3,
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
				URL:     "https://zhentube.com/torokase-orgasm-episode-1/",
				Title:   "Torokase Orgasm Episode 1",
				Quality: "1920x1080",
				Size:    442901680,
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
