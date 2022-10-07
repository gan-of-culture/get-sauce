package hentaitv

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
			Name: "Single Episode Raw",
			URL:  "https://hentaitv.fun/episode/2509/deribari-chinko-o-tanomitai-onee-san-episode-1",
			Want: 1,
		}, {
			Name: "Single Episode Eng Sub",
			URL:  "https://hentaitv.fun/episode/2502/korashime-2-kyouikuteki-depaga-shidou-episode-1",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaitv.fun/hentai/1036/summer-inaka-no-seikatsu/",
			Want: 4,
		}, {
			Name: "Tag",
			URL:  "https://hentaitv.fun/list?genre[]=Large%20Breasts",
			Want: 25,
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
				URL:     "https://hentaitv.fun/episode/2503/summer-inaka-no-seikatsu-episode-2",
				Title:   "Summer: Inaka no Seikatsu Episode 2 English",
				Quality: "",
				Size:    60871679,
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
