package hentaiff

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode Raw",
			URL:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-raw/",
			Want: 1,
		}, {
			Name: "Single Episode Eng Sub",
			URL:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-english-subbed/",
			Want: 1,
		}, {
			Name: "Single Episode Eng Dub",
			URL:  "https://hentaiff.com/a-kite-episode-02-english-dubbed/",
			Want: 1,
		}, {
			Name: "Single Episode Preview",
			URL:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-previews/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaiff.com/anime/a-kite/",
			Want: 6,
		}, {
			// this is the same logic for all extensions that group shows e.g. /genres/
			// its hard to make a test for the other groups since the number of episodes always changes
			Name: "Studio",
			URL:  "https://hentaiff.com/studio/arms/",
			// 5 show with a sum of 28 episodes
			Want: 28,
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
			URL:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-english-subbed/",
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
