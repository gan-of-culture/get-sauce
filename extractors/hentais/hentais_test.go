package hentais

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single default extraction",
			URL:  "https://www.hentais.tube/episodes/shishunki-sex-episode-4/",
			Want: 1,
		},
		{
			Name: "Whole default series extraction",
			URL:  "https://www.hentais.tube/tvshows/shishunki-sex/",
			Want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) != tt.Want {
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
			Name: "Single default extraction",
			URL:  "https://www.hentais.tube/episodes/shishunki-sex-episode-4",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) != tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
