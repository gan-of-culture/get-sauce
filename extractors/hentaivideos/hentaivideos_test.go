package hentaivideos

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
			URL:  "https://hentaivideos.net/ouji-no-honmei-wa-akuyaku-reijou-episode-4",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaivideos.net/hentai/ouji-no-honmei-wa-akuyaku-reijou",
			Want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want && tt.Want != 0 {
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
				URL:   "https://hentaivideos.net/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue-episode-3",
				Title: "Usamimi Bouken-tan: Sekuhara Shinagara Sekai o Sukue Episode 3",
				Size:  213691465,
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
