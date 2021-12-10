package hentaihd

import "testing"

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
			Want: 11,
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
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://v2.hentaihd.net/77536/chizuru-chan-kaihatsu-nikki-episodio-6-en-espanol/",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
