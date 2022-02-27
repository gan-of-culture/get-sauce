package htstreaming

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
			Name: "Single Episode hentai.pro",
			URL:  "https://hentai.pro/knight-of-erin-episode-2/",
			Want: 1,
		}, {
			Name: "Overview hentai.pro",
			URL:  "https://hentai.pro/tag/breasts/",
			Want: 50,
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
			Name: "Single Episode hentai.pro",
			Args: test.Args{
				URL:     "https://hentai.pro/imaizumin-chi-wa-douyara-gal-no-tamariba-ni-natteru-rashii-episode-2/",
				Title:   "Imaizumin-chi wa Douyara Gal no Tamariba ni Natteru Rashii Episode 2",
				Quality: "1920x1080",
				Size:    288810112,
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
