package hentaimama

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
			URL:  "https://hentaimama.io/episodes/katainaka-ni-totsui-de-kita-russia-musume-h-shimakuru-ohanashi-episode-1/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaimama.io/tvshows/katainaka-ni-totsui-de-kita-russia-musume-h-shimakuru-ohanashi/",
			Want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want || len(URLs) < 1 {
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
			Name: "Single Episode using HLS only",
			Args: test.Args{
				URL:   "https://hentaimama.io/episodes/kuroinu-ii-animation-episode-1/",
				Title: "Kuroinu II The Animation Episode 1",
				Size:  172537611,
			},
		},
		{
			Name: "Single Episode using a single mp4 file",
			Args: test.Args{
				URL:   "https://hentaimama.io/episodes/ura-jutaijima-episode-1/",
				Title: "Ura Jutaijima Episode 1",
				Size:  77530809,
			},
		},
		{
			Name: "Single Episode using a both mp4 and HLS",
			Args: test.Args{
				URL:   "https://hentaimama.io/episodes/torokase-orgasm-animation-episode-1/",
				Title: "Torokase Orgasm The Animation Episode 1",
				Size:  186261816,
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
