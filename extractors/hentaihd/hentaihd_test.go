package hentaihd

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
			URL:  "https://v2.hentaihd.net/74888/shuumatsu-no-harem-episode-1-raw/",
			Want: 1,
		}, {
			Name: "Single Episode Eng Sub",
			URL:  "https://v2.hentaihd.net/74923/shuumatsu-no-harem-episode-1-english-subbed/",
			Want: 1,
		}, {
			Name: "Single Episode Preview",
			URL:  "https://v2.hentaihd.net/74886/shuumatsu-no-harem-previews/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://v2.hentaihd.net/anime/accelerando-datenshi-tachi-no-sasayaki/",
			Want: 11,
		}, {
			// this is the same logic for all extensions that group shows e.g. /genres/
			// its hard to make a test for the other groups since the number of episodes always changes
			Name: "Studio",
			URL:  "https://v2.hentaihd.net/studio/flavors-soft/",
			Want: 23,
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
				URL:     "https://v2.hentaihd.net/77536/chizuru-chan-kaihatsu-nikki-episodio-6-en-espanol/",
				Title:   "Chizuru-chan Kaihatsu Nikki, Episodio 6 En Español",
				Quality: "1080p",
				Size:    279009210,
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
