package hentaiyes

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://hentaiyes.com/watch/hime-sama-love-life-episode-03/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaiyes.com/series/hime-sama-love-life/",
			Want: 3,
		}, {
			Name: "Tag",
			URL:  "https://hentaiyes.com/tag/1080p/",
			Want: 20,
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
			URL:  "https://hentaiyes.com/watch/kanpeki-ojou-sama-no-watakushi-ga-dogeza-de-mazo-ochisuru-choroin-na-wakenai-desu-wa-episode-03/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaiyes.com/series/akiko/",
			Want: 2,
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
